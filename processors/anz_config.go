package processors

import (
	"dario.cat/mergo"
	"github.com/airtonix/bank-downloaders/core"
	"gopkg.in/yaml.v3"
)

type AnzCredentials struct {
	// username is the customer registration number for ANZ internet banking.
	Username string `json:"username" yaml:"username"`
	// password is the password for ANZ internet banking.
	Password string `json:"password" yaml:"password"`
}

func NewAnzCredentials(credentials map[string]interface{}) *AnzCredentials {
	var anzCredentials = &AnzCredentials{}

	// parse the anz credentials
	if credentials["username"] != nil {
		anzCredentials.Username = credentials["username"].(string)
	}

	if credentials["password"] != nil {
		anzCredentials.Password = credentials["password"].(string)
	}

	return anzCredentials
}

func (credentials *AnzCredentials) UnmarshalYAML(node *yaml.Node) error {
	var raw interface{}
	if err := node.Decode(&raw); err != nil {
		return err
	}

	// parse the anz credentials
	if err := mergo.Merge(credentials, &AnzCredentials{
		Username: raw.(map[string]interface{})["username"].(string),
		Password: raw.(map[string]interface{})["password"].(string),
	}); err != nil {
		return err
	}

	return nil
}

type AnzConfig struct {
	Credentials      *AnzCredentials `json:"credentials" yaml:"credentials"`
	*ProcessorConfig `json:",inline" yaml:",inline"`
}

var defaultAnzConfig = AnzConfig{
	Credentials: &AnzCredentials{
		Username: "",
		Password: "",
	},
	ProcessorConfig: &ProcessorConfig{
		Domain:      "https://login.anz.com",
		Format:      "Quicken(QIF)",
		DaysToFetch: 30,
	},
}

func NewAnzConfig(config map[string]interface{}) (*AnzConfig, error) {
	var anzConfig AnzConfig
	var err error

	// merge default config
	err = mergo.Merge(
		&anzConfig,
		defaultAnzConfig,
	)

	if core.AssertErrorToNilf("could not merge default config: %w", err) {
		return nil, err
	}

	// merge credentials config
	credentials := config["credentials"].(map[string]interface{})

	// var anzCredentials AnzCredentials
	err = mergo.Merge(
		&anzConfig,
		AnzConfig{
			Credentials:     NewAnzCredentials(credentials),
			ProcessorConfig: NewProcessorConfig(config),
		},
		mergo.WithSliceDeepCopy,
	)
	if core.AssertErrorToNilf("could not merge config: %w", err) {
		return nil, err
	}

	return &anzConfig, nil
}
