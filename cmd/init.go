package cmd

import (
	"github.com/airtonix/bank-downloaders/store"
	"github.com/spf13/cobra"
)

var initDmd = &cobra.Command{
	Use:   "init",
	Short: "initialise configuration",
	Run: func(cmd *cobra.Command, args []string) {

		// test if chrome can be found by executable name

		store.CreateNewConfigFile()
		store.CreateNewHistoryFile()
	},
}

func init() {
	rootCmd.AddCommand(initDmd)
}
