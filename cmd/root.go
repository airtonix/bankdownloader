package cmd

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/airtonix/bank-downloaders/meta"
	"github.com/airtonix/bank-downloaders/store"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"github.com/spf13/cobra"
)

var configFileArg string
var historyFileArg string
var debugFlag bool
var headlessFlag bool

var envvarPrefix string = strings.ToUpper(meta.Name)
var debugEnvVarName string = fmt.Sprintf("%s_DEBUG", envvarPrefix)
var disableLogEnvVarName string = fmt.Sprintf("%s_DISABLE_LOG_FILE", envvarPrefix)

var rootCmd = &cobra.Command{
	Use:   meta.Name,
	Short: meta.Description,
}

func init() {
	cobra.OnInitialize(Initialize)
	rootCmd.PersistentFlags().StringVar(&configFileArg, "config", "", "config file (default is ./%s.yaml)")
	rootCmd.PersistentFlags().StringVar(&historyFileArg, "history", "", "history file (default is ./%s-history.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "shwo debug messages")
	rootCmd.PersistentFlags().BoolVar(&headlessFlag, "headless", true, "run browser in headless mode?")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logrus.Fatal(err)
	}
}

func Initialize() {
	InitLogger(nil)
	store.InitialiseSchemas()
	store.InitConfig()
	store.InitHistory()
}

func InitLogger(hook logrus.Hook) {
	formatter := logrus.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
		PadLevelText:     true,
	}
	debugEnabled := debugFlag || os.Getenv(debugEnvVarName) == "true"

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
