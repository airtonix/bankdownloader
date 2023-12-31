package store

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Account struct {
	Name           string
	Number         string
	OutputTemplate string `mapstructure:",omitempty"`
	Format         string `mapstructure:",omitempty"`
}

type SourceConfig struct {
	Domain         string
	ExportFormat   string
	OutputTemplate string
	DaysToFetch    int
	Credentials    map[string]interface{}
}

type SourceType string

var (
	AnzSourceType SourceType = "anz"
)

type Source struct {
	Type     SourceType
	Accounts []Account
	Config   SourceConfig
}

type Configuration struct {
	DateFormat string   `mapstructure:"dateformat"`
	Sources    []Source `mapstructure:"sources"`
}

var conf Configuration

func GetConfig() *Configuration {
	return &conf
}

var configReader *viper.Viper

func NewConfigReader(configFileArg string) *viper.Viper {
	configReader = viper.New()

	configReader.SetEnvPrefix(appname)
	configReader.AutomaticEnv()
	configReader.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	configReader.BindEnv("config")

	var configFileName = "config"
	var configFileExt = "json"
	if configFileArg != "" {
		// get the extension of the config file arg
		configFileExt = strings.TrimLeft(path.Ext(configFileArg), ".")
		configFileName = strings.TrimSuffix(configFileArg, path.Ext(configFileArg))
	} else {
		configFileArg = fmt.Sprintf("%s.%s", configFileName, configFileExt)
	}
	configFileDir := path.Dir(configFileArg)

	// lower case the name of the config file
	configReader.SetConfigName(configFileName)                           // name of config file (without extension)
	configReader.SetConfigType(configFileExt)                            // REQUIRED if the config file does not have the extension in the name
	configReader.AddConfigPath(configFileDir)                            // optionally look for config in the working directory
	configReader.AddConfigPath(".")                                      // optionally look for config in the working directory
	configReader.AddConfigPath(fmt.Sprintf("$HOME/.config/%s", appname)) // call multiple times to add many search paths
	configReader.AddConfigPath(fmt.Sprintf("/etc/%s/", appname))         // path to look for the config file in

	configReader.SetDefault("$schema", "https://raw.githubusercontent.com/airtonix/bankdownloader/master/store/config-schema.json")

	if err := configReader.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Errorf("Config file not found: %s", configFileArg)
		} else {
			// Config file was found but another error was produced
			logrus.Fatalf("Error reading config file: %s", err)
		}
	}
	return configReader
}

func CreateNewConfigFile() {
	// current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}
	configFilePath := configReader.Get("config")
	if configFilePath == nil {
		configFilePath = fmt.Sprintf("%s/config.json", cwd)
	}

	// if the file exists, don't overwrite it
	if _, err := os.Stat(configFilePath.(string)); err == nil {
		return
	}

	logrus.Infof("Creating new  config file: %s", configFilePath)
	if err := configReader.SafeWriteConfigAs(configFilePath.(string)); err != nil {
		logrus.Fatal(err)
	}
}

func InitConfig(configFileArg string) {
	configReader = NewConfigReader(configFileArg)
	err := configReader.Unmarshal(&conf)

	core.AssertErrorToNilf("could not unmarshal config: %w", err)
	logrus.Debugln("config file", configReader.ConfigFileUsed())
}
