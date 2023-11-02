package cmd

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/airtonix/bank-downloaders/meta"
	"github.com/airtonix/bank-downloaders/store"
	log "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"github.com/spf13/cobra"
)

var prefix = "bankscraper"

var configFileArg string
var historyFileArg string
var now string

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
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func Initialize() {
	InitLogger(nil)
	store.NewConfig(configFileArg)
	store.NewHistory(historyFileArg)
}

func InitLogger(hook log.Hook) {
	formatter := log.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
		PadLevelText:     true,
	}

	if os.Getenv(debugEnvVarName) == "true" {
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
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
				Level:      log.InfoLevel,
				Formatter:  &log.JSONFormatter{},
			})
			if err == nil {
				log.AddHook(rotateFileHook)
			}
		}
	}

	log.SetFormatter(&formatter)

	if hook != nil {
		log.AddHook(hook)
	}
}
