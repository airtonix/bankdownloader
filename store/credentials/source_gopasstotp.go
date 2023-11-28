package credentials

import (
	"time"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/store/credentials/clients"
)

// CredentialsGopassTotp is a struct that contains the credentials for a source.
type CredentialsGopassTotpSource struct {
	Secret      string
	UsernameKey string
	PasswordKey string
	TotpKey     string
	Api         *clients.GopassSecretResolver
	timestampFn func() time.Time
}

// ensure that CredentialsGopassTotp implements the ICredentials interface
var _ ICredentialsSource = (*CredentialsGopassTotpSource)(nil)

func (c *CredentialsGopassTotpSource) Type() CredentialSourceType {
	return CredentialSourceTypeGopassTotp
}
func (c *CredentialsGopassTotpSource) Resolve() ResolvedCredentials {
	Password, err := c.Api.GetPassword(c.Secret)
	core.AssertErrorToNilf("failed to get password for: %s", err)

	Username, err := c.Api.GetUsername(c.Secret)
	core.AssertErrorToNilf("failed to get username for: %s", err)

	Totp, err := c.Api.GetOtp(c.Secret, c.timestampFn())
	core.AssertErrorToNilf("failed to get totp for: %s", err)

	return ResolvedCredentials{
		UsernameAndPasswordAndTotp: UsernameAndPasswordAndTotp{
			Username,
			Password,
			Totp,
		},
		Type: c.Type(),
	}
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
