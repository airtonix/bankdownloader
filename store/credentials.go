package store

import (
	"fmt"
	"os"
	"time"

	"github.com/airtonix/bank-downloaders/store/clients"
)

// Load credentials from various sources

type CredentialSourceType string

type UsernameAndPassword struct {
	Username string
	Password string
}

type UsernameAndPasswordAndTotp struct {
	Username string
	Password string
	Totp     string
}

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

// CredentialsFile is a struct that contains the credentials for a source.
type CredentialsFileSource struct {
	Username string
	Password string
}

// ensure that CredentialsFile implements the ICredentials interface
var _ ICredentialsSource = (*CredentialsFileSource)(nil)

func (c *CredentialsFileSource) Type() CredentialSourceType { return CredentialSourceTypeFile }
func (c *CredentialsFileSource) Resolve() (ResolvedCredentials, error) {
	return ResolvedCredentials{
		UsernameAndPassword: UsernameAndPassword{
			Username: c.Username,
			Password: c.Password,
		},
		Type: c.Type(),
	}, nil
}

// CredentialsEnv is a struct that contains the credentials for a source.
type CredentialsEnvSource struct {
	UsernameKey string
	PasswordKey string
}

// ensure that CredentialsEnv implements the ICredentials interface
var _ ICredentialsSource = (*CredentialsEnvSource)(nil)

func (c *CredentialsEnvSource) Type() CredentialSourceType { return CredentialSourceTypeEnv }
func (c *CredentialsEnvSource) Resolve() (ResolvedCredentials, error) {
	return ResolvedCredentials{
		UsernameAndPassword: UsernameAndPassword{
			Username: os.Getenv(c.UsernameKey),
			Password: os.Getenv(c.PasswordKey),
		},
		Type: c.Type(),
	}, nil
}

// CredentialsGopass is a struct that contains the credentials for a source.
type CredentialsGopassSource struct {
	Secret      string
	UsernameKey string
	PasswordKey string
	Api         clients.GopassSecretResolver
}

// ensure that CredentialsGopass implements the ICredentials interface
var _ ICredentialsSource = (*CredentialsGopassSource)(nil)

func (c *CredentialsGopassSource) Type() CredentialSourceType { return CredentialSourceTypeGopass }
func (c *CredentialsGopassSource) Resolve() (ResolvedCredentials, error) {
	Password, err := c.Api.GetPassword(c.Secret)
	if err != nil {
		return ResolvedCredentials{}, fmt.Errorf("failed to get password for: %s", c.Secret)
	}
	Username, err := c.Api.GetUsername(c.Secret)
	if err != nil {
		return ResolvedCredentials{}, fmt.Errorf("failed to get username for: %s", c.Secret)
	}

	return ResolvedCredentials{
		UsernameAndPassword: UsernameAndPassword{
			Username,
			Password,
		},
		Type: c.Type(),
	}, nil
}

// CredentialsGopassTotp is a struct that contains the credentials for a source.
type CredentialsGopassTotpSource struct {
	Secret      string
	UsernameKey string
	PasswordKey string
	TotpKey     string
	Api         clients.GopassSecretResolver
	timestampFn func() time.Time
}

// ensure that CredentialsGopassTotp implements the ICredentials interface
var _ ICredentialsSource = (*CredentialsGopassTotpSource)(nil)

func (c *CredentialsGopassTotpSource) Type() CredentialSourceType {
	return CredentialSourceTypeGopassTotp
}
func (c *CredentialsGopassTotpSource) Resolve() (ResolvedCredentials, error) {
	Password, err := c.Api.GetPassword(c.Secret)
	if err != nil {
		return ResolvedCredentials{}, fmt.Errorf("failed to get password for: %s", c.Secret)
	}
	Username, err := c.Api.GetUsername(c.Secret)
	if err != nil {
		return ResolvedCredentials{}, fmt.Errorf("failed to get username for: %s", c.Secret)
	}
	Totp, err := c.Api.GetOtp(c.Secret, c.timestampFn())
	if err != nil {
		return ResolvedCredentials{}, fmt.Errorf("failed to get totp for: %s. Reason: %+v", c.Secret, err)
	}

	return ResolvedCredentials{
		UsernameAndPasswordAndTotp: UsernameAndPasswordAndTotp{
			Username,
			Password,
			Totp,
		},
		Type: c.Type(),
	}, nil
}
func (c *CredentialsGopassTotpSource) GetTimestamp() time.Time {
	fn := c.timestampFn
	if fn == nil {
		return time.Now().UTC()
	}
	return fn()
}

func (c *CredentialsGopassTotpSource) SetTimestampFn(fn func() time.Time) {
	c.timestampFn = fn
}

