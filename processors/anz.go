package processors

import (
	"fmt"
	"time"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/kr/pretty"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
)

type AnzProcessor struct {
	Config *AnzConfig
	*Processor
}

// ensure that AnzProcessor implements the Processor interface
var _ IProcessor = (*AnzProcessor)(nil)

func (processor *AnzProcessor) GetFormat() string {
	return processor.Config.Format
}

func (processor *AnzProcessor) GetDaysToFetch() int {
	return processor.Config.DaysToFetch
}

func (processor *AnzProcessor) Render() error {
	pretty.Println(processor)
	return nil
}

func (processor *AnzProcessor) Login(automation *core.Automation) error {
	var err error
	page := automation.Page

	loginDetails := processor.Config.Credentials

	url := fmt.Sprintf(
		"%s/internetbanking",
		processor.Config.Domain,
	)
	logrus.Info("logging into ", url)

	page.BringToFront()

	// start at the login page
	_, err = page.Goto(url)
	core.AssertErrorToNilf(
		fmt.Sprintf("could not goto: %s", url),
		err)
	page.SetViewportSize(1200, 900)

	logrus.Debugln("waiting for login page to load...")
	// wait for the login page to load
	header := page.Locator(pageObjects.LoginHeader)
	header.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	core.AssertHasMatchingElements(header, "LoginHeader")
	logrus.Debugln("header", core.PrintMatchingElements(header))

	// Login
	logrus.Debugln("finding login fields...")

	// Username
	usernameField := page.Locator(pageObjects.LoginUsernameInput)
	usernameField.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	core.AssertHasMatchingElements(usernameField, "LoginUsernameInput")
	usernameField.Fill(loginDetails.Username)

	// Password
	passwordField := page.Locator(pageObjects.LoginPasswordInput)
	passwordField.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	core.AssertHasMatchingElements(usernameField, "LoginPasswordInput")
	passwordField.Fill(loginDetails.Password)

	// LoginButton
	loginButton := page.Locator(pageObjects.LoginButton)
	loginButton.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	core.AssertHasMatchingElements(loginButton, "LoginButton")
	loginButton.Click()

	logrus.Info("authenticating...")

	// Accounts Page
	// wait for the account page to load
	accountsPageHeader := page.Locator(pageObjects.AccountsPageHeader)
	accountsPageHeader.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	core.AssertHasMatchingElements(accountsPageHeader, "accounts page header")
	logrus.Debugln("accounts header", core.PrintMatchingElements(accountsPageHeader))
	logrus.Info("authenticated")
	return nil
}

