package cmd

import (
	"github.com/airtonix/bank-downloaders/config"
	"github.com/airtonix/bank-downloaders/core"
	"github.com/playwright-community/playwright-go"
	"github.com/spf13/cobra"
)

var initDmd = &cobra.Command{
	Use:   "init",
	Short: "initialise configuration",
	Run: func(cmd *cobra.Command, args []string) {
		userconfig := config.GetConfig()
		userconfig.Save()

		err := playwright.Install()
		if core.AssertErrorToNilf("could not install playwright: %w", err) {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(initDmd)
}
