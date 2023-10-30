package sources

import (
	"fmt"
	"time"

	"github.com/airtonix/bank-downloaders/core"
)

type AnzCredentials struct {
	// username is the customer registration number for ANZ internet banking.
	Username string `json:"username" yaml:"username"`
	// password is the password for ANZ internet banking.
	Password string `json:"password" yaml:"password"`
}

type AnzSource struct {
	*Source
}

func (source *AnzSource) Login(credentials any) error {
	page := source.Page

	loginDetails := credentials.(AnzCredentials)
	url :=
		fmt.Sprintf("%s/internetbanking", source.Domain)
	core.LogLine("visiting: %s", url)

	// start at the login page
	var _, err = page.Goto(url)
	core.AssertErrorToNilf(
		fmt.Sprintf("could not goto: %s", url),
		err)

	// Login
	page.GetByLabel("Customer Registration Number").Fill(loginDetails.Username)
	page.GetByLabel("Password").Fill(loginDetails.Password)
	page.Locator("//button[type='submit']").Click()

	// wait for the account page to load
	page.Locator("//h1[text()='Your accounts']")

	return nil
}

func (source *AnzSource) DownloadTransactions(
	accountNumber string,
	accountName string,
	format string,
	fromDate time.Time,
	toDate time.Time,
) (string, error) {
	page := source.Page

	domain, err := core.GetDomainFromUrl(page.URL())

	page.Goto(domain + "/IBUI/#/download-transaction")

	// pick the account by clicking the label "Account"
	page.Locator("//label[text()='Account']").Click()

	// then click the ul > li named account number entry that is the adjacent sibling of the labels parent parent
	page.Locator("//label[text()='Account']/parent::div/parent::div/following-sibling::ul/li[text()='" + accountNumber + "']").Click()

	// change to date range mode

	page.Locator("//[@data-testid='Date range']").Click()

	// select the date range fromDate
	page.Locator("//input[@data-testid='fromdate-textfield']").Fill(fromDate.Format("12/09/2023"))

	// select the date range toDate
	page.Locator("//input[@data-testid='todate-textfield']").Fill(toDate.Format("12/09/2023"))

	// select the downlaod format by clicking the label "Software package"
	page.Locator("//label[text()='Software package']").Click()
	page.Locator("//label[text()='Software package']/parent::div/parent::div/following-sibling::ul/li[text()='" + format + "']").Click()

	// click the download button
	download, err := page.ExpectDownload(func() error {
		return page.Locator("//button[text()='Download']").Click()
	})

	core.AssertErrorToNilf("could not download: %w", err)
	filename := download.SuggestedFilename()
	download.SaveAs(filename)

	return filename, nil
}

// ensure AnzSource satisfies the SourceCommand interface
var _ SourceCommand = &AnzSource{}

func NewAnzSource(params NewSourceParams) SourceCommand {
	return &AnzSource{
		Source: &Source{
			Name:   "anz",
			Domain: params.Domain,
		},
	}
}
