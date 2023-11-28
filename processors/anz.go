package processors

import (
	"fmt"
	"time"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/store"
	"github.com/airtonix/bank-downloaders/store/credentials"
	"github.com/sirupsen/logrus"
)

type AnzProcessor struct {
	Credentials credentials.UsernameAndPassword
	store.SourceConfig
	Processor
	Automation *core.Automation
}

// ensure that AnzProcessor implements the Processor interface
var _ IProcessor = (*AnzProcessor)(nil)

func (processor *AnzProcessor) Login() error {
	loginDetails := processor.Credentials
	url := fmt.Sprintf(
		"%s/internetbanking",
		processor.SourceConfig.Domain,
	)

	automation := processor.Automation

	logrus.Info("logging into ", url)

	// start at the login page
	automation.Goto(url)
	automation.SetViewportSize(1200, 900)

	logrus.Debugln("waiting for login page to load...")
	// wait for the login page to load
	automation.Find(pageObjects.LoginHeader)

	// Username
	automation.Find(pageObjects.LoginUsernameInput)
	automation.Focus(pageObjects.LoginUsernameInput)
	automation.Fill(pageObjects.LoginUsernameInput, loginDetails.Username)

	// Password
	automation.Find(pageObjects.LoginPasswordInput)
	automation.Focus(pageObjects.LoginPasswordInput)
	automation.FillSensitive(pageObjects.LoginPasswordInput, loginDetails.Password)

	// LoginButton
	automation.Click(pageObjects.LoginButton)

	logrus.Info("authenticating...")

	// Accounts Page
	// wait for the account page to load
	automation.Find(pageObjects.AccountsPageHeader)
	logrus.Info("authenticated")

	return nil
}

func (processor *AnzProcessor) DownloadTransactions(
	accountName string,
	accountNumber string,
	fromDate time.Time,
	toDate time.Time,
) (string, error) {
	automation := processor.Automation
	dateFormat := "02/01/2006"

	var format = processor.SourceConfig.ExportFormat
	fromDateString := fromDate.Format(dateFormat)
	toDateString := toDate.Format(dateFormat)

	// ANZ web app uses the transactions page for two purposes: searching and downloading.
	// it's only possible to be in one mode or the other as a result of clicking the right button.
	// As such, when we want to download transactions for an account, we first need to go to the
	// home page, then click the account button, then click the download button.

	logrus.Infoln(
		fmt.Sprintf(
			"Fetching transactions for: %s [%s]: %s - %s",
			accountName,
			accountNumber,
			fromDateString,
			toDateString,
		),
	)
	automation.Find(pageObjects.NavigateToHomeButton)
	automation.Click(pageObjects.NavigateToHomeButton)
	// ANZ web app uses responsive design, so we need to set the viewport size
	// otherwise we get a different set of selectors (we use the desktop version)
	automation.SetViewportSize(1200, 900)

	// find the account button
	automation.Click(fmt.Sprintf(pageObjects.AccountsListAccountButton, accountNumber))

	automation.Find(fmt.Sprintf(pageObjects.AccountDetailHeader, accountNumber))
	automation.Find(pageObjects.AccountTransactionTabButton)
	// click the transaction tab button
	automation.Click(pageObjects.AccountTransactionTabButton)

	// find the account button
	automation.Click(pageObjects.AccountGotoExportButton)

	// Transactions Page
	// wait for the page to load
	automation.Find(pageObjects.ExportPageHeader)

	// pick the account by clicking the label "Account"
	automation.Click(pageObjects.ExportAccountDropdownLabel)
	// then click the account option
	automation.Click(fmt.Sprintf(pageObjects.ExportAccountDropdownOption, accountNumber))
	logrus.Debug("selected account: ", accountNumber)

	// change to date range mode
	automation.Click(pageObjects.ExportDateRangeModeButton)
	// select the date range fromDate
	automation.Fill(pageObjects.ExportDateRangeFromDateInput, fromDateString)
	// select the date range toDate
	automation.Fill(pageObjects.ExportDateRangeToDateInput, toDateString)
	logrus.Debugf(
		"selected date range: %s - %s",
		fromDateString, toDateString,
	)

	// select the downlaod format by clicking the label "Software package"
	automation.Click(pageObjects.ExportDownloadFormatDropdownLabel)
	// select the download format
	automation.Click(fmt.Sprintf(pageObjects.ExportDownloadFormatDropdownOption, format))
	logrus.Debug("selected format: ", format)

	filenameContext := store.NewFilenameTemplateContext(
		processor.Name,
		accountName,
		accountNumber,
		fromDate,
		toDate,
	)

	filenameTemplate := store.NewFilenameTemplate(processor.OutputTemplate)

	// click the download button
	filename, err := automation.DownloadFile(
		filenameTemplate.Render(filenameContext),
		func() error {
			automation.Click(pageObjects.ExportDownloadButton)
			// if we get this far it didn't panic
			return nil
		},
	)
	core.AssertErrorToNilf(
		fmt.Sprintf("could not download file: %s", filename),
		err)

	logrus.Info("Downloaded", filename)

	return filename, nil
}

