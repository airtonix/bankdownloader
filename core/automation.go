package core

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

type Automation struct {
	Context context.Context
	Cleanup context.CancelFunc
}

type AutomationOptionator func(*context.Context) context.Context

var (
	allocCtx context.Context
)

var allocateOnce sync.Once

func NewAutomation(
	options ...AutomationOptionator,
) *Automation {
	// Start the browser exactly once, as needed.
	allocateOnce.Do(func() {
		ctx, _ := chromedp.NewExecAllocator(
			context.Background(),
			chromedp.Headless,
			chromedp.NoSandbox,
		)

		allocCtx, _ = chromedp.NewContext(ctx)

		logrus.Infof("Allocated context: %v", &allocCtx)

		if err := chromedp.Run(allocCtx); err != nil {
			logrus.Panic(err)
		}

		chromedp.ListenBrowser(allocCtx, func(ev interface{}) {
			if ev, ok := ev.(*runtime.EventExceptionThrown); ok {
				logrus.Panicf("%+v\n", ev.ExceptionDetails)
			}
		})
	})

	// create a timeout as a safety net to prevent any infinite wait loops
	ctx, cleanup := context.WithTimeout(allocCtx, 60*time.Second)

	automation := &Automation{
		Context: ctx,
		Cleanup: cleanup,
	}

	return automation
}

func (a *Automation) CloseBrowser() {
	a.Cleanup()
}

func (a *Automation) SetViewportSize(width int64, height int64) error {
	logrus.Debugf("Setting viewport size to: %dx%d", width, height)
	err := chromedp.Run(a.Context,
		chromedp.Sleep(100*time.Millisecond),
		chromedp.EmulateViewport(width, height),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not set viewport size: %dx%d", width, height),
		err)

	return err
}

func (a *Automation) GetLocation() url.URL {
	var urlstr string
	err := chromedp.Run(a.Context,
		chromedp.Location(&urlstr),
	)
	if err != nil {
		logrus.Errorln("Could not get location: ", err)
		return url.URL{}
	}
	obj, err := url.Parse(urlstr)
	if err != nil {
		logrus.Errorln("Could not parse url: ", err)
		return url.URL{}
	}

	return *obj
}

func (a *Automation) Goto(url string) error {
	logrus.Debugf("Going to %s", url)

	// Navigate to the url and wait for the url to change
	err := chromedp.Run(a.Context,
		chromedp.Sleep(100*time.Millisecond),
		chromedp.Navigate(url),
		chromedp.WaitVisible("body"),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not goto: %s", url),
		err)

	logrus.Debugf("Went to %s", url)

	return err
}

func (a *Automation) Find(selector string) error {
	logrus.Debugf("Looking for %s", selector)
	err := chromedp.Run(a.Context,
		chromedp.Sleep(100*time.Millisecond),
		chromedp.WaitVisible(selector),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not find: %s", selector),
		err)
	logrus.Debugf("Found %s", selector)

	return err
}

func (a *Automation) Click(selector string) error {
	logrus.Debugf("Clicking %s", selector)
	err := chromedp.Run(a.Context,
		chromedp.Sleep(100*time.Millisecond),
		chromedp.Click(selector),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not click: %s", selector),
		err)
	logrus.Debugf("Clicked %s", selector)

	return err
}

func (a *Automation) Focus(selector string) error {
	logrus.Debugf("Focusing %s", selector)
	err := chromedp.Run(a.Context,
		chromedp.Sleep(100*time.Millisecond),
		chromedp.Focus(selector),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not focus: %s", selector),
		err)

	logrus.Debugf("Focused %s", selector)
	return err
}

func (a *Automation) Fill(selector string, value string) error {
	logrus.Debugf("Filling %s with %s", selector, value)
	err := chromedp.Run(a.Context,
		chromedp.Sleep(100*time.Millisecond),
		chromedp.WaitVisible(selector),
		chromedp.Sleep(1000),
		chromedp.SetValue(selector, value),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not fill: %s", selector),
		err)

	logrus.Debugf("Filled %s with %s", selector, value)

	return err
}

