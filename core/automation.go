package core

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/gookit/color"
	"github.com/gosimple/slug"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/playwright-community/playwright-go"
	"github.com/yogiis/golang-web-automation/helpers"
)

type Automation struct {
	Pw      *playwright.Playwright
	Browser playwright.Browser
	Context playwright.BrowserContext
	Page    playwright.Page
}

func (this *Automation) OpenBrowser() error {
	cwd := GetCwd()
	downloadsPath := path.Join(cwd, "downloads")

	pw, err := playwright.Run(&playwright.RunOptions{
		Browsers: []string{"firefox"},
	})
	helpers.LogPanicln(err)

	this.Pw = pw

	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		DownloadsPath: &downloadsPath,
	})
	helpers.LogPanicln(err)
	this.Browser = browser

	context, err := browser.NewContext()
	helpers.LogPanicln(err)
	this.Context = context

	page, err := context.NewPage()
	helpers.LogPanicln(err)
	this.Page = page

	log.Info("Opened browser")

	return nil
}

func (this *Automation) CloseBrowser() error {
	err := this.Browser.Close()
	helpers.LogPanicln(err)

	log.Info("Closed browser")

	return nil
}

func (this *Automation) GetPageUrlObject() url.URL {
	page := this.Page
	obj, err := url.Parse(page.URL())
	if err != nil {
		logrus.Errorln("Could not parse url: ", err)
		return url.URL{}
	}
	return *obj
}

func (this *Automation) DownloadFile(
	filename string,
	action func() error,
) (string, error) {
	page := this.Page
	download, err := page.ExpectDownload(action)
	AssertErrorToNilf("could not expect download: %w", err)

	downloadedFilename := download.SuggestedFilename()
	ext := path.Ext(downloadedFilename)

	logrus.Debugln("Suggested filename: ", downloadedFilename)
	logrus.Debugln("Download ext: ", ext)

	storagePath := ResolveFileArg(
		filename,
		"BANKDOWNLOADER_DOWNLOADDIR",
		path.Join("downloads", filename),
	)
	savedFilename := path.Join(storagePath, fmt.Sprintf(
		"%s.%s", filename, ext,
	))
	err = download.SaveAs(savedFilename)

	AssertErrorToNilf("could not save file: %w", err)

	return savedFilename, nil
}

func (this *Automation) PickElements(selector string) (playwright.Locator, error) {
	page := this.Page

	locator := page.Locator(selector)
	locator.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	if !AssertHasMatchingElements(locator, selector) {
		return nil, fmt.Errorf("could not pick elements matching: %s", selector)
	}
	return locator, nil
}
func (this *Automation) Find(selector string) error {
	locator, err := this.PickElements(selector)
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
func (this *Automation) Click(selector string) error {
	locator, err := this.PickElements(selector)
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
func (this *Automation) Focus(selector string) error {
	locator, err := this.PickElements(selector)
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
func (this *Automation) Fill(selector string, value string) error {
	locator, err := this.PickElements(selector)
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
func (this *Automation) FillSensitive(selector string, value string) error {
	locator, err := this.PickElements(selector)
	if err != nil {
		return err
	}
	element := locator.First()
	element.Fill(value)
	typedValue, err := element.InputValue()
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
		log.Errorln("Could not get count of elements: ", err)
		return 0
	}

	return count
}

func HasMatchingElements(locator playwright.Locator) bool {
	count := CountOfElements(locator)
	if count > 0 {
		return true
	}
	return false
}

func AssertHasMatchingElements(locator playwright.Locator, itemName string) bool {
	if !HasMatchingElements(locator) {
		page, err := locator.Page()
		if err == nil {
			cwd := GetCwd()
			if _, err := os.Stat(path.Join(cwd, "screenshots")); os.IsNotExist(err) {
				os.Mkdir(path.Join(cwd, "screenshots"), 0755)
			}
			screenshotPath := path.Join(cwd, "screenshots", fmt.Sprintf("%s.png", slug.Make(itemName)))
			if _, err := page.Screenshot(playwright.PageScreenshotOptions{
				Path:     playwright.String(screenshotPath),
				FullPage: playwright.Bool(true),
			}); err != nil {
				logrus.Errorln("Could not take screenshot: ", err)
			}
		}
		logrus.Panic(
			color.FgRed.Render(fmt.Sprintf("could not find item: %s", itemName)),
		)
		return false
	}
	return true
}

// function that prints the inner texts of matching elements
func PrintMatchingElements(locator playwright.Locator) string {
	elements, err := locator.AllInnerTexts()
	if err != nil {
		log.Errorln("Could not get elements: ", err)
		return ""
	}
	return JoinStrings(elements)
}

// function that prints the input values of matching elements
func PrintMatchingInputValues(locator playwright.Locator) string {
	elements, err := locator.All()
	if err != nil {
		log.Errorln("Could not get elements: ", err)
		return ""
	}
	var values []string
	for _, element := range elements {
		value, err := element.InputValue()
		if err != nil {
			log.Errorln("Could not get input value: ", err)
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
