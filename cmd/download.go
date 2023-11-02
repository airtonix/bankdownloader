package cmd

import (
	"fmt"

	"github.com/airtonix/bank-downloaders/store"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "dwnloads transactions from a source",
	Run: func(cmd *cobra.Command, args []string) {
		for _, job := range store.GetJobs() {
			fmt.Printf("Processing job: %s\n", job.SourceName)
			fmt.Printf("Processing jobConfig: %s\n", job.Config)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}

// func ProcessAccounts(
// 	accounts []store.Account,
// 	source SourceCommand,
// ) {
// 	userhistory := store.GetHistory()

// 	for _, account := range accounts {
// 		// determine the next date to fetch transactions from
// 		// based on the number of days to fetch
// 		// and the last date we fetched transactions from
// 		// if the last date is empty, then we fetch from today - daysToFetch
// 		// if the last date is not empty, then we fetch from last date - daysToFetch
// 		nextFromDate, err := userhistory.GetNextDate(
// 			job.Source,
// 			account.Number,
// 			account.Name,
// 			job.DaysToFetch,
// 		)
// 		if core.AssertErrorToNilf("could not get next date: %w", err) {
// 			continue
// 		}

// 		nextToDate := core.ToStartOfDay(nextFromDate.AddDate(0, 0, job.DaysToFetch))

// 		source.OpenBrowser()

// 		// login to the source
// 		err = source.Login(
// 			job.Credentials,
// 		)

// 		// download the transactions
// 		transactionFilename, err := source.DownloadTransactions(
// 			account.Number,
// 			account.Name,
// 			job.Format,
// 			nextFromDate,
// 			nextToDate,
// 		)

// 		if core.AssertErrorToNilf("could not download transactions: %w", err) {
// 			continue
// 		}

// 		log.Printf(
// 			"Downloaded transactions for %s from %s to %s as %s",
// 			account.Name, nextFromDate, nextToDate, transactionFilename,
// 		)

// 		userhistory.SaveEvent(
// 			job.Source,
// 			account.Number,
// 			account.Name,
// 			nextFromDate,
// 			nextToDate,
// 		)
// 	}
// }
