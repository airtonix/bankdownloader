package cmd

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/meta"
	"github.com/airtonix/bank-downloaders/store"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var envvarPrefix string = strings.ToUpper(meta.Name)
var debugEnvVarName string = fmt.Sprintf("%s_DEBUG", envvarPrefix)
var disableLogEnvVarName string = fmt.Sprintf("%s_DISABLE_LOG_FILE", envvarPrefix)

var rootCmd = &cobra.Command{
	Use:   meta.Name,
	Short: meta.Description,
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Show debug messages")
	rootCmd.PersistentFlags().StringP("config", "", "$HOME/.config/bankdownloader/config.json", "config file")
	rootCmd.PersistentFlags().StringP("history", "", "$HOME/.config/bankdownloader/history.json", "history file")
	rootCmd.PersistentFlags().BoolP("no-headless", "", false, "Don't run browser in headless mode?")

	err := viper.BindPFlag("configFile", rootCmd.PersistentFlags().Lookup("config"))
	core.AssertErrorToNilf("could not bind flags: %w", err)
	err = viper.BindPFlag("historyFile", rootCmd.PersistentFlags().Lookup("history"))
	core.AssertErrorToNilf("could not bind flags: %w", err)
	err = viper.BindPFlag("noHeadless", rootCmd.PersistentFlags().Lookup("no-headless"))
	core.AssertErrorToNilf("could not bind flags: %w", err)
	err = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	core.AssertErrorToNilf("could not bind flags: %w", err)

	cobra.OnInitialize(Initialize)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logrus.Fatal(err)
	}
}

func Initialize() {
	InitLogger(nil)
	core.EnsureChromeExists()

	err := store.InitialiseSchemas()
	core.AssertErrorToNilf("could not initialise schemas: %w", err)

	store.InitConfig(viper.GetString("configFile"))
	store.InitHistory(viper.GetString("historyFile"))
}

func InitLogger(hook logrus.Hook) {
	formatter := logrus.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
		PadLevelText:     true,
	}
	debugEnabled := viper.GetBool("debug") || os.Getenv(debugEnvVarName) == "true"

	if debugEnabled {
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.DebugLevel)
		formatter.CallerPrettyfier = func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return funcName, fmt.Sprintf(" [%s:%d]", path.Base(f.File), f.Line)
		}
	}

	if os.Getenv(disableLogEnvVarName) != "true" {
		p, err := store.EnsureLogFilePath()
		if err == nil {
			rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
				Filename:   p,
				MaxSize:    50,
				MaxBackups: 7,
				MaxAge:     30,
				Level:      logrus.InfoLevel,
				Formatter:  &logrus.JSONFormatter{},
			})
			if err == nil {
				logrus.AddHook(rotateFileHook)
			}
		}
	}

	logrus.SetFormatter(&formatter)

	if hook != nil {
		logrus.AddHook(hook)
	}
}
