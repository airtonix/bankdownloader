package clients

import "github.com/zalando/go-keyring"

type MockStoredKeychainSecret struct {
	Name   string
	Secret string
}

func MockKeychainResolver(
	serviceName string,
	secrets []MockStoredKeychainSecret,
) *KeychainSecretResolver {
	// sets up a mock keychain client
	keyring.MockInit()

	for _, secret := range secrets {
		keyring.Set(serviceName, secret.Name, secret.Secret)
	}

	return NewKeychainResolver()
}
