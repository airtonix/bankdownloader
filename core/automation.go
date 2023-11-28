package core

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

type Automation struct {
	Context          context.Context
	AllocatorOptions []chromedp.ExecAllocatorOption
	Cleanup          context.CancelFunc
	Delay            time.Duration
	Fs               Filesystem
}

type AutomationOptionator func(automation *Automation)

func WithHeadless() AutomationOptionator {
	return func(automation *Automation) {
		automation.AllocatorOptions = append(automation.AllocatorOptions, chromedp.Headless)
	}
}

func WithNoSandbox() AutomationOptionator {
	return func(automation *Automation) {
		automation.AllocatorOptions = append(automation.AllocatorOptions, chromedp.NoSandbox)
	}
}

func WithDownloadMemFs() AutomationOptionator {
	return func(automation *Automation) {
		automation.Fs = NewMockFilesystem("/home/user")
	}
}

// var allocateOnce sync.Once

func NewAutomation(
	options ...AutomationOptionator,
) *Automation {
	automation := &Automation{
		Context: context.Background(),
		Delay:   200 * time.Millisecond,
		Fs:      NewRealFilesystem(),
		AllocatorOptions: []chromedp.ExecAllocatorOption{
			chromedp.NoFirstRun,
			chromedp.NoDefaultBrowserCheck,

			// After Puppeteer's default behavior.
			chromedp.Flag("disable-background-networking", true),
			chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
			chromedp.Flag("disable-background-timer-throttling", true),
			chromedp.Flag("disable-backgrounding-occluded-windows", true),
			chromedp.Flag("disable-breakpad", true),
			chromedp.Flag("disable-client-side-phishing-detection", true),
			chromedp.Flag("disable-default-apps", true),
			chromedp.Flag("disable-dev-shm-usage", true),
			chromedp.Flag("disable-extensions", true),
			chromedp.Flag("disable-features", "site-per-process,Translate,BlinkGenPropertyTrees"),
			chromedp.Flag("disable-hang-monitor", true),
			chromedp.Flag("disable-ipc-flooding-protection", true),
			chromedp.Flag("disable-popup-blocking", true),
			chromedp.Flag("disable-prompt-on-repost", true),
			chromedp.Flag("disable-renderer-backgrounding", true),
			chromedp.Flag("disable-sync", true),
			chromedp.Flag("force-color-profile", "srgb"),
			chromedp.Flag("metrics-recording-only", true),
			chromedp.Flag("safebrowsing-disable-auto-update", true),
			chromedp.Flag("enable-automation", true),
			chromedp.Flag("password-store", "basic"),
			chromedp.Flag("use-mock-keychain", true),
		},
	}

	// loop over options and apply
	for _, option := range options {
		option(automation)
	}

	ctx, _ := chromedp.NewExecAllocator(
		automation.Context,
		automation.AllocatorOptions...,
	)

	// create a timeout as a safety net to prevent any infinite wait loops
	ctx, cleanup := context.WithTimeout(ctx, 60*time.Second)
	automation.Context, _ = chromedp.NewContext(ctx)

	logrus.Infof("Allocated context: %v", &automation.Context)

	if err := chromedp.Run(automation.Context); err != nil {
		cleanup()
		logrus.Panic(err)
	}
	automation.Cleanup = cleanup

	chromedp.ListenBrowser(automation.Context, func(ev interface{}) {
		if ev, ok := ev.(*runtime.EventExceptionThrown); ok {
			logrus.Panicf("%+v\n", ev.ExceptionDetails)
		}
	})

	return automation
}

func (a *Automation) CloseBrowser() {
	a.Cleanup()
}

