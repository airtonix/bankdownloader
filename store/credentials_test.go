package store

import (
	"testing"
	"time"

	"github.com/airtonix/bank-downloaders/store/clients"
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
		Secret:      "pathtosecret",
		UsernameKey: "username",
		PasswordKey: "password",
		Api:         store,
	}

	resolved, err := credentials.Resolve()
	assert.NoError(t, err)

	assert.Equal(t, "someguy", resolved.UsernameAndPassword.Username)
	assert.Equal(t, "somepassword", resolved.UsernameAndPassword.Password)
}

func TestGopassTotpCredentials(t *testing.T) {
	when := time.Date(2022, 1, 1, 2, 1, 1, 1, time.UTC)
	totpURL := "otpauth://totp/github-fake-account?secret=rpna55555qyho42j"
	totpSecret := clients.NewMockGopassSecret(t, "somepassword\nusername: someguy\ntotp: "+totpURL)
	expectedToken, err := clients.ResolveOtp(totpSecret, when)
	assert.NoError(t, err)

	store, err := clients.NewMockGopassSecretResolver([]clients.MockStoredGopassSecret{
		{
			Name:   []string{"pathtosecret"},
			Secret: totpSecret,
		},
	})
	assert.NoError(t, err)

	credentials := &CredentialsGopassTotpSource{
		Secret:      "pathtosecret",
		UsernameKey: "username",
		PasswordKey: "password",
		TotpKey:     "totp",
		Api:         store,
	}
	credentials.SetTimestampFn(func() time.Time {
		return when
	})

	resolved, err := credentials.Resolve()
	assert.NoError(t, err)

	assert.Equal(t, "someguy", resolved.UsernameAndPasswordAndTotp.Username)
	assert.Equal(t, "somepassword", resolved.UsernameAndPasswordAndTotp.Password)
	assert.Equal(t, expectedToken, resolved.UsernameAndPasswordAndTotp.Totp)
}

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

	resolved, err := credentials.Resolve()
	assert.NoError(t, err)

	assert.Equal(t, "someguy", resolved.UsernameAndPassword.Username)
	assert.Equal(t, "somepassword", resolved.UsernameAndPassword.Password)
}
