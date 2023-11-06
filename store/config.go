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

// A job is a set of instructions for downloading transactions from a source
// We would download transactions for a set of accounts for a number of days
type ConfigSource struct {
	Name     string       `json:"name" yaml:"name"`         // the name of the source. This is used to lookup the source in the registry
	Accounts []JobAccount `json:"accounts" yaml:"accounts"` // the accounts to download transactions for
	Config   any          // the source specific config, ignore it when marshalling
}

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
		mergo.WithOverrideEmptySlice,
	)
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
