package credentials

import (
	"testing"

	"github.com/airtonix/bank-downloaders/store/credentials/clients"
	"github.com/stretchr/testify/assert"
)

func TestGopassCredentials(t *testing.T) {
	store, err := clients.NewMockGopassSecretResolver([]clients.MockStoredGopassSecret{
		{
			Name:   []string{"pathtosecret"},
			Secret: clients.NewMockGopassSecret(t, "somepassword\nusername: someguy"),
		},
	})
	assert.NoError(t, err)

	credentials := &CredentialsGopassSource{
		CredentialsGopassInput: CredentialsGopassInput{
			Secret:      "pathtosecret",
			UsernameKey: "username",
			PasswordKey: "password",
		},
		Api: store,
	}

	resolved := credentials.Resolve()
	assert.NoError(t, err)

	assert.Equal(t, "someguy", resolved.UsernameAndPassword.Username)
	assert.Equal(t, "somepassword", resolved.UsernameAndPassword.Password)
}
