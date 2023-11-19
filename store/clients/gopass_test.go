package clients

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGopassClient(t *testing.T) {
	store, err := MockGopassApi([]MockStoredGopassSecret{
		{
			Name:   []string{"pathtosecret"},
			Secret: MockGopassSecret(t, "somepassword\nusername: someguy"),
		},
	})
	assert.NoError(t, err)

	username, err := store.GetUsername("pathtosecret")
	assert.NoError(t, err)
	assert.Equal(t, "someguy", username)

	password, err := store.GetPassword("pathtosecret")
	assert.NoError(t, err)
	assert.Equal(t, "somepassword", password)
}

func TestGopassOtpClient(t *testing.T) {
	when := time.Date(2022, 1, 1, 2, 1, 1, 1, time.UTC)
	totpURL := "otpauth://totp/github-fake-account?secret=rpna55555qyho42j"
	secret := MockGopassSecret(t, "somepassword\nusername: someguy\ntotp: "+totpURL)
	store, err := MockGopassApi([]MockStoredGopassSecret{
		{
			Name:   []string{"pathtosecret"},
			Secret: secret,
		},
	})
	assert.NoError(t, err)

	expectedOtp, err := ResolveOtp(secret, when)
	assert.NoError(t, err)

	retrievedOtp, err := store.GetOtp("pathtosecret", when)
	assert.NoError(t, err)

	assert.Equal(t, expectedOtp, retrievedOtp, "otp should be the same")
}
