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

type KeychainClientGetOptions struct {
	version string
	path    string
}

// a keychan secret path is composed of `{servicename}/{secretname}`
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

func (k *KeychainSecretResolver) GetPassword(secretpath string) (string, error) {
	secret, err := k.get(secretpath)

	if err != nil {
		return "", err
	}

	return secret, nil
}

// in a keychain the username is the entry name
func (k *KeychainSecretResolver) GetUsername(secretpath string) (string, error) {
	_, name := k.getRequestArgs(secretpath)
	return name, nil
}

// not implemented. since not sure how to implement this
func (k *KeychainSecretResolver) GetOtp(name string, timestamp time.Time) (string, error) {
	return "", errors.New("not implemented")
}

func NewKeychainResolver() *KeychainSecretResolver {
	return &KeychainSecretResolver{}
}
