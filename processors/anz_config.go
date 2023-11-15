package processors

import (
	"dario.cat/mergo"
	"github.com/airtonix/bank-downloaders/core"
	"github.com/kr/pretty"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type AnzCredentials struct {
	// username is the customer registration number for ANZ internet banking.
	Username string
	// password is the password for ANZ internet banking.
	Password string
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

// func (credentials *AnzCredentials) UnmarshalYAML(node *yaml.Node) error {
// 	var raw interface{}
// 	if err := node.Decode(&raw); err != nil {
// 		return err
// 	}

// 	// parse the anz credentials
// 	if err := mergo.Merge(credentials, &AnzCredentials{
// 		Username: raw.(map[string]interface{})["username"].(string),
// 		Password: raw.(map[string]interface{})["password"].(string),
// 	}); err != nil {
// 		return err
// 	}

// 	return nil
// }

type AnzConfig struct {
	Credentials     AnzCredentials
	ProcessorConfig `mapstructure:",squash"`
}

var defaultAnzConfig = AnzConfig{
	Credentials: AnzCredentials{
		Username: "",
		Password: "",
	},
	ProcessorConfig: ProcessorConfig{
		Domain:         "https://login.anz.com",
		ExportFormat:   "Quicken(QIF)",
		OutputTemplate: "{{ .SourceSlug }}_{{ .Account.NumberSlug }}_{{ .DateRange.FromUnix }}-{{ .DateRange.ToUnix }}.qif",
		DaysToFetch:    30,
	},
}

func NewAnzConfig(config map[string]interface{}) (*AnzConfig, error) {
	var output AnzConfig
	var err error

	// merge default config
	err = mergo.Merge(
		&output,
		defaultAnzConfig,
	)

	if core.AssertErrorToNilf("could not merge default config: %w", err) {
		return nil, err
	}

	// merge credentials config
	var anzConfig AnzConfig
	mapstructure.Decode(config, &anzConfig)

	// var anzCredentials AnzCredentials
	err = mergo.Merge(
		&output, anzConfig,
		mergo.WithSliceDeepCopy,
	)
	if core.AssertErrorToNilf("could not merge config: %w", err) {
		return nil, err
	}
	logrus.Debugf("anzConfig: %v", pretty.Sprint(output))

	return &anzConfig, nil
}
