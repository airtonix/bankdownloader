package store

import (
	_ "embed"

	"dario.cat/mergo"
	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/schemas"
	"github.com/airtonix/bank-downloaders/sources"
	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DateFormat string `json:"dateFormat" yaml:"dateFormat"` // the format to use for dates
	Jobs       []Job  `json:"jobs" yaml:"jobs"`
}

func (c *Config) UnmarshalYAML(data []byte) error {
	log.Info("unmarshalling json config: ", c)
	return nil
}

type JobAccount struct {
	Name   string `json:"name" yaml:"name"`
	Number string `json:"number" yaml:"number"`
}

// A job is a set of instructions for downloading transactions from a source
// We would download transactions for a set of accounts for a number of days
type Job struct {
	SourceName string         `json:"source" yaml:"source"`     // the name of the source. This is used to lookup the source in the registry
	Config     any            `json:"config" yaml:"config"`     // the source specific config (gets parsed further by source factory)
	Accounts   []JobAccount   `json:"accounts" yaml:"accounts"` // the accounts to download transactions for
	Source     sources.Source `json:"-" yaml:"-"`               // the actual source once we unmarshal ignore it when marshalling
}

func (job *Job) UnmarshalYAML(value *yaml.Node) error {
	var raw interface{}
	if err := value.Decode(&raw); err != nil {
		return err
	}

	sourceName := raw.(map[string]interface{})["source"].(string)
	job.SourceName = sourceName

	source, err := sources.GetSourceFactory(sourceName)

	core.AssertErrorToNilf("could not get source factory", err)

	job.Source = source
	job.Source.LoadConfig(value)
	job.Accounts = []JobAccount{}

	for _, account := range raw.(map[string]interface{})["accounts"].([]interface{}) {
		job.Accounts = append(job.Accounts, JobAccount{
			Name:   account.(map[string]interface{})["name"].(string),
			Number: account.(map[string]interface{})["number"].(string),
		})
	}

	return nil
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
					Value: "DD/MM/YYYY",
				},
				{
					Kind:  yaml.ScalarNode,
					Value: "jobs",
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

	err := mergo.Merge(
		&config,
		defaultConfig,
		mergo.WithOverrideEmptySlice)

	if core.AssertErrorToNilf("could not ensure default config values: %w", err) {
		return config, err
	}

	if !core.FileExists(filepath) {
		CreateDefaultConfig(filepath)
	}
	var configObject Config
	err = LoadYamlFile[Config](
		filepath,
		schemas.GetConfigSchema(),
		&configObject,
	)
	if core.AssertErrorToNilf("could not load config file: %w", err) {
		return config, err
	}

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

func GetJobs() []Job {
	return config.Jobs
}
