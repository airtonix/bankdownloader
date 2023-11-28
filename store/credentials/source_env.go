package credentials

import "os"

// CredentialsEnv is a struct that contains the credentials for a source.
type CredentialsEnvSource struct {
	UsernameKey string
	PasswordKey string
}

// ensure that CredentialsEnv implements the ICredentials interface
var _ ICredentialsSource = (*CredentialsEnvSource)(nil)

func (c *CredentialsEnvSource) Type() CredentialSourceType { return CredentialSourceTypeEnv }
func (c *CredentialsEnvSource) Resolve() ResolvedCredentials {
	return ResolvedCredentials{
		UsernameAndPassword: UsernameAndPassword{
			Username: os.Getenv(c.UsernameKey),
			Password: os.Getenv(c.PasswordKey),
		},
		Type: c.Type(),
	}
}
