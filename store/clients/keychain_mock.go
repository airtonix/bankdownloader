package clients

import (
	"path"

	"github.com/zalando/go-keyring"
)

type MockStoredKeychainSecret struct {
	// secretname/username
	Name string
	// password
	Secret string
}

func NewMockKeychainSecretResolver(
	serviceName string,
	secrets []MockStoredKeychainSecret,
) *KeychainSecretResolver {
	// sets up a mock keychain client
	keyring.MockInit()

	for _, secret := range secrets {
		name, username := path.Split(secret.Name)
		keyring.Set(name, username, secret.Secret)
	}

	return NewKeychainResolver()
}
