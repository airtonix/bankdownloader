package clients

import (
	"github.com/sirupsen/logrus"
	"github.com/zalando/go-keyring"
)

type MockStoredKeychainSecret struct {
	// secretname/username
	Name string
	// username
	Username string
	// password
	Password string
}

func NewMockKeychainSecretResolver(
	secrets []MockStoredKeychainSecret,
) *KeychainSecretResolver {
	// sets up a mock keychain client
	keyring.MockInit()

	for _, secret := range secrets {
		keyring.Set(secret.Name, secret.Username, secret.Password)
		logrus.Infof("keychain: %s, username: %s, password: %s", secret.Name, secret.Username, secret.Password)
	}

	return NewKeychainResolver()
}