func (processor *AnzProcessor) DownloadTransactions(
	accountName string,
	accountNumber string,
	fromDate time.Time,
	toDate time.Time,
	automation *core.Automation,
) (string, error) {
	var err error
	page := automation.Page

	var format = processor.Config.Format
	fromDateString := fromDate.Format("02/01/2006")
	toDateString := toDate.Format("02/01/2006")

	// get the current hostname for the current page
	pageUrl := automation.GetPageUrlObject()

	url := fmt.Sprintf(
		"%s://%s/IBUI/#/download-transaction",
		pageUrl.Scheme,
		pageUrl.Host,
	)

	logrus.Infoln(
		fmt.Sprintf(
			"Fetching transactions from: %s \n %s [%s]: %s - %s",
			url,
			accountName,
			accountNumber,
			fromDateString,
			toDateString,
		),
	)

	_, err = page.Goto(url)
	core.AssertErrorToNilf(
		fmt.Sprintf("could not goto: %s", url),
		err)
	page.SetViewportSize(1200, 900)

	// Transactions Page
	// wait for the page to load
	exportPageHeader := page.Locator(pageObjects.ExportPageHeader)
	exportPageHeader.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	core.AssertHasMatchingElements(exportPageHeader, "ExportPageHeader")
	logrus.Debugln("ExportPageHeader >", core.PrintMatchingElements(exportPageHeader))

	// pick the account by clicking the label "Account"
	accountDropdownLabel := page.Locator(pageObjects.ExportAccountDropdownLabel)
	accountDropdownLabel.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	core.AssertHasMatchingElements(accountDropdownLabel, "ExportAccountDropdownLabel")
	logrus.Debug("ExportAccountDropdownLabel >", core.PrintMatchingElements(accountDropdownLabel))
	accountDropdownLabel.Click()

	// then click the ul > li named account number entry that is the adjacent sibling of the labels parent parent
	accountDropdownOption := page.Locator(fmt.Sprintf(pageObjects.ExportAccountDropdownOption, accountNumber))
	accountDropdownOption.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	core.AssertHasMatchingElements(accountDropdownOption, "ExportAccountDropdownOption")
	logrus.Debug("ExportAccountDropdownOption > ", core.PrintMatchingElements(accountDropdownOption))
	accountDropdownOption.Click()
	logrus.Debug("selected account: ", accountNumber)

	// change to date range mode
	exportDateRangeMode := page.Locator(pageObjects.ExportDateRangeModeButton)
	exportDateRangeMode.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	core.AssertHasMatchingElements(exportDateRangeMode, "ExportDateRangeModeButton")
	logrus.Debug("ExportDateRangeModeButton > ", core.PrintMatchingElements(exportDateRangeMode))
	exportDateRangeMode.Click()

	// select the date range fromDate
	fromDateInput := page.Locator(pageObjects.ExportDateRangeFromDateInput)
	fromDateInput.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	core.AssertHasMatchingElements(fromDateInput, "ExportDateRangeFromDateInput")
	logrus.Debug("ExportDateRangeFromDateInput > ", core.PrintMatchingElements(fromDateInput))
	fromDateInput.Fill(fromDateString)
	logrus.Debug("ExportDateRangeFromDateInput.value > ", core.PrintMatchingInputValues(fromDateInput))

	// select the date range toDate
	toDateInput := page.Locator(pageObjects.ExportDateRangeToDateInput)
	toDateInput.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	core.AssertHasMatchingElements(toDateInput, "ExportDateRangeToDateInput")
	logrus.Debug("ExportDateRangeToDateInput > ", core.PrintMatchingElements(toDateInput))
	toDateInput.Fill(toDateString)
	logrus.Debug("ExportDateRangeToDateInput.value > ", core.PrintMatchingInputValues(toDateInput))

	logrus.Debugf(
		"selected date range: %s - %s",
		fromDateString, toDateString,
	)
	// close the date popups
	exportPageHeader.Click()

	// select the downlaod format by clicking the label "Software package"
	formatDropdownLabel := page.Locator(pageObjects.ExportDownloadFormatDropdownLabel)
	formatDropdownLabel.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateAttached,
		Timeout: playwright.Float(1000),
	})
	core.AssertHasMatchingElements(formatDropdownLabel, "ExportDownloadFormatDropdownLabel")
	logrus.Debug("ExportDownloadFormatDropdownLabel > ", core.PrintMatchingInputValues(formatDropdownLabel))
	formatDropdownLabel.Click()

	formatDropdownOption := page.Locator(fmt.Sprintf(pageObjects.ExportDownloadFormatDropdownOption, format))
	formatDropdownOption.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	core.AssertHasMatchingElements(formatDropdownOption, "ExportDownloadFormatDropdownOption")
	logrus.Debug("ExportDownloadFormatDropdownOption > ", core.PrintMatchingInputValues(formatDropdownOption))
	formatDropdownOption.Click()

	logrus.Debug("selected format: ", format)

	// click the download button
	downloadButton := page.Locator(pageObjects.ExportDownloadButton)
	downloadButton.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
	core.AssertHasMatchingElements(downloadButton, "ExportDownloadButton")
	logrus.Debug("downloadButton", core.PrintMatchingInputValues(downloadButton))
	download, err := page.ExpectDownload(func() error {
		logrus.Debug("downloading  transactions")
		return downloadButton.Click()
	})
	core.AssertErrorToNilf("could not download: %w", err)

	filename := download.SuggestedFilename()
	download.SaveAs(filename)
	logrus.Infof("Downloaded %s for %s[%s] %s - %s \n",
		filename,
		accountName,
		accountNumber,
		fromDateString,
		toDateString,
	)

	return filename, nil
}

func NewAnzParsedProcessor(config map[string]interface{}) (*AnzProcessor, error) {
	anzConfig, err := NewAnzConfig(config)
	if err != nil {
		return nil, err
	}

	return NewAnzProcessor(anzConfig), nil
}

func NewAnzProcessor(anzConfig *AnzConfig) *AnzProcessor {
	processor := &AnzProcessor{
		Config: anzConfig,
		Processor: &Processor{
			Name: "anz",
		},
	}

	return processor
}

// AnzPageObjects is a struct that contains the page objects for the ANZ internet banking website.
type AnzPageObjects struct {
	LoginHeader                        string `json:"login_header" yaml:"login_header"`
	LoginUsernameInput                 string `json:"login_username_input" yaml:"login_username_input"`
	LoginPasswordInput                 string `json:"login_password_input" yaml:"login_password_input"`
	LoginButton                        string `json:"login_button" yaml:"login_button"`
	AccountsPageHeader                 string `json:"accounts_page_header" yaml:"accounts_page_header"`
	ExportPageHeader                   string `json:"export_page_header" yaml:"export_page_header"`
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
	LoginHeader:                        "h1#login-header",
	LoginUsernameInput:                 "input[name='customerRegistrationNumber']",
	LoginPasswordInput:                 "input[name='password']",
	LoginButton:                        "button[data-test-id='log-in-btn']",
	AccountsPageHeader:                 "h1[id='home-title']",
	ExportPageHeader:                   "h1[id='search-transaction']",
	ExportAccountDropdownLabel:         "label[for='drop-down-search-transaction-account1-dropdown-field']",
	ExportAccountDropdownOption:        "//ul[@data-test-id='drop-down-search-transaction-account1-dropdown-results']/li[contains(.,'%s')]",
	ExportDateRangeModeButton:          "ul[role='tablist'] li[id='Date rangetab']",
	ExportDateRangeFromDateInput:       "input[id='fromdate-textfield']",
	ExportDateRangeToDateInput:         "input[id='todate-textfield']",
	ExportDownloadFormatDropdownLabel:  "label[data-test-id='drop-down-search-software-dropdown-field-input-text-label']",
	ExportDownloadFormatDropdownOption: "//ul[@data-test-id='drop-down-search-software-dropdown-results']/li[contains(. '%s')]",
	ExportDownloadButton:               "//button[contains(., 'Download')]",
}
