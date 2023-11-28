package credentials

import (
	"fmt"

	"github.com/airtonix/bank-downloaders/core"
)

type Credentials struct {
	ResolvedCredentials
	CredentialsSource
	Type CredentialSourceType
}

// accepts a generic object, inspects a key "type", and returns a struct with the embeded struct filled out.
func NewCredentials(source map[string]interface{}) Credentials {
	var output Credentials
	output.Type = CredentialSourceType(source["type"].(string))
	core.Action(fmt.Sprintf("Resolving credentials: %s", output.Type))

	switch output.Type {
	case CredentialSourceTypeFile:
		output.CredentialsFileSource = CredentialsFileSource{
			Username: source["username"].(string),
			Password: source["password"].(string),
		}
		output.ResolvedCredentials = output.CredentialsFileSource.Resolve()

	case CredentialSourceTypeEnv:
		output.CredentialsEnvSource = CredentialsEnvSource{
			UsernameKey: source["usernameKey"].(string),
			PasswordKey: source["passwordKey"].(string),
		}
		output.ResolvedCredentials = output.CredentialsEnvSource.Resolve()

	case CredentialSourceTypeGopass:
		output.CredentialsGopassSource = *NewCredentialsGopassSource(source)
		output.ResolvedCredentials = output.CredentialsGopassSource.Resolve()

	case CredentialSourceTypeGopassTotp:
		output.CredentialsGopassTotpSource = CredentialsGopassTotpSource{
			UsernameKey: source["usernameKey"].(string),
			PasswordKey: source["passwordKey"].(string),
			TotpKey:     source["totpKey"].(string),
		}
		output.ResolvedCredentials = output.CredentialsGopassTotpSource.Resolve()

	case CredentialSourceTypeKeychain:
		output.CredentialsKeychainSource = CredentialsKeychainSource{
			ServiceName: source["serviceName"].(string),
			Username:    source["username"].(string),
		}
		output.ResolvedCredentials = output.CredentialsKeychainSource.Resolve()

	default:
		panic(fmt.Sprintf("Unknown credential source type: %s", source["type"]))
	}

	return output
}
