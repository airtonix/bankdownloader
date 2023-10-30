package cmd

import (
	"log"

	"github.com/airtonix/bank-downloaders/config"
	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/sources"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "dwnloads transactions from a source",
	Run: func(cmd *cobra.Command, args []string) {
		registry := sources.GetRegistry()
		userhistory := config.GetHistory()
		jobs := config.GetJobs()

		// loop through the jobs
		// for each job, download the transactions
		for _, job := range jobs {
			// get the source
			source := registry.GetSource(job.Source)

			for _, account := range job.Accounts {
				// determine the next date to fetch transactions from
				// based on the number of days to fetch
				// and the last date we fetched transactions from
				// if the last date is empty, then we fetch from today - daysToFetch
				// if the last date is not empty, then we fetch from last date - daysToFetch
				nextFromDate, err := userhistory.GetNextDate(
					job.Source,
					account.Number,
					account.Name,
					job.DaysToFetch,
				)
				if core.AssertErrorToNilf("could not get next date: %w", err) {
					continue
				}

				nextToDate := core.ToStartOfDay(nextFromDate.AddDate(0, 0, job.DaysToFetch))

				source.OpenBrowser()

				// login to the source
				err = source.Login(
					job.Credentials,
				)

				// download the transactions
				transactionFilename, err := source.DownloadTransactions(
					account.Number,
					account.Name,
					job.Format,
					nextFromDate,
					nextToDate,
				)

				if core.AssertErrorToNilf("could not download transactions: %w", err) {
					continue
				}

				log.Printf(
					"Downloaded transactions for %s from %s to %s as %s",
					account.Name, nextFromDate, nextToDate, transactionFilename,
				)

				userhistory.SaveEvent(
					job.Source,
					account.Number,
					account.Name,
					nextFromDate,
					nextToDate,
				)
			}
		}

		userhistory.Save()
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