func (a *Automation) SetViewportSize(width int64, height int64) {
	logrus.Debugf("Setting viewport size to: %dx%d", width, height)
	err := chromedp.Run(a.Context,
		chromedp.Sleep(100*time.Millisecond),
		chromedp.EmulateViewport(width, height),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not set viewport size: %dx%d", width, height),
		err)
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

func (a *Automation) Goto(url string) {
	logrus.Debugf("Going to %s", url)

	// Navigate to the url and wait for the url to change
	err := chromedp.Run(a.Context,
		chromedp.Sleep(a.Delay),
		chromedp.Navigate(url),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not goto: %s", url),
		err)

	logrus.Debugf("Went to %s", url)
}

func (a *Automation) Find(selector string) {
	logrus.Debugf("Looking for %s", selector)
	err := chromedp.Run(a.Context,
		chromedp.Sleep(a.Delay),
		chromedp.WaitVisible(selector),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not find: %s", selector),
		err)
	logrus.Debugf("Found %s", selector)
}

func (a *Automation) Click(selector string) {
	logrus.Debugf("Clicking %s", selector)
	err := chromedp.Run(a.Context,
		chromedp.Sleep(a.Delay),
		chromedp.Click(selector),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not click: %s", selector),
		err)
	logrus.Debugf("Clicked %s", selector)
}

func (a *Automation) Focus(selector string) {
	logrus.Debugf("Focusing %s", selector)
	err := chromedp.Run(a.Context,
		chromedp.Sleep(a.Delay),
		chromedp.Focus(selector),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not focus: %s", selector),
		err)

	logrus.Debugf("Focused %s", selector)
}

func (a *Automation) Fill(selector string, value string) {
	logrus.Debugf("Filling %s with %s", selector, value)
	err := chromedp.Run(a.Context,
		chromedp.Sleep(a.Delay),
		chromedp.WaitVisible(selector),
		chromedp.Sleep(a.Delay),
		chromedp.SendKeys(selector, value),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not fill: %s", selector),
		err)

	logrus.Debugf("Filled %s with %s", selector, value)
}

func (a *Automation) FillSensitive(selector string, value string) {
	// make a string of stars the same length as the value
	stars := Stars(value)
	logrus.Debugf("Filling %s with %s", selector, stars)
	err := chromedp.Run(a.Context,
		chromedp.Sleep(a.Delay),
		chromedp.WaitVisible(selector),
		chromedp.Sleep(a.Delay),
		chromedp.SendKeys(selector, value),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not fill: %s", selector),
		err)

	logrus.Debugf("Filled %s with %s", selector, stars)
}

func (a *Automation) Pause(ms int) {
	logrus.Debugf("Pausing for %d sec", ms)
	err := chromedp.Run(a.Context,
		chromedp.Sleep(time.Duration(ms)*time.Millisecond),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not pause: %d sec", ms),
		err)

	logrus.Debugf("Paused for %d sec", ms)
}

func (a *Automation) DownloadFile(
	downloadpath string,
	action func() error,
) (string, error) {
	resolvedDownloadPath := a.Fs.ExpandPathWithHome(downloadpath)
	logrus.Debugf("Downloading: %s", downloadpath)

	targetDir, targetFilename := path.Split(
		resolvedDownloadPath,
	)
	savedFilename := path.Join(targetDir, targetFilename)
	download_queue := make(chan string, 1)

	// set up a listener to watch the download events
	// filename downloaded is communicated through the channel
	chromedp.ListenTarget(a.Context, func(v interface{}) {
		ev, ok := v.(*browser.EventDownloadProgress)
		if ok {
			completed := "(unknown)"
			if ev.TotalBytes != 0 {
				completed = fmt.Sprintf("%0.2f%%", ev.ReceivedBytes/ev.TotalBytes*100.0)
			}

			logrus.Debugf("state: %s, completed: %s\n", ev.State.String(), completed)
			if ev.State == browser.DownloadProgressStateCompleted {
				download_queue <- ev.GUID
			}
		}
	})
	logrus.Debugf("Listening for download event")

	if err := chromedp.Run(a.Context,
		browser.
			SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(targetDir).
			WithEventsEnabled(true),
	); err != nil && !strings.Contains(err.Error(), "net::ERR_ABORTED") {
		AssertErrorToNilf(
			fmt.Sprintf("could not set download behavior: %s", targetDir),
			err)
	}

	AssertErrorToNilf("Problem initiating download: %w", action())

	downloaded_filename := <-download_queue
	downloadedPath := path.Join(targetDir, downloaded_filename)

	// check if the file exists
	if _, err := a.Fs.GetFs().Stat(downloadedPath); os.IsNotExist(err) {
		return "", fmt.Errorf("could not download file: %s", downloadedPath)
	}

	// move the file to the expected location
	if err := a.Fs.GetFs().Rename(downloadedPath, savedFilename); err != nil {
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
