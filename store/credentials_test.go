package store

import (
	"testing"

	"github.com/airtonix/bank-downloaders/store/clients"
	"github.com/stretchr/testify/assert"
)

func TestGopassCredentials(t *testing.T) {
	expectedUsername := "someguy"
	expectedPassword := "somepassword"

	client := clients.NewMockGopassResolver(
		map[string]interface{}{
			"websites/test.com/someguy": map[string]interface{}{
				"username": expectedUsername,
				"password": expectedPassword,
			},
		},
	)

	credentials := &CredentialsGopassSource{
		Secret:      "websites/test.com/someguy",
		UsernameKey: "username",
		PasswordKey: "password",
		Api:         client,
	}

	resolved, err := credentials.Resolve()
	if err != nil {
		t.Errorf("failed to resolve credentials: %v", err)
	}

	assert.Equal(t, expectedUsername, resolved.UsernameAndPassword.Username)
	assert.Equal(t, expectedPassword, resolved.UsernameAndPassword.Password)
}
func TestGopassTotpCredentials(t *testing.T) {
	expectedUsername := "someguy"
	expectedPassword := "somepassword"

	client := clients.NewMockGopassResolver(
		map[string]interface{}{
			"websites/test.com/someguy": map[string]interface{}{
				"username": expectedUsername,
				"password": expectedPassword,
				"totp":     "123456",
			},
		},
	)

	credentials := &CredentialsGopassTotpSource{
		Secret:      "websites/test.com/someguy",
		UsernameKey: "username",
		PasswordKey: "password",
		TotpKey:     "totp",
		Api:         client,
	}

	resolved, err := credentials.Resolve()
	if err != nil {
		t.Errorf("failed to resolve credentials: %v", err)
	}

	assert.Equal(t, expectedUsername, resolved.UsernameAndPasswordAndTotp.Username)
	assert.Equal(t, expectedPassword, resolved.UsernameAndPasswordAndTotp.Password)
	assert.Len(t, resolved.UsernameAndPasswordAndTotp.Totp, 6)

}
