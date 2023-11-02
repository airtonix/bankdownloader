package sources

import (
	"errors"
	"time"

	"gopkg.in/yaml.v3"
)

type Source interface {
	LoadConfig(configNode *yaml.Node) error

	Process() error

	// function to login to the source
	Login() error

	// function to download the transactions
	DownloadTransactions(
		accountName string,
		accountNumber string,
		format string,
		fromDate time.Time,
		toDate time.Time,
	) (string, error)

	// function to open the browser
	OpenBrowser() error
}

func GetSourceFactory(source string) (Source, error) {
	switch source {
	case "anz":
		return &AnzSource{}, nil
	// case "commbank":
	// 	return &CommbankSource{}, nil
	default:
		return nil, errors.New("unsupported source")
	}
}

type SourceProps struct {
	Name string // name of the source
}

type SourceConfig struct {
	Domain      string `json:"domain" yaml:"domain"`           // the domain of the source
	Format      string `json:"format" yaml:"format"`           // format to download transactions in
	DaysToFetch int    `json:"daysToFetch" yaml:"daysToFetch"` // the number of days to fetch transactions for
}

func (s *SourceProps) Process() {
	// 	userhistory := store.GetHistory()

	// 	for _, account := range accounts {
	// 		// determine the next date to fetch transactions from
	// 		// based on the number of days to fetch
	// 		// and the last date we fetched transactions from
	// 		// if the last date is empty, then we fetch from today - daysToFetch
	// 		// if the last date is not empty, then we fetch from last date - daysToFetch
	// 		nextFromDate, err := userhistory.GetNextDate(
	// 			job.Source,
	// 			account.Number,
	// 			account.Name,
	// 			job.DaysToFetch,
	// 		)
	// 		if core.AssertErrorToNilf("could not get next date: %w", err) {
	// 			continue
	// 		}

	// 		nextToDate := core.ToStartOfDay(nextFromDate.AddDate(0, 0, job.DaysToFetch))

	// 		source.OpenBrowser()

	// 		// login to the source
	// 		err = source.Login(
	// 			job.Credentials,
	// 		)

	// 		// download the transactions
	// 		transactionFilename, err := source.DownloadTransactions(
	// 			account.Number,
	// 			account.Name,
	// 			job.Format,
	// 			nextFromDate,
	// 			nextToDate,
	// 		)

	// 		if core.AssertErrorToNilf("could not download transactions: %w", err) {
	// 			continue
	// 		}

	// 		log.Printf(
	// 			"Downloaded transactions for %s from %s to %s as %s",
	// 			account.Name, nextFromDate, nextToDate, transactionFilename,
	// 		)

	//		userhistory.SaveEvent(
	//			job.Source,
	//			account.Number,
	//			account.Name,
	//			nextFromDate,
	//			nextToDate,
	//		)
	//	}
}

// type Entity struct {
// 	Pw      *playwright.Playwright
// 	Browser playwright.Browser
// 	Page    playwright.Page
// }

// type NewSourceParams struct {
// 	Domain string
// }

// func (s *SourceProps) OpenBrowser() error {
// 	cwd := core.GetCwd()
// 	pw, err := playwright.Run()
// 	helpers.LogPanicln(err)

// 	s.Entity.Pw = pw

// 	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
// 		DownloadsPath: cwd,
// 		Headless:      playwright.Bool(true),
// 	})
// 	helpers.LogPanicln(err)
// 	s.Entity.Browser = browser

// 	page, err := browser.NewPage()
// 	helpers.LogPanicln(err)
// 	err = page.SetViewportSize(1920, 1440)
// 	helpers.LogPanicln(err)
// 	s.Entity.Page = page

// 	return nil
// }
