package credentials

import (
	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/store/credentials/clients"
)

// CredentialsGopass is a struct that contains the credentials for a source.
type CredentialsGopassInput struct {
	Secret      string
	UsernameKey string
	PasswordKey string
}
type CredentialsGopassSource struct {
	CredentialsGopassInput
	Api *clients.GopassSecretResolver
}

// ensure that CredentialsGopass implements the ICredentials interface
var _ ICredentialsSource = (*CredentialsGopassSource)(nil)

func NewCredentialsGopassSource(input map[string]interface{}) *CredentialsGopassSource {
	resolvedInput := CredentialsGopassInput{
		UsernameKey: "username",
		PasswordKey: "password",
	}
	if input["secret"] == nil {
		panic("Using gopass requires a secret path to be set.")
	}

	resolvedInput.Secret = input["secret"].(string)

	if input["usernameKey"] != nil {
		resolvedInput.UsernameKey = input["usernameKey"].(string)
	}

	if input["passwordKey"] != nil {
		resolvedInput.PasswordKey = input["passwordKey"].(string)
	}

	return &CredentialsGopassSource{
		CredentialsGopassInput: resolvedInput,
		Api:                    clients.NewGopassResolver(),
	}
}

func (c *CredentialsGopassSource) Type() CredentialSourceType { return CredentialSourceTypeGopass }
func (c *CredentialsGopassSource) Resolve() ResolvedCredentials {
	Password, err := c.Api.GetPassword(c.Secret)
	core.AssertErrorToNilf("failed to get password: %s", err)

	Username, err := c.Api.GetUsername(c.Secret)
	core.AssertErrorToNilf("failed to get username: %s", err)

	return ResolvedCredentials{
		UsernameAndPassword: UsernameAndPassword{
			Username,
			Password,
		},
		Type: c.Type(),
	}
}
