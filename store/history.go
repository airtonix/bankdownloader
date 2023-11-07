package store

import (
	"errors"
	"sort"
	"time"

	"dario.cat/mergo"
	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/schemas"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type HistoryEvent struct {
	Source          string `yaml:"source"`
	lastDateFetched string `yaml:"lastDateFetched"`
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
		there := events[i].lastDateFetched
		here := events[j].lastDateFetched
		return there < here
	})

	// get the last event
	event := events[len(events)-1]

	return event, nil
}

type HistoryStrategy interface {
	Strategy() strategy
}

const (
	DaysAgo strategy = iota
	SinceLastDownload
)

type strategy int

// confirm that strategy implements HistoryStrategy
var _ HistoryStrategy = strategy(0)

func (h strategy) Strategy() strategy {
	return h
}

func NewHistoryStrategy(input string) strategy {
	switch input {
	case "days-ago":
		return DaysAgo
	case "since-last-download":
		return SinceLastDownload
	default:
		return DaysAgo
	}
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
	var fromDate time.Time
	var toDate time.Time

	// Strategy: DaysAgo
	if strategy.Strategy() == DaysAgo {
		fromDate, toDate = h.GetDaysAgo(daysToFetch)
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
			fromDate, toDate = h.GetDaysAgo(daysToFetch)
			return fromDate, toDate, nil
		}

		// fromDate will be the last toDate
		fromDate := core.StringToDate(event.lastDateFetched, time.RFC3339)
		// toDate will be fromDate plus daysToFetch
		toDate = fromDate.AddDate(0, 0, daysToFetch)

		// if toDate is beyond yesterday, default to yesterday
		if toDate.After(core.GetTodayMinusDays(1)) {
			toDate = core.GetTodayMinusDays(1)
		}

		return fromDate, toDate, nil
	}

	return fromDate, toDate, errors.New("Unable to calculate next date range")
}

func (this *History) GetDaysAgo(
	daysToFetch int,
) (time.Time, time.Time) {
	toDate := core.GetTodayMinusDays(1)
	fromDate := core.GetTodayMinusDays(daysToFetch)
	return fromDate, toDate
}

// save the event
func (this *History) SaveEvent(source string, accountNo string, fromDate time.Time, toDate time.Time) {
	format := GetDateFormat()
	event := HistoryEvent{
		Source:          source,
		AccountNumber:   accountNo,
		lastDateFetched: toDate.Format(format),
	}
	this.Events = append(this.Events, event)
	this.Save()
}

func (this *History) Save() error {
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
		this,
		mergo.WithOverrideEmptySlice,
	)
	if core.AssertErrorToNilf("Problem preparing history to save: %w", err) {
		return err
	}

	SaveYamlFile(output, historyFilePath)
	return nil
}

// History Singleton
var history History
var historyFilePath string
var defaultHistory = History{}
var defaultHistoryTree = &yaml.Node{
	Kind: yaml.DocumentNode,
	Content: []*yaml.Node{
		{
			Kind: yaml.MappingNode,
			Content: []*yaml.Node{
				{
					Kind:        yaml.ScalarNode,
					Value:       "events",
					HeadComment: "# yaml-language-server: $schema=https://raw.githubusercontent.com/airtonix/bankdownloader/master/schemas/history.json",
				},
				{
					Kind:    yaml.SequenceNode,
					Content: []*yaml.Node{},
				},
			},
		},
	},
}

func NewHistory(filepathArg string) (History, error) {
	filename := "history.yaml"
	filepath := core.ResolveFileArg(
		filepathArg,
		"BANKDOWNLOADER_HISTORY",
		filename,
	)

	err := mergo.Merge(
		&history,
		defaultHistory,
		mergo.WithOverrideEmptySlice)

	if core.AssertErrorToNilf("could not ensure default history values: %w", err) {
		return history, err
	}

	if !core.FileExists(filepath) {
		CreateDefaultHistory(filepath)
	}

	var historyObject History
	err = LoadYamlFile[History](
		filepath,
		schemas.GetHistorySchema(),
		&historyObject,
	)
	if core.AssertErrorToNilf("could not load history file: %w", err) {
		return history, err
	}
	log.Info("history ready: ", filepath)

	// store the history as a singleton
	history = historyObject
	historyFilePath = filepath

	// also return it
	return history, nil
}

func CreateDefaultHistory(historyFilePath string) History {
	var defaultHistory History

	log.Info("creating default config: ", configFilePath)

	content, err := yaml.Marshal(defaultHistoryTree)
	WriteFile(historyFilePath, content)

	if core.AssertErrorToNilf("could not marshal default history: %w", err) {
		return defaultHistory
	}
	return defaultHistory
}

func GetHistory() History {
	return history
}
