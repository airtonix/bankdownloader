package cmd

import (
	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/processors"
	"github.com/airtonix/bank-downloaders/store"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config related commands",
	Run: func(cmd *cobra.Command, args []string) {

		core.Header("Configured Sources")

		for _, item := range store.GetConfigSources() {
			source, err := processors.GetProcecssorFactory(
				item.Name,
				item.Config.(map[string]interface{}),
			)
			if err != nil {
				continue
			}
			source.Render()
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