func (a *Automation) FillSensitive(selector string, value string) error {
	// make a string of stars the same length as the value
	stars := Stars(value)
	logrus.Debugf("Filling %s with %s", selector, stars)
	err := chromedp.Run(a.Context,
		chromedp.Sleep(100*time.Millisecond),
		chromedp.SetValue(selector, value),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not fill: %s", selector),
		err)

	logrus.Debugf("Filled %s with %s", selector, stars)

	return err
}

func (a *Automation) Pause(ms int) error {
	logrus.Debugf("Pausing for %d sec", ms)
	err := chromedp.Run(a.Context,
		chromedp.Sleep(time.Duration(ms)*time.Millisecond),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not pause: %d sec", ms),
		err)

	logrus.Debugf("Paused for %d sec", ms)

	return err
}

func (a *Automation) DownloadFile(
	downloadpath string,
	action func() error,
) (string, error) {
	logrus.Debugf("Downloading: %s", downloadpath)

	targetDir, targetFilename := path.Split(downloadpath)
	storagePath := ResolveFileArg(
		"",
		"BANKDOWNLOADER_DOWNLOADDIR",
		path.Join("downloads", targetDir),
	)
	savedFilename := path.Join(storagePath, targetFilename)
	is_downloaded := make(chan string, 1)

	// set up a listener to watch the download events and close the channel
	// when complete this could be expanded to handle multiple downloads
	// through creating a guid map, monitor download urls via
	// EventDownloadWillBegin, etc
	chromedp.ListenTarget(a.Context, func(v interface{}) {
		ev, ok := v.(*browser.EventDownloadProgress)
		if ok {
			completed := "(unknown)"
			if ev.TotalBytes != 0 {
				completed = fmt.Sprintf("%0.2f%%", ev.ReceivedBytes/ev.TotalBytes*100.0)
			}

			log.Printf("state: %s, completed: %s\n", ev.State.String(), completed)
			if ev.State == browser.DownloadProgressStateCompleted {
				is_downloaded <- ev.GUID
				close(is_downloaded)
			}
		}
	})
	logrus.Debugf("Listening for download event")

	err := chromedp.Run(a.Context,
		browser.
			SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(storagePath).
			WithEventsEnabled(true))
	AssertErrorToNilf("could not save file: %w", err)

	err = action()
	AssertErrorToNilf("Problem initiating download: %w", err)

	downloaded := <-is_downloaded
	downloadedPath := path.Join(storagePath, downloaded)

	// check if the file exists
	if _, err := os.Stat(downloadedPath); os.IsNotExist(err) {
		return "", fmt.Errorf("could not download file: %s", downloadedPath)
	}

	// move the file to the expected location
	if err := os.Rename(downloadedPath, savedFilename); err != nil {
		return "", fmt.Errorf("could not move file: %s", err)
	}
	logrus.Debugf("Downloaded: %s", savedFilename)

	return savedFilename, nil
}

var possibleChromePaths = []string{
	"chromium",
	"chromium-browser",
	"google-chrome",
	"google-chrome-stable",
	"google-chrome-beta",
}

func FindChrome() (string, error) {
	var err error

	// for each path, check if it exists
	for _, path := range possibleChromePaths {
		actualpath, err := exec.LookPath(path)
		if err == nil {
			logrus.Debugf("Found chrome: %s", actualpath)
			return actualpath, nil
		}
	}

	return "", fmt.Errorf("could not find chrome: %w", err)

}

func EnsureChromeExists() {
	_, err := FindChrome()
	AssertErrorToNilf("could not find chrome: %w", err)
	if err != nil {
		panic(err)
	}
}

// func TakeScreenshot(page playwright.Page, topic string) {
// 	cwd := GetCwd()
// 	if _, err := os.Stat(path.Join(cwd, "screenshots")); os.IsNotExist(err) {
// 		os.Mkdir(path.Join(cwd, "screenshots"), 0755)
// 	}
// 	screenshotPath := path.Join(cwd, "screenshots", fmt.Sprintf("%s.png", slug.Make(topic)))
// 	if _, err := page.Screenshot(playwright.PageScreenshotOptions{
// 		Path:     playwright.String(screenshotPath),
// 		FullPage: playwright.Bool(true),
// 	}); err != nil {
// 		logrus.Errorln("Could not take screenshot: ", err)
// 	}
// }
