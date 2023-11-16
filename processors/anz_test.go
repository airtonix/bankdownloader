package processors

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/store"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func server(t *testing.T) *httptest.Server {
	r := mux.NewRouter()

	var tpl = template.Template{}

	tpl.New("login").Parse(`
		<html>
		<body>
		<h1 id="login-header">Login</h1>
		<form action="/accounts">
		<input id="customerRegistrationNumber" name="customerRegistrationNumber" type="text" />
		<input id="password" name="password" type="password" />
		<button data-test-id="log-in-btn" type="submit">Login</button>
		</form>
		</body>
		</html>
		`)

	tpl.New("accounts-list").Parse(`
		<html>
		<body>
		<h1 id="home-title">Accounts</h1>
		<div id="main-div">
		<ul>
		<li>
		<a href="/accounts/123456789">123456789</a>
		</li>
		<li>
		<a href="/accounts/987654321">987654321</a>
		</li>
		</ul>
		</body>
		</html>
		`)

	tpl.New("account-detail").Parse(`
		<html>
		<body>
		<h1 id="home-title">Account 123456789</h1>
		<div id="main-div">
		<ul aria-label="Account Overview">
		<li id="Transactionstab">Transactions</li>
		</ul>
		</div>
		</body>
		</html>
		`)

	tpl.New("transactions").Parse(`
		<html>
		<body>
		<h1 id="home-title">Account {{.account}}</h1>
		<div id="main-div">
		<ul aria-label="Account Overview">
		<li id="Transactionstab">Transactions</li>
		</ul>
		</div>
		</body>
		</html>
		`)

	tpl.New("download").Parse(`
		<html>
			<body>
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
							<li role="tab">Date range</li>
							<li role="tab">Duration</li>
						</ul>
						<div id="Date rangepanel" role="tappanel">
							<label for="fromdate-textfield">From</label>
							<input id="fromdate-textfield" name="daterange-fromdate" />
							<label for="todate-textfield">To</label>
							<input id="todate-textfield" name="daterange-todate" />
						</div>
					</div>
					
					<div>
						<label for="drop-down-search-software-dropdown-field" data-test-id="drop-down-search-transaction-account1-dropdown-field-input-text-label">Software package</label>
						<input id="drop-down-search-software-dropdown-field" name="format" />
						<ul data-test-id="drop-down-search-software-dropdown-field">
							<li role="option" onClick="document.getElementById('drop-down-search-software-dropdown-field').value=this.innerText">CSV</li>
							<li role="option" onClick="document.getElementById('drop-down-search-software-dropdown-field').value=this.innerText">Agrimaster(CSV)</li>
							<li role="option" onClick="document.getElementById('drop-down-search-software-dropdown-field').value=this.innerText">MYOB(OFX)</li>
							<li role="option" onClick="document.getElementById('drop-down-search-software-dropdown-field').value=this.innerText">MYOB(QIF)</li>
							<li role="option" onClick="document.getElementById('drop-down-search-software-dropdown-field').value=this.innerText">Quicken(OFX)</li>
							<li role="option" onClick="document.getElementById('drop-down-search-software-dropdown-field').value=this.innerText">Quicken(QIF)</li>
							<li role="option" onClick="document.getElementById('drop-down-search-software-dropdown-field').value=this.innerText">Microsoft Excel(CSV)</li>
						</ul>
					</div>

					<button type="submit">Download</button>
				</form>
			</body>
		</html>
		`)

	r.HandleFunc("/internetbanking", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		t.Logf("[MockRequest] %s", path)
		tpl.ExecuteTemplate(w, "login", nil)
	})

	r.HandleFunc("/IBUI/#/home", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		t.Logf("[MockRequest] %s", path)
		tpl.ExecuteTemplate(w, "accounts-list", nil)
	})

	r.HandleFunc("/IBUI/#account/{account}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := r.URL.Query().Get("path")
		t.Logf("[MockRequest] %s", path)
		tpl.ExecuteTemplate(w, "account-detail", vars)
	})

	r.HandleFunc("/IBUI/#/account/{account}/transactions", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := r.URL.Query().Get("path")
		t.Logf("[MockRequest] %s", path)
		tpl.ExecuteTemplate(w, "transactions", vars)
	})

	r.HandleFunc("/IBUI/#/download-transaction", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		t.Logf("[MockRequest] %s", path)
		tpl.ExecuteTemplate(w, "download", nil)
	})

	r.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	t.Log("Creating test server")
	s := httptest.NewServer(r)
	return s
}

func TestAnzSourceLogin(t *testing.T) {
	var err error

	s := server(t)

	chromedp.Flag("headless", false)
	automation := core.NewAutomation()

	sourceConfig := store.SourceConfig{
		Domain:         s.URL,
		ExportFormat:   "csv",
		OutputTemplate: "{{.AccountName}}-{{.AccountNumber}}.csv",
		DaysToFetch:    30,
	}
	credentials := store.UsernameAndPassword{
		Username: "username",
		Password: "password",
	}
	source := NewAnzProcessor(sourceConfig, credentials, automation)

	t.Log("Created processor")

	err = source.Login()

	assert.NoError(t, err, "error")
}

func TestAnzSourceDownload(t *testing.T) {
	var err error

	s := server(t)
	defer s.Close()

	automation := core.NewAutomation()

	sourceConfig := store.SourceConfig{
		Domain:         s.URL,
		ExportFormat:   "csv",
		OutputTemplate: "{{.AccountName}}-{{.AccountNumber}}.csv",
		DaysToFetch:    30,
	}
	credentials := store.UsernameAndPassword{
		Username: "username",
		Password: "password",
	}
	source := NewAnzProcessor(sourceConfig, credentials, automation)
	filename, err := source.DownloadTransactions(
		"My Account",
		"123456789",
		time.Now().Add(-time.Hour*24*30),
		time.Now(),
	)

	assert.NoError(t, err, "couldn't download transactions")
	assert.Equal(t, "anz-123456789.csv", filename, "filename")
}
