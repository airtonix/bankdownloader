package clients

import (
	"errors"
	"path"
	"time"

	"github.com/zalando/go-keyring"
)

type KeychainSecretResolver struct {
}

// ensure that KeychainSecretResolver implements the SecretsResolver interface
var _ SecretsResolver = (*KeychainSecretResolver)(nil)

// a keychan secret path is composed of `{secretname}/{username for secret}`
func (k *KeychainSecretResolver) getRequestArgs(secretpath string) (string, string) {
	service, name := path.Split(secretpath)
	return service, name
}

func (k *KeychainSecretResolver) get(secretpath string) (string, error) {
	service, name := k.getRequestArgs(secretpath)
	secret, err := keyring.Get(service, name)
	if err != nil {
		return "", err
	}

	return secret, nil
}

// Get the password for a given secret path
func (k *KeychainSecretResolver) GetPassword(secretpath string) (string, error) {
	secret, err := k.get(secretpath)

	if err != nil {
		return "", err
	}

	return secret, nil
}

// Get the username for a given secret path
// However, in a keychain the username is the secret path
func (k *KeychainSecretResolver) GetUsername(secretpath string) (string, error) {
	_, name := k.getRequestArgs(secretpath)
	return name, nil
}

// Get the OTP code. This is currently not implemented.
// since not sure how to implement this in the context of a secret being a
// record holding all the credentials for a given service.
func (k *KeychainSecretResolver) GetOtp(name string, timestamp time.Time) (string, error) {
	return "", errors.New("not implemented")
}

func NewKeychainResolver() *KeychainSecretResolver {
	return &KeychainSecretResolver{}
}
