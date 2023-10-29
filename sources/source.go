package sources

import (
	"github.com/airtonix/bank-downloaders/core"
	"github.com/playwright-community/playwright-go"
)

type Source struct {
	//
	Name string

	//
	Config any

	//
	Page playwright.Page
}

func (s *Source) OpenBrowser() error {
	pw, err := playwright.Run()
	core.AssertErrorToNilf("could not launch playwright: %w", err)
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	core.AssertErrorToNilf("could not launch Chromium: %w", err)
	context, err := browser.NewContext()
	core.AssertErrorToNilf("could not create context: %w", err)
	page, err := context.NewPage()
	core.AssertErrorToNilf("could not create page: %w", err)

	s.Page = page

	return nil
}
