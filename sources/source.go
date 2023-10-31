package sources

import (
	"github.com/airtonix/bank-downloaders/core"
	"github.com/playwright-community/playwright-go"
	"github.com/yogiis/golang-web-automation/helpers"
)

type Entity struct {
	Pw      *playwright.Playwright
	Browser playwright.Browser
	Page    playwright.Page
}

type Source struct {
	Name    string // name of the source
	Domain  string // domain name of the source
	Config  any    // configuration for the source
	*Entity        // entity for the source
}

type NewSourceParams struct {
	Domain string
}

func (s *Source) OpenBrowser() error {
	cwd := core.GetCwd()
	pw, err := playwright.Run()
	helpers.LogPanicln(err)

	s.Entity.Pw = pw

	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		DownloadsPath: cwd,
		Headless:      playwright.Bool(true),
	})
	helpers.LogPanicln(err)
	s.Entity.Browser = browser

	page, err := browser.NewPage()
	helpers.LogPanicln(err)
	err = page.SetViewportSize(1920, 1440)
	helpers.LogPanicln(err)
	s.Entity.Page = page

	return nil
}
