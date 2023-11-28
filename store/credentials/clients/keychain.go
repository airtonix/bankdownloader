package clients

import (
	"errors"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zalando/go-keyring"
)

type KeychainSecretResolver struct {
}

// ensure that KeychainSecretResolver implements the SecretsResolver interface
var _ SecretsResolver = (*KeychainSecretResolver)(nil)

// Get the service name and secret name from the secret path
func (k *KeychainSecretResolver) getRequestArgs(secretpath string) (string, string) {
	parts := strings.Split(secretpath, "/")
	service := parts[0]
	name := parts[1]
	return service, name
}

// Get the password for a given secret path
func (k *KeychainSecretResolver) GetPassword(secretpath string) (string, error) {
	servicename, secretname := k.getRequestArgs(secretpath)
	logrus.Infof("Getting password for %s/%s", servicename, secretname)
	secret, err := keyring.Get(servicename, secretname)
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
