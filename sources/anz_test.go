package sources

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnzSourceLogin(t *testing.T) {

	content := `
		<html>
			<body>
				<form>
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
			switch strings.TrimSpace(request.URL.Path) {

			case "/internetbanking":
				w.Write([]byte(strings.ToUpper(content)))

			default:
				w.WriteHeader(http.StatusNotFound)
			}

		}),
	)

	defer s.Close()

	source := NewAnzSource(NewSourceParams{
		Domain: s.URL,
	})

	source.OpenBrowser()

	err := source.Login(AnzCredentials{
		Username: "username",
		Password: "password",
	})

	assert.NoError(t, err, "error")

}
