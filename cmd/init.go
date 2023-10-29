package cmd

import (
	"github.com/airtonix/bank-downloaders/config"
	"github.com/spf13/cobra"
)

var initDmd = &cobra.Command{
	Use:   "init",
	Short: "initialise configuration",
	Run: func(cmd *cobra.Command, args []string) {
		userconfig := config.GetConfig()
		userconfig.Save()
	},
}

func init() {
	rootCmd.AddCommand(initDmd)
}
