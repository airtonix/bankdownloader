package credentials

import (
	"testing"

	"github.com/airtonix/bank-downloaders/store/credentials/clients"
	"github.com/stretchr/testify/assert"
)

func TestKeychainCredentials(t *testing.T) {
	store := clients.NewMockKeychainSecretResolver([]clients.MockStoredKeychainSecret{
		{
			Name:     "pathtosecret",
			Username: "someguy",
			Password: "somepassword",
		},
	})

	credentials := &CredentialsKeychainSource{
		ServiceName: "pathtosecret",
		Username:    "someguy",
		Api:         store,
	}

	resolved := credentials.Resolve()

	assert.Equal(t, "someguy", resolved.UsernameAndPassword.Username)
	assert.Equal(t, "somepassword", resolved.UsernameAndPassword.Password)
}
