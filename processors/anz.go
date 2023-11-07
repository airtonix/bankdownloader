package processors

import (
	"fmt"
	"time"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/store"
	"github.com/kr/pretty"
	"github.com/sirupsen/logrus"
)

type AnzProcessor struct {
	Config *AnzConfig
	*Processor
}

// ensure that AnzProcessor implements the Processor interface
var _ IProcessor = (*AnzProcessor)(nil)

func (processor *AnzProcessor) GetFormat() string {
	return processor.Config.ExportFormat
}
func (processor *AnzProcessor) GetOutputTemplate() string {
	return processor.Config.OutputTemplate
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
	automation.Find(pageObjects.LoginHeader)

	// Username
	automation.Focus(pageObjects.LoginUsernameInput)
	automation.Fill(pageObjects.LoginUsernameInput, loginDetails.Username)

	// Password
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
	automation *core.Automation,
) (string, error) {
	var err error
	page := automation.Page
	dateFormat := "02/01/2006"

	var format = processor.Config.ExportFormat
	fromDateString := fromDate.Format(dateFormat)
	toDateString := toDate.Format(dateFormat)

	// get the current hostname for the current page
	pageUrl := automation.GetPageUrlObject()

	// ANZ web app uses the transactions page for two purposes: searching and downloading.
	// it's only possible to be in one mode or the other as a result of clicking the right button.
	// As such, when we want to download transactions for an account, we first need to go to the
	// home page, then click the account button, then click the download button.
	url := fmt.Sprintf(
		"%s://%s/IBUI/#/home",
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
	// ANZ web app uses responsive design, so we need to set the viewport size
	// otherwise we get a different set of selectors (we use the desktop version)
	page.SetViewportSize(1200, 900)

	// find the account button
	automation.Click(fmt.Sprintf(pageObjects.AccountsListAccountButton, accountNumber))

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

	filenameTemplate := store.NewFilenameTemplate(processor.Config.OutputTemplate)

	// click the download button
	filename, err := automation.DownloadFile(
		filenameTemplate.Render(filenameContext),
		func() error {
			return automation.Click(pageObjects.ExportDownloadButton)
		},
	)

	logrus.Info("Downloaded", filename)

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
	AccountsListAccountButton          string `json:"accounts_list_account_button" yaml:"accounts_list_account_button"`
	AccountGotoExportButton            string `json:"account_goto_export_button" yaml:"account_goto_export_button"`
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
	AccountsListAccountButton:          "//div[@id='main-div'] //li[contains(., '%s')]",
	AccountGotoExportButton:            "//div[@id='search-download'] //button[contains(., 'Download')]",
	ExportPageHeader:                   "h1[id='search-transaction']",
	ExportAccountDropdownLabel:         "label[for='drop-down-search-transaction-account1-dropdown-field']",
	ExportAccountDropdownOption:        "//ul[@data-test-id='drop-down-search-transaction-account1-dropdown-results']/li[contains(.,'%s')]",
	ExportDateRangeModeButton:          "ul[role='tablist'] li[id='Date rangetab']",
	ExportDateRangeFromDateInput:       "input[id='fromdate-textfield']",
	ExportDateRangeToDateInput:         "input[id='todate-textfield']",
	ExportDownloadFormatDropdownLabel:  "//label[@for='drop-down-search-software-dropdown-field'][contains(., 'Software package')]",
	ExportDownloadFormatDropdownOption: "//ul[@data-test-id='drop-down-search-software-dropdown-results']/li[contains(., '%s')]",
	ExportDownloadButton:               "//button[contains(., 'Download')]",
}
