package credentials

// CredentialsFile is a struct that contains the credentials for a source.
type CredentialsFileSource struct {
	Username string
	Password string
}

// ensure that CredentialsFile implements the ICredentials interface
var _ ICredentialsSource = (*CredentialsFileSource)(nil)

func (c *CredentialsFileSource) Type() CredentialSourceType { return CredentialSourceTypeFile }
func (c *CredentialsFileSource) Resolve() ResolvedCredentials {
	return ResolvedCredentials{
		UsernameAndPassword: UsernameAndPassword{
			Username: c.Username,
			Password: c.Password,
		},
		Type: c.Type(),
	}
}
