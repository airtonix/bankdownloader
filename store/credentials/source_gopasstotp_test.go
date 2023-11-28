package credentials

import (
	"testing"
	"time"

	"github.com/airtonix/bank-downloaders/store/credentials/clients"
	"github.com/stretchr/testify/assert"
)

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

	resolved := credentials.Resolve()

	assert.Equal(t, "someguy", resolved.UsernameAndPasswordAndTotp.Username)
	assert.Equal(t, "somepassword", resolved.UsernameAndPasswordAndTotp.Password)
	assert.Equal(t, expectedToken, resolved.UsernameAndPasswordAndTotp.Totp)
}
