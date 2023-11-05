package store

import (
	_ "embed"
	"time"

	"dario.cat/mergo"
	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/schemas"
	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"
)

type JobAccount struct {
	Name   string `json:"name" yaml:"name"`
	Number string `json:"number" yaml:"number"`
}

// func (this *JobAccount) UnmarshalYAML(node *yaml.Node) error {
// 	var raw interface{}
// 	if err := node.Decode(&raw); err != nil {
// 		return err
// 	}

// 	name := raw.(map[string]interface{})["name"].(string)
// 	number := raw.(map[string]interface{})["number"].(string)

// 	this.Name = name
// 	this.Number = number

// 	return nil
// }

// A job is a set of instructions for downloading transactions from a source
// We would download transactions for a set of accounts for a number of days
type ConfigSource struct {
	Name     string       `json:"name" yaml:"name"`         // the name of the source. This is used to lookup the source in the registry
	Accounts []JobAccount `json:"accounts" yaml:"accounts"` // the accounts to download transactions for
	Config   any          // the source specific config, ignore it when marshalling
}

// func (source *ConfigSource) UnmarshalYAML(node *yaml.Node) error {
// 	var raw interface{}
// 	if err := node.Decode(&raw); err != nil {
// 		return err
// 	}

// 	sourceName := raw.(map[string]interface{})["name"].(string)
// 	sourceAccounts := raw.(map[string]interface{})["accounts"].([]interface{})
// 	sourceConfig := raw.(map[string]interface{})["config"].(map[string]interface{})

// 	// set the name
// 	source.Name = sourceName

// 	// set the accounts
// 	for _, account := range sourceAccounts {
// 		// create a new job account
// 		source.Accounts = append(source.Accounts, JobAccount{
// 			Name:   account.(map[string]interface{})["name"].(string),
// 			Number: account.(map[string]interface{})["number"].(string),
// 		})
// 	}

// 	source.Config = sourceConfig

// 	return nil
// }

type Config struct {
	DateFormat string         `json:"dateFormat" yaml:"dateFormat"` // the format to use for dates
	Sources    []ConfigSource `json:"sources" yaml:"sources"`
}

func (this *Config) Save() error {
	// marshal contents into bytes[]
	log.Info("saving config: ", configFilePath)
	SaveYamlFile(this, configFilePath)
	return nil
}

// Config Singleton
var config Config
var configFilePath string
var defaultConfig = &Config{}
var defaultConfigTree = &yaml.Node{
	Kind: yaml.DocumentNode,
	Content: []*yaml.Node{
		{
			Kind: yaml.MappingNode,
			Content: []*yaml.Node{
				{
					Kind:        yaml.ScalarNode,
					Value:       "dateFormat",
					HeadComment: "# yaml-language-server: $schema=https://raw.githubusercontent.com/airtonix/bankdownloader/master/schemas/config.json",
				},
				{
					Kind:  yaml.ScalarNode,
					Style: yaml.DoubleQuotedStyle,
					Value: time.RFC3339,
				},
				{
					Kind:  yaml.ScalarNode,
					Value: "sources",
				},
				{
					Kind:    yaml.SequenceNode,
					Content: []*yaml.Node{},
				},
			},
		},
	},
}

func NewConfig(filepathArg string) (Config, error) {
	filename := "config.yaml"
	filepath := core.ResolveFileArg(
		filepathArg,
		"BANKDOWNLOADER_CONFIG",
		filename,
	)

	// Initialise the config with default values
	err := mergo.Merge(
		&config,
		defaultConfig,
		mergo.WithOverrideEmptySlice)
	if core.AssertErrorToNilf("could not ensure default config values: %w", err) {
		return config, err
	}

	// Check if the config file exists
	// If it doesn't, create it
	if !core.FileExists(filepath) {
		CreateDefaultConfig(filepath)
	}

	// Load the config file and parse it
	var configObject Config
	err = LoadYamlFile[Config](
		filepath,
		schemas.GetConfigSchema(),
		&configObject,
	)
	if core.AssertErrorToNilf("could not load config file: %w", err) {
		return config, err
	}
	// merge it on top of the config
	err = mergo.Merge(
		&config,
		configObject,
		mergo.WithOverrideEmptySlice)
	if core.AssertErrorToNilf("could not merge user config values: %w", err) {
		return config, err
	}

	log.Info("config ready: ", filepath)

	// store the config as a singleton
	config = configObject
	configFilePath = filepath

	// also return it
	return config, nil
}

func CreateDefaultConfig(configFilePath string) Config {
	var defaultConfig Config

	log.Info("creating default config: ", configFilePath)

	content, err := yaml.Marshal(defaultConfigTree)
	WriteFile(configFilePath, content)

	if core.AssertErrorToNilf("could not marshal default config: %w", err) {
		return defaultConfig
	}
	return defaultConfig
}

func GetConfig() Config {
	return config
}

func GetDateFormat() string {
	return config.DateFormat
}

func GetConfigSources() []ConfigSource {
	return config.Sources
}
