package core

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

type Automation struct {
	Context context.Context
	Cleanup context.CancelFunc
	Cancel  context.CancelFunc
}

type NewAutomationOptions struct {
	Headless bool
}

func NewAutomation() *Automation {

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	// create a timeout as a safety net to prevent any infinite wait loops
	ctx, cleanup := context.WithTimeout(ctx, 60*time.Second)

	automation := &Automation{
		Context: ctx,
		Cleanup: cleanup,
		Cancel:  cancel,
	}

	return automation
}

func (a *Automation) Close() {
	a.Cleanup()
	// a.Cancel()
}

func (a *Automation) SetViewportSize(width int64, height int64) error {
	err := chromedp.Run(a.Context,
		chromedp.EmulateViewport(width, height),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not set viewport size: %dx%d", width, height),
		err)

	return err
}

func (a *Automation) GetLocation() url.URL {
	var urlstr string
	chromedp.Run(a.Context,
		chromedp.Location(&urlstr),
	)
	obj, err := url.Parse(urlstr)
	if err != nil {
		logrus.Errorln("Could not parse url: ", err)
		return url.URL{}
	}

	return *obj
}

func (a *Automation) Goto(url string) error {
	err := chromedp.Run(a.Context,
		chromedp.Navigate(url),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not goto: %s", url),
		err)

	return err
}

func (a *Automation) Find(selector string) error {
	err := chromedp.Run(a.Context,
		chromedp.WaitVisible(selector),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not find: %s", selector),
		err)

	return err
}

func (a *Automation) Click(selector string) error {
	err := chromedp.Run(a.Context,
		chromedp.Click(selector),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not click: %s", selector),
		err)

	return err
}

func (a *Automation) Focus(selector string) error {
	err := chromedp.Run(a.Context,
		chromedp.Focus(selector),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not focus: %s", selector),
		err)

	return err
}

func (a *Automation) Fill(selector string, value string) error {
	err := chromedp.Run(a.Context,
		chromedp.WaitVisible(selector),
		chromedp.Sleep(1000),
		chromedp.SetValue(selector, value),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not fill: %s", selector),
		err)

	return err
}

func (a *Automation) FillSensitive(selector string, value string) error {
	err := chromedp.Run(a.Context,
		chromedp.SetValue(selector, value),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not fill: %s", selector),
		err)

	return err
}

func (a *Automation) Pause(duration time.Duration) error {
	err := chromedp.Run(a.Context,
		chromedp.Sleep(duration),
	)

	AssertErrorToNilf(
		fmt.Sprintf("could not pause: %s", duration),
		err)

	return err
}

func (a *Automation) DownloadFile(
	downloadpath string,
	action func() error,
) (string, error) {
	logrus.Debugln("Target filename: ", downloadpath)
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

	err := chromedp.Run(a.Context,
		browser.
			SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(storagePath).
			WithEventsEnabled(true))
	AssertErrorToNilf("could not save file: %w", err)

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

// func (a *Automation) PickElements(selector string) (playwright.Locator, error) {
// 	page := a.Page

// 	locator := page.Locator(selector)
// 	locator.WaitFor(playwright.LocatorWaitForOptions{
// 		State: playwright.WaitForSelectorStateAttached,
// 	})
// 	if !AssertHasMatchingElements(locator, selector) {
// 		return nil, fmt.Errorf("could not pick elements matching: %s", selector)
// 	}
// 	return locator, nil
// }
// func (a *Automation) Find(selector string) error {
// 	locator, err := a.PickElements(selector)
// 	if err != nil {
// 		return err
// 	}

// 	logrus.Debugf(
// 		"[Found] %s > %s \n",
// 		selector,
// 		PrintMatchingElements(locator),
// 	)
// 	return nil
// }
// func (a *Automation) Click(selector string) error {
// 	locator, err := a.PickElements(selector)
// 	if err != nil {
// 		return err
// 	}
// 	locator.First().Click()
// 	logrus.Debugf(
// 		"[Clicked] %s > %s \n",
// 		selector,
// 		PrintMatchingElements(locator),
// 	)
// 	return nil
// }
// func (a *Automation) Focus(selector string) error {
// 	locator, err := a.PickElements(selector)
// 	if err != nil {
// 		return err
// 	}
// 	locator.First().Focus()
// 	logrus.Debugf(
// 		"[Focused] %s > %s \n",
// 		selector,
// 		PrintMatchingElements(locator),
// 	)
// 	return nil
// }
// func (a *Automation) Fill(selector string, value string) error {
// 	locator, err := a.PickElements(selector)
// 	if err != nil {
// 		return err
// 	}
// 	element := locator.First()
// 	element.Fill(value)
// 	logrus.Debugf(
// 		"[Focused] %s > %s \n",
// 		selector,
// 		PrintMatchingInputValues(locator),
// 	)
// 	return nil
// }
// func (a *Automation) FillSensitive(selector string, value string) error {
// 	locator, err := a.PickElements(selector)
// 	if err != nil {
// 		return err
// 	}
// 	element := locator.First()
// 	element.Fill(value)
// 	typedValue, err := element.InputValue()
// 	if err != nil {
// 		return err
// 	}
// 	logrus.Debugf(
// 		"[FilledSensitive] %s > %t \n",
// 		selector,
// 		typedValue == value,
// 	)

// 	return nil
// }

// func CountOfElements(locator playwright.Locator) int {
// 	count, err := locator.Count()
// 	if err != nil {
// 		logrus.Errorln("Could not get count of elements: ", err)
// 		return 0
// 	}

// 	return count
// }

// func HasMatchingElements(locator playwright.Locator) bool {
// 	count := CountOfElements(locator)
// 	return count > 0
// }

// func AssertHasMatchingElements(locator playwright.Locator, itemName string) bool {
// 	if !HasMatchingElements(locator) {
// 		page, err := locator.Page()
// 		if err == nil {
// 			TakeScreenshot(page, itemName)
// 		}
// 		logrus.Panic(
// 			color.FgRed.Render(fmt.Sprintf("could not find item: %s", itemName)),
// 		)
// 		return false
// 	}
// 	return true
// }

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

// // function that prints the inner texts of matching elements
// func PrintMatchingElements(locator playwright.Locator) string {
// 	elements, err := locator.AllInnerTexts()
// 	if err != nil {
// 		logrus.Errorln("Could not get elements: ", err)
// 		return ""
// 	}
// 	return JoinStrings(elements)
// }

// // function that prints the input values of matching elements
// func PrintMatchingInputValues(locator playwright.Locator) string {
// 	elements, err := locator.All()
// 	if err != nil {
// 		logrus.Errorln("Could not get elements: ", err)
// 		return ""
// 	}
// 	var values []string
// 	for _, element := range elements {
// 		value, err := element.InputValue()
// 		if err != nil {
// 			logrus.Errorln("Could not get input value: ", err)
// 			return ""
// 		}
// 		values = append(values, value)
// 	}

// 	return JoinStrings(values)
// }

// function that joins strings
func JoinStrings(strings []string) string {
	var result string
	for _, str := range strings {
		result = result + str
	}
	return result
}
