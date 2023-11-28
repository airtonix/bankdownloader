package cmd

import (
	"fmt"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/store"
	"github.com/airtonix/bank-downloaders/store/credentials"
	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

var configDmd = &cobra.Command{
	Use:   "config",
	Short: "show configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config := store.GetConfig()

		core.Header("Config")
		pretty.Println(config)

		core.Header("History")
		pretty.Println(store.GetHistory())

		core.Header("Credentials")
		for _, item := range config.Sources {

			core.Action(fmt.Sprintf("Resolving credentials: %s", item.Type))

			credentials := credentials.NewCredentials(
				item.Config.Credentials,
			)
			if credentials.ConfirmResolved() {
				core.Success("resolved")
			} else {
				core.Failure("did not resolved")
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(configDmd)
}