// CredentialsKeychain is a struct that contains the credentials for a source.
type CredentialsKeychainSource struct {
	UsernameKey string
	PasswordKey string
}

// ensure that CredentialsKeychain implements the ICredentials interface
var _ ICredentialsSource = (*CredentialsKeychainSource)(nil)

func (c *CredentialsKeychainSource) Type() CredentialSourceType { return CredentialSourceTypeKeychain }
func (c *CredentialsKeychainSource) Resolve() (ResolvedCredentials, error) {
	return ResolvedCredentials{
		UsernameAndPassword: UsernameAndPassword{
			// TODO: implement keychain integration
			Username: os.Getenv(c.UsernameKey),
			Password: os.Getenv(c.PasswordKey),
		},
		Type: c.Type(),
	}, nil
}

// CredentialsKeychainTotp is a struct that contains the credentials for a source.
type CredentialsKeychainTotpSource struct {
	UsernameKey string
	PasswordKey string
	TotpKey     string
}

// ensure that CredentialsKeychainTotp implements the ICredentials interface
var _ ICredentialsSource = (*CredentialsKeychainTotpSource)(nil)

func (c *CredentialsKeychainTotpSource) Type() CredentialSourceType {
	return CredentialSourceTypeKeychainTotp
}
func (c *CredentialsKeychainTotpSource) Resolve() (ResolvedCredentials, error) {
	return ResolvedCredentials{
		UsernameAndPasswordAndTotp: UsernameAndPasswordAndTotp{
			// TODO: implement keychain integration
			Username: os.Getenv(c.UsernameKey),
			Password: os.Getenv(c.PasswordKey),
			Totp:     os.Getenv(c.TotpKey),
		},
		Type: c.Type(),
	}, nil
}

// A fat union
type CredentialsSource struct {
	CredentialsFileSource
	CredentialsEnvSource
	CredentialsGopassSource
	CredentialsGopassTotpSource
	CredentialsKeychainSource
	CredentialsKeychainTotpSource
	Type CredentialSourceType
}

// ICredentials is an interface that contains the credentials for a source.
type ICredentialsSource interface {
	Type() CredentialSourceType
	Resolve() (ResolvedCredentials, error)
}

type ResolvedCredentials struct {
	UsernameAndPassword
	UsernameAndPasswordAndTotp
	Type CredentialSourceType
}

type Credentials struct {
	ResolvedCredentials
	CredentialsSource
	Type CredentialSourceType
}

// accepts a generic object, inspects a key "type", and returns a struct with the embeded struct filled out.
func NewCredentials(source map[string]interface{}) Credentials {
	var output Credentials
	output.Type = CredentialSourceType(source["type"].(string))

	switch source["type"] {
	case CredentialSourceTypeFile:
		output.CredentialsFileSource = CredentialsFileSource{
			Username: source["username"].(string),
			Password: source["password"].(string),
		}
		output.ResolvedCredentials, _ = output.CredentialsFileSource.Resolve()

	case CredentialSourceTypeEnv:
		output.CredentialsEnvSource = CredentialsEnvSource{
			UsernameKey: source["usernameKey"].(string),
			PasswordKey: source["passwordKey"].(string),
		}
		output.ResolvedCredentials, _ = output.CredentialsEnvSource.Resolve()

	case CredentialSourceTypeGopass:
		output.CredentialsGopassSource = CredentialsGopassSource{
			UsernameKey: source["usernameKey"].(string),
			PasswordKey: source["passwordKey"].(string),
		}
		output.ResolvedCredentials, _ = output.CredentialsGopassSource.Resolve()

	case CredentialSourceTypeGopassTotp:
		output.CredentialsGopassTotpSource = CredentialsGopassTotpSource{
			UsernameKey: source["usernameKey"].(string),
			PasswordKey: source["passwordKey"].(string),
			TotpKey:     source["totpKey"].(string),
		}
		output.ResolvedCredentials, _ = output.CredentialsGopassTotpSource.Resolve()

	case CredentialSourceTypeKeychain:
		output.CredentialsKeychainSource = CredentialsKeychainSource{
			UsernameKey: source["usernameKey"].(string),
			PasswordKey: source["passwordKey"].(string),
		}
		output.ResolvedCredentials, _ = output.CredentialsKeychainSource.Resolve()

	case CredentialSourceTypeKeychainTotp:
		output.CredentialsKeychainTotpSource = CredentialsKeychainTotpSource{
			UsernameKey: source["usernameKey"].(string),
			PasswordKey: source["passwordKey"].(string),
			TotpKey:     source["totpKey"].(string),
		}
		output.ResolvedCredentials, _ = output.CredentialsKeychainTotpSource.Resolve()

	}

	return output
}
