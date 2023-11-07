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
		strategy := store.NewHistoryStrategy(cmd.Flag("range-strategy").Value.String())

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
				fromDate, toDate, err := history.GetDownloadDateRange(
					source.GetName(),
					account.Number,
					source.GetDaysToFetch(),
					strategy,
				)

				daysSinceLastEvent := core.GetDaysBetweenDates(toDate, core.GetToday())

				// if the days since the last event is less than a day
				// then we don't need to download transactions
				if daysSinceLastEvent < 1 {
					logrus.Infoln(
						fmt.Sprintf(
							"Skipping %s, Since at least one day has not passed since last download. (%d days)",
							account.Name, daysSinceLastEvent,
						),
					)
					continue
				}

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
				history.SaveEvent(
					source.GetName(),
					account.Number,
					toDate,
				)
			}
		}
	},
}

func init() {
	// TODO: https://github.com/spf13/pflag/issues/236#issuecomment-931600452
	strategyEnum := core.EnumFlag([]string{"days-ago", "since-last-download"}, "days-ago")
	downloadCmd.Flags().VarP(strategyEnum, "range-strategy", "r", "strategy to use when determining the date range to download")

	rootCmd.AddCommand(downloadCmd)
}
