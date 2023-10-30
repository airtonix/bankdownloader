package sources

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAnzSourceLogin(t *testing.T) {

	content := `
		<html>
			<body>
				<form action="/internetbanking">
					<label for="username">Customer Registration Number</label>
					<input id="username" name="username" type="text" />
					<label for="password">Password</label>
					<input id="password" name="password" type="password" />
					<button type="submit">Login</button>
				</form>
			</body>
		</html>
	`

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			path := strings.TrimSpace(request.URL.Path)
			switch path {

			case "/internetbanking":
				t.Log("[MockRequest] /internetbanking")
				w.Write([]byte(strings.ToUpper(content)))

			default:
				w.Write([]byte(""))
			}

		}),
	)
	t.Log("Creating test server")

	defer s.Close()

	source := NewAnzSource(NewSourceParams{
		Domain: s.URL,
	})
	t.Log("Created source")

	source.OpenBrowser()
	t.Log("Opened browser")

	err := source.Login(AnzCredentials{
		Username: "username",
		Password: "password",
	})

	assert.NoError(t, err, "error")
}

func TestAnzSourceDownload(t *testing.T) {

	content := `
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
	`

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			switch strings.TrimSpace(request.URL.Path) {

			case "/IBUI/#/download-transaction":
				w.Write([]byte(strings.ToUpper(content)))

			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}),
	)
	t.Log("Started test server on", s.URL)

	defer s.Close()

	source := NewAnzSource(NewSourceParams{
		Domain: s.URL,
	})
	t.Log("Created source")

	source.OpenBrowser()
	t.Log("Opened browser")

	filename, err := source.DownloadTransactions(
		"My Account",
		"123456789",
		"Agrimaster(CSV)",
		time.Now().Add(-time.Hour*24*30),
		time.Now(),
	)

	assert.NoError(t, err, "error")
	assert.Equal(t, "anz-123456789.csv", filename, "filename")
}
