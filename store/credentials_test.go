package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGopassCredentials(t *testing.T) {
	credentials := &CredentialsGopassSource{
		Path:        "websites/test.com/someguy",
		UsernameKey: "username",
		PasswordKey: "password",
	}

	resolved, err := credentials.Resolve()
	if err != nil {
		t.Errorf("failed to resolve credentials: %v", err)
	}

	assert.NotEmpty(t, resolved.UsernameAndPassword.Username)
	assert.NotEmpty(t, resolved.UsernameAndPassword.Password)
}
func TestGopassTotpCredentials(t *testing.T) {
	credentials := &CredentialsGopassTotpSource{
		Path:        "websites/test.com/someguy",
		UsernameKey: "username",
		PasswordKey: "password",
		TotpKey:     "totp",
	}

	resolved, err := credentials.Resolve()
	if err != nil {
		t.Errorf("failed to resolve credentials: %v", err)
	}

	assert.NotEmpty(t, resolved.UsernameAndPasswordAndTotp.Username)
	assert.NotEmpty(t, resolved.UsernameAndPasswordAndTotp.Password)
	assert.NotEmpty(t, resolved.UsernameAndPasswordAndTotp.Totp)

}