func NewAnzProcessor(
	config store.SourceConfig,
	credentials credentials.UsernameAndPassword,
	automation *core.Automation,
) *AnzProcessor {
	processor := Processor{
		Name: "anz",
	}

	return &AnzProcessor{
		Processor:    processor,
		SourceConfig: config,
		Automation:   automation,
		Credentials:  credentials,
	}
}

// AnzPageObjects is a struct that contains the page objects for the ANZ internet banking website.
type AnzPageObjects struct {
	LoginHeader                        string
	LoginUsernameInput                 string
	LoginPasswordInput                 string
	LoginButton                        string
	NavigateToHomeButton               string
	AccountsPageHeader                 string
	AccountsListAccountButton          string
	AccountTransactionTabButton        string
	AccountDetailHeader                string
	AccountGotoExportButton            string
	ExportPageHeader                   string
	ExportAccountDropdownLabel         string
	ExportAccountDropdownOption        string
	ExportDateRangeModeButton          string
	ExportDateRangeFromDateInput       string
	ExportDateRangeToDateInput         string
	ExportDownloadFormatDropdownLabel  string
	ExportDownloadFormatDropdownOption string
	ExportDownloadButton               string
}

var pageObjects = AnzPageObjects{
	LoginHeader:                        "h1#login-header",
	LoginUsernameInput:                 "input[name='customerRegistrationNumber']",
	LoginPasswordInput:                 "input[name='password']",
	LoginButton:                        "button[data-test-id='log-in-btn']",
	NavigateToHomeButton:               "div[data-test-id='navbar-container'] [role='button'][aria-label='Home']",
	AccountsPageHeader:                 "h1[id='home-title']",
	AccountsListAccountButton:          "//div[@id='main-div'] //*[@id='main-details-wrapper'][contains(., '%s')]",
	AccountDetailHeader:                "//div[@id='account-overview'][contains(., '%s')]",
	AccountTransactionTabButton:        "//ul[@role='tablist'][@aria-label='Account Overview'] //li[@role='tab'] //*[contains(., 'Transactions')]",
	AccountGotoExportButton:            "//div[@id='search-download'] //span[contains(., 'Download')]",
	ExportPageHeader:                   "//h1[@id='search-transaction'][contains(., 'Download transactions')]",
	ExportAccountDropdownLabel:         "label[for='drop-down-search-transaction-account1-dropdown-field']",
	ExportAccountDropdownOption:        "//ul[@data-test-id='drop-down-search-transaction-account1-dropdown-results']/li[contains(.,'%s')]",
	ExportDateRangeModeButton:          "//ul[@role='tablist'] //li[@aria-controls='Date rangepanel'] //div[contains(., 'Date range')]",
	ExportDateRangeFromDateInput:       "input[id='fromdate-textfield']",
	ExportDateRangeToDateInput:         "input[id='todate-textfield']",
	ExportDownloadFormatDropdownLabel:  "//label[@for='drop-down-search-software-dropdown-field'][contains(., 'Software package')]",
	ExportDownloadFormatDropdownOption: "//ul[@data-test-id='drop-down-search-software-dropdown-results']/li[@role='option'][contains(., '%s')]",
	ExportDownloadButton:               "//*[@data-test-id='footer-primary-button_button'][contains(., 'Download')]",
}
