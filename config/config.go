package config

import (
	"errors"
	"fmt"
	"os"

	_ "embed"

	"dario.cat/mergo"
	"github.com/airtonix/bank-downloaders/core"

	"gopkg.in/yaml.v3"
)

type Account struct {
	Name   string `json:"name" yaml:"name"`
	Number string `json:"number" yaml:"number"`
}

// A job is a set of instructions for downloading transactions from a source
// We would download transactions for a set of accounts for a number of days
type Job struct {
	Source      string    `json:"source" yaml:"source"`           // the name of the source. This is used to lookup the source in the registry
	Credentials any       `json:"credentials" yaml:"credentials"` // the credentials to use for the source
	Format      string    `json:"format" yaml:"format"`           // the format to download the transactions in
	Accounts    []Account `json:"accounts" yaml:"accounts"`       // the accounts to download transactions for
	DaysToFetch int       `json:"daysToFetch" yaml:"daysToFetch"` // the number of days to fetch transactions for
}

type Config struct {
	DateFormat string `json:"dateFormat" yaml:"dateFormat"` // the format to use for dates
	Jobs       []Job  `json:"jobs" yaml:"jobs"`
}

// Config Singleton
var config Config
var configFileName string
var defaultConfig = Config{}

func (c *Config) Save() error {
	// marshal contents into bytes[]
	SaveYamlFile(c, configFileName)
	return nil
}

func LoadConfig(configFile string) {
	configFilename := "config.yaml"
	configFilepath := core.GetUserFilePath(configFilename)

	// envvar runtime override
	if envConfigFile := os.Getenv("BANKSCRAPER_CONFIG"); envConfigFile != "" {
		NewConfig(envConfigFile)

		// args filename override
	} else if configFile != "" {
		NewConfig(configFile)

		// config file in current directory
	} else if core.FileExists(configFilename) {
		NewConfig(configFilename)

		// config file in XDG directory
	} else if core.FileExists(configFilepath) {
		NewConfig(configFilepath)

	} else {
		InitializeConfig(config, configFilepath)
		config.Save()
	}
}

func NewConfig(configFilePath string) (Config, error) {
	var configJson interface{}
	var err error

	content, err := LoadYamlFile(configFilePath)

	err = yaml.Unmarshal(content, &configJson)
	if core.AssertErrorToNilf("could not unmarshal config file: %w", err) {
		return config, err
	}

	err = schema.Validate(configJson)
	if core.AssertErrorToNilf("could not validate config file: %w", err) {
		return config, errors.New(fmt.Sprintf("Invalid configuration\n%#v", err))
	}

	err = yaml.Unmarshal(content, &config)
	if core.AssertErrorToNilf("could not unmarshal config file: %w", err) {
		return config, err
	}

	InitializeConfig(config, configFilePath)

	return config, nil
}

func InitializeConfig(c Config, filepath string) error {
	var err error

	err = mergo.Merge(&config, defaultConfig, mergo.WithOverrideEmptySlice)
	if core.AssertErrorToNilf("could not merge config file: %w", err) {
		return err
	}

	configFileName = filepath
	return nil
}

func GetConfig() Config {
	return config
}

func GetDateFormat() string {
	return config.DateFormat
}

func GetJobs() []Job {
	return config.Jobs
}
