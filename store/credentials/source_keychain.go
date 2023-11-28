package credentials

import (
	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/store/credentials/clients"
)

// CredentialsKeychain is a struct that contains the credentials for a source.
type CredentialsKeychainSource struct {
	ServiceName string
	Username    string
	Api         *clients.KeychainSecretResolver
}

// ensure that CredentialsKeychain implements the ICredentials interface
var _ ICredentialsSource = (*CredentialsKeychainSource)(nil)

func (c *CredentialsKeychainSource) Type() CredentialSourceType { return CredentialSourceTypeKeychain }
func (c *CredentialsKeychainSource) Resolve() ResolvedCredentials {
	secretpath := c.ServiceName + "/" + c.Username

	username, err := c.Api.GetUsername(secretpath)
	core.AssertErrorToNilf("failed to get username for: %s", err)

	password, err := c.Api.GetPassword(secretpath)
	core.AssertErrorToNilf("failed to get password for: %s", err)

	return ResolvedCredentials{
		UsernameAndPassword: UsernameAndPassword{
			Username: username,
			Password: password,
		},
		Type: c.Type(),
	}
}
