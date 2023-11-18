package processors

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"time"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func CreatePathLoggingMiddleware(t *testing.T) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf(`
      [MockRequest]
        Path: %s
      `,
				r.URL.Path,
			)
			next.ServeHTTP(w, r)
		})
	}
}

func MockServer(t *testing.T) *httptest.Server {
	r := mux.NewRouter()
	r.Use(CreatePathLoggingMiddleware(t))

	tpl := template.New("root")
	tpl.New("navbar").Parse(`'
  <div data-test-id='navbar-container'>
    <nav>
      <ul>
        <li>
          <a role='button' aria-label='Home' href="/accounts">🏡</a>
        </li>
      </ul>
    </nav>
  </div>
  `)

	tpl.New("login").Parse(`
    <html>
        <body>
            {{ template "navbar" }}
            <h1 id="login-header">Login</h1>
            <form action="/accounts">
                <input id="customerRegistrationNumber" name="customerRegistrationNumber" type="text" />
                <input id="password" name="password" type="password" />
                <button data-test-id="log-in-btn" type="submit">Login</button>
            </form>
        </body>
    </html>
    `)

	tpl.New("accounts").Parse(`
    <html>
        <head>
          <style>
          [onclick] {
            cursor: pointer;
            color: blue;
          }
          </style>
        </head>
        <body>
            {{ template "navbar" }}
            <h1 id="home-title">Accounts</h1>
            <div id="main-div">
                <ul>
                    <li><a id="main-details-wrapper" href="/accounts/123456789">123456789</a></li>
                    <li><a id="main-details-wrapper" href="/accounts/987654321">987654321</a></li>
                </ul>
            </div>

        </body>
    </html>
    `)

	tpl.New("account-detail").Parse(`
    <html>
        <head>
          <style>
            input[type="radio"] {
              display: none;
            }
            [role="tabpanel"] {
              display: none;
            }
            input[type="radio"]:checked + [role="tabpanel"] {
              display: block;
            }
          </style>
        </head>
        <body>
            {{ template "navbar" }}
            <div id="account-overview">
              <h1>Account Overview</h1>
              <div>
                <span>Account Number</span>
                <span>{{.account}}</span>
              </div>
            </div>
            <div id="main-div">
                <ul aria-label="Account Overview" role="tablist">
                    <li role="tab"><label for="Transactionspanelswitch">Transactions</label></li>
                    <li role="tab"><label for="Detailspanelswitch">Details</label></li>
                </ul>

                <input type="radio" id="Detailspanelswitch" checked="checked" name="tabs"/>
                <div id="Detailspanel" role="tabpanel">
                  overview panel we don't want
                </div>

                <input type="radio" id="Transactionspanelswitch" name="tabs" />
                <div id="Transactionspanel" role="tabpanel" >
                  <div id="search-download">
                    <a href="/search-transactions"><span>Search</span></a>
                    <a href="/download-transactions"><span>Download</span></a>
                  </div>
                  <div>transactions panel we want</div>
                </div>
            </div>
        </body>
    </html>
    `)

	tpl.New("download-transactions").Parse(`
    <html>
      <head>
            <style>
              [onclick] {
                cursor: pointer;
                color: blue;
              }

              [role="tappanel"] {
                display: none;
              }
              input[type="radio"][name="tabs"] {
                display: none;
              }
              input[type="radio"]:checked + [role="tappanel"] {
                display: block;
              }
            </style>
      </head>
      <body>
        {{ template "navbar" }}
        <h1 id="search-transaction">Download transactions</h1>
        <form>
          <div>
            <label for="drop-down-search-transaction-account1-dropdown-field">Account</label>
            <input id="drop-down-search-transaction-account1-dropdown-field" name="transaction-account1" />
            <ul data-test-id="drop-down-search-transaction-account1-dropdown-results">
              <li role="option" onClick="document.getElementById('drop-down-search-transaction-account1-dropdown-field').value=this.innerText">123456789</li>
              <li role="option" onClick="document.getElementById('drop-down-search-transaction-account1-dropdown-field').value=this.innerText">987654321</li>
              <li role="option" onClick="document.getElementById('drop-down-search-transaction-account1-dropdown-field').value=this.innerText">543216789</li>
              <li role="option" onClick="document.getElementById('drop-down-search-transaction-account1-dropdown-field').value=this.innerText">678954321</li>
            <ul>
          </div>

          <div>
            <ul aria-label="Search period" role="tablist">
              <li role="tab" aria-controls="Date rangepanel"><label for="Daterangepanelswitch"><div>Date range</div></label></li>
              <li role="tab" aria-controls="Durationpanel"><label for="Durationrangepanelswitch"><div>Duration</div></label></li>
            </ul>
            
            <input type="radio" id="Durationrangepanelswitch" name="tabs" checked="checked" />
            <div id="Durationpanel" role="tappanel">
              Duration panel we dont want
            </div>

            <input type="radio" id="Daterangepanelswitch" name="tabs" />
            <div id="Date rangepanel" role="tappanel">
              <label for="fromdate-textfield">From</label>
              <input id="fromdate-textfield" name="daterange-fromdate" />
              <label for="todate-textfield">To</label>
              <input id="todate-textfield" name="daterange-todate" />
            </div>
          </div>
          
          <div>
            <label for="drop-down-search-software-dropdown-field">Software package</label>
            <input id="drop-down-search-software-dropdown-field" name="format"/>
            <ul data-test-id="drop-down-search-software-dropdown-results">
              <li role="option" onClick="document.getElementById('drop-down-search-software-dropdown-field').value=this.innerText">CSV</li>
              <li role="option" onClick="document.getElementById('drop-down-search-software-dropdown-field').value=this.innerText">Agrimaster(CSV)</li>
              <li role="option" onClick="document.getElementById('drop-down-search-software-dropdown-field').value=this.innerText">MYOB(OFX)</li>
              <li role="option" onClick="document.getElementById('drop-down-search-software-dropdown-field').value=this.innerText">MYOB(QIF)</li>
              <li role="option" onClick="document.getElementById('drop-down-search-software-dropdown-field').value=this.innerText">Quicken(OFX)</li>
              <li role="option" onClick="document.getElementById('drop-down-search-software-dropdown-field').value=this.innerText">Quicken(QIF)</li>
              <li role="option" onClick="document.getElementById('drop-down-search-software-dropdown-field').value=this.innerText">Microsoft Excel(CSV)</li>
            </ul>
          </div>
          <a data-test-id="footer-primary-button_button" download="SomeFile.txt" href="data:text/plain;charset=utf8;,hello world">Download</a>
        </form>
      </body>
    </html>
    `)

	// processors decide what the paths are.
	r.HandleFunc("/internetbanking", func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		tpl.ExecuteTemplate(w, "login", nil)
	})

	r.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		tpl.ExecuteTemplate(w, "accounts", nil)
	})

	r.HandleFunc("/accounts/{account}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.WriteHeader(http.StatusOK)
		tpl.ExecuteTemplate(w, "account-detail", vars)
	})

	r.HandleFunc("/accounts/{account}/transactions", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.WriteHeader(http.StatusOK)
		tpl.ExecuteTemplate(w, "transactions", vars)
	})

	r.HandleFunc("/download-transactions", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		tpl.ExecuteTemplate(w, "download-transactions", nil)
	})
	r.HandleFunc("/search-transactions", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		tpl.ExecuteTemplate(w, "search-transactions", nil)
		t.FailNow()
	})

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	t.Log("Creating test server")
	s := httptest.NewServer(r)
	return s
}

