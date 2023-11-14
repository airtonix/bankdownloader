package store

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"errors"
	"sort"
	"time"

	"dario.cat/mergo"
	"github.com/airtonix/bank-downloaders/core"
)

type HistoryEvent struct {
	Source          string `yaml:"source"`
	LastDateFetched string `yaml:"lastDateFetched"`
	AccountNumber   string `yaml:"accountNumber"`
}

type History struct {
	Events []HistoryEvent `yaml:"events"`
}

func (h *History) GetEvents(
	source string,
	accountNo string,
) []HistoryEvent {
	events := []HistoryEvent{}

	for _, event := range h.Events {
		if event.Source == source && event.AccountNumber == accountNo {

			// push event onto events
			events = append(events, event)
		}
	}
	return events
}

func (h *History) GetLatestEvent(
	source string,
	accountNo string,
) (HistoryEvent, error) {

	events := h.GetEvents(
		source,
		accountNo,
	)

	// if there are no events, return zero dates
	if len(events) == 0 {
		return HistoryEvent{}, errors.New("no events found")
	}

	// sort events by lastDateFetched
	sort.Slice(events, func(i, j int) bool {
		there := events[i].LastDateFetched
		here := events[j].LastDateFetched
		return there < here
	})

	// get the last event
	event := events[len(events)-1]

	return event, nil
}

// Calculate the next timeframe to download transactions for.
//
// Returns a tuple of `from` and `to` dates.
//
// There's always two dates: from and to.
// Both use the `daysToFetch` config to calculate the date range.
//
// The strategy is one of:
//
//   - `days-ago`: fetch transactions from `daysToFetch` days ago to yesterday.
//     `to` is always yesterday, and `from` is yesterday minus `daysToFetch` days ago.
//
//   - `since-last-download`: `from` always the last downloaded transaction `lastDateFetched`, and `to` is `lastDateFetched` plus `daysToFetch` days.
//     If `lastDateFetched` is not available, it will default to `days-ago`.
//     If `lastDateFetched` plus `daysToFetch` is beyond yesterday, it will default to yeserday.
func (h *History) GetDownloadDateRange(
	source string,
	accountNo string,
	daysToFetch int,
	strategy HistoryStrategy,
) (time.Time, time.Time, error) {
	toDate := core.GetTodayMinusDays(1)
	fromDate := core.GetDaysAgo(toDate, daysToFetch)

	// Strategy: DaysAgo
	if strategy.Strategy() == DaysAgo {
		return fromDate, toDate, nil
	}

	// Strategy: SinceLastDownload
	if strategy.Strategy() == SinceLastDownload {
		// try to detect a previously recorded event
		event, err := h.GetLatestEvent(
			source,
			accountNo,
		)

		// if there's no event, use the default
		if err != nil {
			logrus.Debugln("No events found, using default")
			return fromDate, toDate, nil
		}

		// fromDate will be the last toDate
		fromDate := core.StringToDate(event.LastDateFetched, time.RFC3339)
		// toDate will be fromDate plus daysToFetch
		toDate = fromDate.AddDate(0, 0, daysToFetch)
		yesterday := core.GetTodayMinusDays(1)

		// if toDate is beyond yesterday, default to yesterday
		if toDate.After(yesterday) {
			toDate = yesterday
		}

		daysSinceLastEvent := core.GetDaysBetweenDates(fromDate, toDate)
		if daysSinceLastEvent < 1 {
			return fromDate, toDate, errors.New("Days since last event is less than 1")
		}

		return fromDate, toDate, nil
	}

	return fromDate, toDate, errors.New("Unable to calculate next date range")
}

// save the event
func (h *History) SaveEvent(source string, accountNo string, toDate time.Time) {
	event := HistoryEvent{
		Source:          source,
		AccountNumber:   accountNo,
		LastDateFetched: toDate.Format(time.RFC3339),
	}
	h.Events = append(h.Events, event)
	h.Save()
}

func (h *History) Save() error {
	var output History
	var err error

	// TODO: not sure how to merge the default history tree with the history object
	// this throws an error that src and dest are not the same type
	// err = mergo.Merge(
	// 	&output,
	// 	&defaultHistoryTree,
	// 	mergo.WithOverrideEmptySlice,
	// )
	// if core.AssertErrorToNilf("Problem preparing history to save: %w", err) {
	// 	return err
	// }

	err = mergo.Merge(
		&output,
		h,
		mergo.WithOverrideEmptySlice,
	)
	if core.AssertErrorToNilf("Problem preparing history to save: %w", err) {
		return err
	}

	// SaveYamlFile(output, historyFilePath)
	return nil
}

var history History

func GetHistory() *History {
	return &history
}

var historyReader *viper.Viper

func NewHistoryReader() *viper.Viper {
	reader := viper.New()

	reader.SetConfigName("history")                            // name of config file (without extension)
	reader.SetConfigType("json")                               // REQUIRED if the config file does not have the extension in the name
	reader.AddConfigPath(configReader.GetString("configpath")) // call multiple times to add many search paths
	reader.AddConfigPath(".")
	reader.AddConfigPath(fmt.Sprintf("$HOME/.config/%s", appname)) // call multiple times to add many search paths
	reader.AddConfigPath(fmt.Sprintf("/etc/%s/", appname))         // path to look for the config file in

	reader.SetDefault("$schema", "https://raw.githubusercontent.com/airtonix/bankdownloader/master/schemas/history.json")

	reader.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	reader.WatchConfig()

	if err := reader.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
		}
	}

	return reader
}

func CreateNewHistoryFile() {
	// current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}
	historyFilePath := configReader.Get("configpath")
	if historyFilePath == nil {
		historyFilePath = fmt.Sprintf("%s/history.json", cwd)
	}

	// if the file exists, don't overwrite it
	if _, err := os.Stat(historyFilePath.(string)); err == nil {
		return
	}

	logrus.Infof("Creating new history file: %s", historyFilePath)
	if err := historyReader.SafeWriteConfigAs(historyFilePath.(string)); err != nil {
		logrus.Fatal(err)
	}
}

func InitHistory() {
	historyReader = NewHistoryReader()
	err := historyReader.Unmarshal(&history)
	core.AssertErrorToNilf("could not unmarshal history: %w", err)
	logrus.Debugln("history file", historyReader.ConfigFileUsed())
}

func init() {
	InitHistory()
}
