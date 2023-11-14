package cmd

import (
	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/store"
	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

var configDmd = &cobra.Command{
	Use:   "config",
	Short: "show configuration",
	Run: func(cmd *cobra.Command, args []string) {
		core.Header("Config")
		pretty.Println(store.GetConfig())

		core.Header("History")
		pretty.Println(store.GetHistory())
	},
}

func init() {
	rootCmd.AddCommand(configDmd)
}
