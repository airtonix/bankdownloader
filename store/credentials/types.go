package credentials

type UsernameAndPassword struct {
	Username string
	Password string
}

func (c *UsernameAndPassword) ConfirmResolved() bool {
	return c.Username != "" && c.Password != ""
}

type UsernameAndPasswordAndTotp struct {
	Username string
	Password string
	Totp     string
}

func (c *UsernameAndPasswordAndTotp) ConfirmResolved() bool {
	return c.Username != "" && c.Password != "" && c.Totp != ""
}

type CredentialSourceType string

const (
	// credentials are provided in the config file
	CredentialSourceTypeFile CredentialSourceType = "file"
	// credentials are provided in the environment
	CredentialSourceTypeEnv CredentialSourceType = "env"
	// credentials come from gopass
	CredentialSourceTypeGopass CredentialSourceType = "gopass"
	// credentials come from gopass and have a totp
	CredentialSourceTypeGopassTotp CredentialSourceType = "gopass-totp"
	// credentials come from the keychain
	CredentialSourceTypeKeychain CredentialSourceType = "libsecret"
	// credentials come from the keychain and have a totp
	CredentialSourceTypeKeychainTotp CredentialSourceType = "libsecret-totp"
)

// A fat union
type CredentialsSource struct {
	CredentialsFileSource
	CredentialsEnvSource
	CredentialsGopassSource
	CredentialsGopassTotpSource
	CredentialsKeychainSource
	Type CredentialSourceType
}

// ICredentials is an interface that contains the credentials for a source.
type ICredentialsSource interface {
	Type() CredentialSourceType
	Resolve() ResolvedCredentials
}

type ResolvedCredentials struct {
	UsernameAndPassword
	UsernameAndPasswordAndTotp
	Type CredentialSourceType
}

func (c *ResolvedCredentials) ConfirmResolved() bool {
	switch c.Type {
	case CredentialSourceTypeGopassTotp:
		return c.UsernameAndPasswordAndTotp.ConfirmResolved()
	case CredentialSourceTypeGopass:
		return c.UsernameAndPassword.ConfirmResolved()
	case CredentialSourceTypeFile:
		return c.UsernameAndPassword.ConfirmResolved()
	case CredentialSourceTypeEnv:
		return c.UsernameAndPassword.ConfirmResolved()
	case CredentialSourceTypeKeychainTotp:
		return c.UsernameAndPasswordAndTotp.ConfirmResolved()
	case CredentialSourceTypeKeychain:
		return c.UsernameAndPassword.ConfirmResolved()
	default:
		return false
	}
}
