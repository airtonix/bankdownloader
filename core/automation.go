package core

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/gookit/color"
	"github.com/gosimple/slug"
	"github.com/sirupsen/logrus"

	"github.com/playwright-community/playwright-go"
	"github.com/yogiis/golang-web-automation/helpers"
)

type Automation struct {
	Pw      *playwright.Playwright
	Browser playwright.Browser
	Context playwright.BrowserContext
	Page    playwright.Page
}

func (a *Automation) OpenBrowser() error {
	cwd := GetCwd()
	downloadsPath := path.Join(cwd, "downloads")

	pw, err := playwright.Run(&playwright.RunOptions{
		Browsers: []string{"firefox"},
	})
	helpers.LogPanicln(err)

	a.Pw = pw

	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		DownloadsPath: &downloadsPath,
	})
	helpers.LogPanicln(err)
	a.Browser = browser

	context, err := browser.NewContext()
	helpers.LogPanicln(err)
	a.Context = context

	page, err := context.NewPage()
	helpers.LogPanicln(err)
	a.Page = page

	logrus.Info("Opened browser")

	return nil
}

func (a *Automation) CloseBrowser() error {
	err := a.Browser.Close()
	helpers.LogPanicln(err)

	logrus.Info("Closed browser")

	return nil
}

func (a *Automation) GetPageUrlObject() url.URL {
	page := a.Page
	obj, err := url.Parse(page.URL())
	if err != nil {
		logrus.Errorln("Could not parse url: ", err)
		return url.URL{}
	}
	return *obj
}

func (a *Automation) DownloadFile(
	downloadpath string,
	action func() error,
) (string, error) {
	page := a.Page
	download, err := page.ExpectDownload(action)
	if err == nil {
		TakeScreenshot(page, downloadpath)
	}

	logrus.Debugln("Target filename: ", downloadpath)
	targetDir, targetFilename := path.Split(downloadpath)

	storagePath := ResolveFileArg(
		"",
		"BANKDOWNLOADER_DOWNLOADDIR",
		path.Join("downloads", targetDir),
	)
	savedFilename := path.Join(storagePath, targetFilename)

	err = download.SaveAs(savedFilename)

	AssertErrorToNilf("could not save file: %w", err)

	return savedFilename, nil
}

func (a *Automation) PickElements(selector string) (playwright.Locator, error) {
	page := a.Page

	locator := page.Locator(selector)
	locator.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	if !AssertHasMatchingElements(locator, selector) {
		return nil, fmt.Errorf("could not pick elements matching: %s", selector)
	}
	return locator, nil
}
func (a *Automation) Find(selector string) error {
	locator, err := a.PickElements(selector)
	if err != nil {
		return err
	}

	logrus.Debugf(
		"[Found] %s > %s \n",
		selector,
		PrintMatchingElements(locator),
	)
	return nil
}
func (a *Automation) Click(selector string) error {
	locator, err := a.PickElements(selector)
	if err != nil {
		return err
	}
	locator.First().Click()
	logrus.Debugf(
		"[Clicked] %s > %s \n",
		selector,
		PrintMatchingElements(locator),
	)
	return nil
}
func (a *Automation) Focus(selector string) error {
	locator, err := a.PickElements(selector)
	if err != nil {
		return err
	}
	locator.First().Focus()
	logrus.Debugf(
		"[Focused] %s > %s \n",
		selector,
		PrintMatchingElements(locator),
	)
	return nil
}
func (a *Automation) Fill(selector string, value string) error {
	locator, err := a.PickElements(selector)
	if err != nil {
		return err
	}
	element := locator.First()
	element.Fill(value)
	logrus.Debugf(
		"[Focused] %s > %s \n",
		selector,
		PrintMatchingInputValues(locator),
	)
	return nil
}
func (a *Automation) FillSensitive(selector string, value string) error {
	locator, err := a.PickElements(selector)
	if err != nil {
		return err
	}
	element := locator.First()
	element.Fill(value)
	typedValue, err := element.InputValue()
	if err != nil {
		return err
	}
	logrus.Debugf(
		"[FilledSensitive] %s > %t \n",
		selector,
		typedValue == value,
	)

	return nil
}

func NewAutomation() *Automation {
	return &Automation{}
}

func CountOfElements(locator playwright.Locator) int {
	count, err := locator.Count()
	if err != nil {
		logrus.Errorln("Could not get count of elements: ", err)
		return 0
	}

	return count
}

func HasMatchingElements(locator playwright.Locator) bool {
	count := CountOfElements(locator)
	return count > 0
}

func AssertHasMatchingElements(locator playwright.Locator, itemName string) bool {
	if !HasMatchingElements(locator) {
		page, err := locator.Page()
		if err == nil {
			TakeScreenshot(page, itemName)
		}
		logrus.Panic(
			color.FgRed.Render(fmt.Sprintf("could not find item: %s", itemName)),
		)
		return false
	}
	return true
}

func TakeScreenshot(page playwright.Page, topic string) {
	cwd := GetCwd()
	if _, err := os.Stat(path.Join(cwd, "screenshots")); os.IsNotExist(err) {
		os.Mkdir(path.Join(cwd, "screenshots"), 0755)
	}
	screenshotPath := path.Join(cwd, "screenshots", fmt.Sprintf("%s.png", slug.Make(topic)))
	if _, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(screenshotPath),
		FullPage: playwright.Bool(true),
	}); err != nil {
		logrus.Errorln("Could not take screenshot: ", err)
	}
}

// function that prints the inner texts of matching elements
func PrintMatchingElements(locator playwright.Locator) string {
	elements, err := locator.AllInnerTexts()
	if err != nil {
		logrus.Errorln("Could not get elements: ", err)
		return ""
	}
	return JoinStrings(elements)
}

// function that prints the input values of matching elements
func PrintMatchingInputValues(locator playwright.Locator) string {
	elements, err := locator.All()
	if err != nil {
		logrus.Errorln("Could not get elements: ", err)
		return ""
	}
	var values []string
	for _, element := range elements {
		value, err := element.InputValue()
		if err != nil {
			logrus.Errorln("Could not get input value: ", err)
			return ""
		}
		values = append(values, value)
	}

	return JoinStrings(values)
}

// function that joins strings
func JoinStrings(strings []string) string {
	var result string
	for _, str := range strings {
		result = result + str
	}
	return result
}
