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
	PageObjects AnzPageObjects
}

type AnzPageObjects struct {
	LoginHeader                        string `json:"login_header" yaml:"login_header"`
	LoginUsernameInput                 string `json:"login_username_input" yaml:"login_username_input"`
	LoginPasswordInput                 string `json:"login_password_input" yaml:"login_password_input"`
	LoginButton                        string `json:"login_button" yaml:"login_button"`
	AccountsPageHeader                 string `json:"accounts_page_header" yaml:"accounts_page_header"`
	ExportAccountDropdownLabel         string `json:"export_account_dropdown_label" yaml:"export_account_dropdown_label"`
	ExportAccountDropdownOption        string `json:"export_account_dropdown_option" yaml:"export_account_dropdown_option"`
	ExportDateRangeModeButton          string `json:"export_date_range_mode_button" yaml:"export_date_range_mode_button"`
	ExportDateRangeFromDateInput       string `json:"export_from_date_label" yaml:"export_from_date_label"`
	ExportDateRangeToDateInput         string `json:"export_to_date_label" yaml:"export_to_date_label"`
	ExportDownloadFormatDropdownLabel  string `json:"export_download_format_dropdown_label" yaml:"export_download_format_dropdown_label"`
	ExportDownloadFormatDropdownOption string `json:"export_download_format_dropdown_option" yaml:"export_download_format_dropdown_option"`
	ExportDownloadButton               string `json:"export_download_button" yaml:"export_download_button"`
}

var pageObjects = AnzPageObjects{
	LoginHeader:                        "//h1[text()='Log in to ANZ Internet Banking']",
	LoginUsernameInput:                 "//input[@name='customerRegistrationNumber']",
	LoginPasswordInput:                 "//input[@name='password']",
	LoginButton:                        "//button[@type='submit']",
	AccountsPageHeader:                 "//h1[text()='Your accounts']",
	ExportAccountDropdownLabel:         "//label[@for='drop-down-search-transaction-account1-dropdown-field'][@text='Account']",
	ExportAccountDropdownOption:        "//ul[@data-test-id='drop-down-search-transaction-account1-dropdown-results'/li[@role='option'][text()='%s']",
	ExportDateRangeModeButton:          "//[@aria-label='Search period'][@role='tablist']/li[@role='tab'][contains(., 'Date range')]",
	ExportDateRangeFromDateInput:       "//[@id='Date rangepanel'][role='tappanel']/input[@id='fromdate-textfield']",
	ExportDateRangeToDateInput:         "//[@id='Date rangepanel'][role='tappanel']/label[@id='todate-textfield']",
	ExportDownloadFormatDropdownLabel:  "//label[@for='drop-down-search-software-dropdown-field'][text()='Software package']",
	ExportDownloadFormatDropdownOption: "//ul[@data-test-id='drop-down-search-software-dropdown-results']/li[@role='option'][text()='%s']",
	ExportDownloadButton:               "//button[contains(., 'Download')]",
}

func (source *AnzSource) Login(credentials any) error {
	var err error
	page := source.Page

	loginDetails := credentials.(AnzCredentials)
	url :=
		fmt.Sprintf("%s/internetbanking", source.Domain)
	core.LogLine("visiting: %s", url)

	// start at the login page
	_, err = page.Goto(url)
	core.AssertErrorToNilf(
		fmt.Sprintf("could not goto: %s", url),
		err)

	// wait for the login page to load
	page.Locator(pageObjects.LoginHeader)

	// Login
	page.Locator(pageObjects.LoginUsernameInput).Fill(loginDetails.Username)
	page.Locator(pageObjects.LoginPasswordInput).Fill(loginDetails.Password)
	core.LogLine("authenticating...")

	page.Locator(pageObjects.LoginButton).Click()

	// wait for the account page to load
	page.Locator(pageObjects.AccountsPageHeader)
	core.LogLine("authenticated")
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
	var err error

	url := fmt.Sprintf("%s/IBUI/#/download-transaction", source.Domain)
	core.LogLine("visiting: %s", url)

	_, err = page.Goto(url)
	core.AssertErrorToNilf(
		fmt.Sprintf("could not goto: %s", url),
		err)

	page.BringToFront()

	// pick the account by clicking the label "Account"
	page.Locator(pageObjects.ExportAccountDropdownLabel).Click()
	// then click the ul > li named account number entry that is the adjacent sibling of the labels parent parent
	page.Locator(fmt.Sprintf(pageObjects.ExportAccountDropdownOption, accountNumber)).Click()
	core.LogLine("selected account: %s", accountNumber)

	// change to date range mode
	page.Locator(pageObjects.ExportDateRangeModeButton).Click()
	// select the date range fromDate
	page.Locator(pageObjects.ExportDateRangeFromDateInput).Fill(fromDate.Format("12/09/2023"))
	// select the date range toDate
	page.Locator(pageObjects.ExportDateRangeToDateInput).Fill(toDate.Format("12/09/2023"))
	core.LogLine("selected date range: %s - %s", fromDate.Format("12/09/2023"), toDate.Format("12/09/2023"))

	// select the downlaod format by clicking the label "Software package"
	page.Locator(pageObjects.ExportDownloadFormatDropdownLabel).Click()
	page.Locator(fmt.Sprintf(pageObjects.ExportDownloadFormatDropdownOption, format)).Click()
	core.LogLine("selected format: %s", format)

	// click the download button
	download, err := page.ExpectDownload(func() error {
		core.LogLine("downloading %s", "transactions")
		return page.Locator(pageObjects.ExportDownloadButton).Click()
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
