package core

import (
	"fmt"
	"net/url"
	"path"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/playwright-community/playwright-go"
	"github.com/yogiis/golang-web-automation/helpers"
)

type Automation struct {
	Pw      *playwright.Playwright
	Browser playwright.Browser
	Page    playwright.Page
}

func (s *Automation) OpenBrowser() error {
	cwd := GetCwd()
	downloadsPath := path.Join(cwd, "downloads")

	pw, err := playwright.Run()
	helpers.LogPanicln(err)

	s.Pw = pw

	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		DownloadsPath: &downloadsPath,
		Headless:      playwright.Bool(true),
	})
	helpers.LogPanicln(err)
	s.Browser = browser

	page, err := browser.NewPage()
	helpers.LogPanicln(err)

	s.Page = page

	log.Info("Opened browser")

	return nil
}

func (s *Automation) CloseBrowser() error {
	err := s.Browser.Close()
	helpers.LogPanicln(err)

	log.Info("Closed browser")

	return nil
}

func (s *Automation) GetPageUrlObject() url.URL {
	page := s.Page
	obj, err := url.Parse(page.URL())
	if err != nil {
		logrus.Errorln("Could not parse url: ", err)
		return url.URL{}
	}
	return *obj
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
			if _, err := page.Screenshot(playwright.PageScreenshotOptions{
				Path:     playwright.String(path.Join(cwd, "screenshots", itemName)),
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