func MakeConfigurations(url string) (store.SourceConfig, store.UsernameAndPassword) {
	sourceConfig := store.SourceConfig{
		Domain:         url,
		ExportFormat:   "Agrimaster(CSV)",
		OutputTemplate: "{{.Account.NameSlug}}-{{.Account.NumberSlug}}.csv",
		DaysToFetch:    30,
	}

	credentials := store.UsernameAndPassword{
		Username: "username",
		Password: "password",
	}
	return sourceConfig, credentials
}

func TestAnzSourceLogin(t *testing.T) {
	var err error

	s := MockServer(t)
	logrus.SetLevel(logrus.DebugLevel)

	automation := core.NewAutomation()

	sourceConfig, credentials := MakeConfigurations(s.URL)

	source := NewAnzProcessor(sourceConfig, credentials, automation)

	t.Log("Created processor")

	err = source.Login()

	assert.NoError(t, err, "error")
	automation.CloseBrowser()
}

func TestAnzSourceDownload(t *testing.T) {
	var err error

	s := MockServer(t)
	defer s.Close()

	automation := core.NewAutomation()

	sourceConfig, credentials := MakeConfigurations(s.URL)

	source := NewAnzProcessor(sourceConfig, credentials, automation)
	// start on accounts page
	automation.Goto(s.URL + "/accounts")
	downloaded, err := source.DownloadTransactions(
		"My Account",
		"123456789",
		time.Now().Add(-time.Hour*24*30),
		time.Now(),
	)

	_, downloadFilename := path.Split(downloaded)

	assert.NoError(t, err, "couldn't download transactions")
	assert.Equal(t, "my-account-123456789.csv",
		// split the path from the filename
		downloadFilename,
		"filename")
}
