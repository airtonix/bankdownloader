package cmd

import (
	"fmt"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/processors"
	"github.com/airtonix/bank-downloaders/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "dwnloads transactions from a source",
	Run: func(cmd *cobra.Command, args []string) {
		history := store.GetHistory()
		automation := core.NewAutomation()
		core.Header("Downloading Sources")

		automation.OpenBrowser()

		for _, item := range store.GetConfigSources() {
			source, err := processors.GetProcecssorFactory(
				item.Name,
				item.Config.(map[string]interface{}),
			)
			if err != nil {
				continue
			}

			err = source.Login(automation)
			if core.AssertErrorToNilf("could not login: %w", err) {
				continue
			}

			for _, account := range item.Accounts {
				logrus.Infof("\nprocessing account: %s [%s]\n", account.Name, account.Number)
				event := history.GetNextEvent(
					source.GetName(),
					account.Number,
					account.Name,
					source.GetDaysToFetch(),
				)

				fromDate := core.StringToDate(event.FromDate, store.GetDateFormat())
				toDate := core.StringToDate(event.ToDate, store.GetDateFormat())

				filename, err := source.DownloadTransactions(
					account.Name,
					account.Number,
					fromDate,
					toDate,
					automation,
				)

				if core.AssertErrorToNilf("could not download transactions: %w", err) {
					continue
				}

				logrus.Infoln(
					fmt.Sprintf(
						"Downloaded transactions for %s from %s to %s as %s",
						account.Name, fromDate, toDate, filename,
					),
				)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
