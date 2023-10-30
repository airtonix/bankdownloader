package sources

import (
	"github.com/airtonix/bank-downloaders/core"
	"github.com/playwright-community/playwright-go"
)

type Source struct {
	Name   string
	Domain string
	Config any
	Page   playwright.Page
}

type NewSourceParams struct {
	Domain string
}

func (s *Source) OpenBrowser() error {
	pw, err := playwright.Run()
	cwd := core.GetCwd()
	core.AssertErrorToNilf("could not launch playwright: %w", err)
	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		DownloadsPath: cwd,
		Headless:      playwright.Bool(true),
	})
	core.AssertErrorToNilf("could not launch browser: %w", err)
	context, err := browser.NewContext()
	core.AssertErrorToNilf("could not create context: %w", err)
	page, err := context.NewPage()
	core.AssertErrorToNilf("could not create page: %w", err)
	isConnected := context.Browser().IsConnected()

	core.LogLine("browser connected: %t", isConnected)
	s.Page = page

	return nil
}
