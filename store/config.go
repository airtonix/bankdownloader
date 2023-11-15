package store

import (
	"fmt"
	"os"
	"strings"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Account struct {
	Name   string
	Number string
}

type Source struct {
	Name     string    `mapstructure:"name"`
	Accounts []Account `mapstructure:"accounts"`
	Config   any       `mapstructure:"config"`
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

func NewConfigReader() *viper.Viper {
	configReader = viper.New()

	configReader.SetEnvPrefix(appname)
	configReader.AutomaticEnv()
	configReader.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	configReader.BindEnv("configpath")

	logrus.Info(configReader.GetString("configpath"))

	// lower case the name of the config file
	configReader.SetConfigName("config")                                 // name of config file (without extension)
	configReader.SetConfigType("json")                                   // REQUIRED if the config file does not have the extension in the name
	configReader.AddConfigPath(configReader.GetString("config"))         // call multiple times to add many search paths
	configReader.AddConfigPath(".")                                      // optionally look for config in the working directory
	configReader.AddConfigPath(fmt.Sprintf("$HOME/.config/%s", appname)) // call multiple times to add many search paths
	configReader.AddConfigPath(fmt.Sprintf("/etc/%s/", appname))         // path to look for the config file in

	configReader.SetDefault("$schema", "https://raw.githubusercontent.com/airtonix/bankdownloader/master/schemas/config.json")

	if err := configReader.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
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
	configFilePath := configReader.Get("configpath")
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

func InitConfig() {
	configReader = NewConfigReader()
	err := configReader.Unmarshal(&conf)
	core.AssertErrorToNilf("could not unmarshal config: %w", err)
	logrus.Debugln("config file", configReader.ConfigFileUsed())
}

func init() {
	InitConfig()
}
