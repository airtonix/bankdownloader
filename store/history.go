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
	Source   string `yaml:"source"`
	Account  string `yaml:"account"`
	FromDate string `yaml:"fromDate"`
	ToDate   string `yaml:"toDate"`
}

type History struct {
	Events []HistoryEvent `yaml:"events"`
}

func (h *History) GetEvents(
	source string,
	accountNo string,
	accountName string,
) []HistoryEvent {
	events := []HistoryEvent{}

	for _, event := range h.Events {
		if event.Source == source && event.Account == accountNo {

			// push event onto events
			events = append(events, event)
		}
	}
	return events
}

func (h *History) GetLatestEvent(
	source string,
	accountNo string,
	accountName string,
) (HistoryEvent, error) {

	events := h.GetEvents(
		source,
		accountNo,
		accountName,
	)

	// if there are no events, return zero dates
	if len(events) == 0 {
		return HistoryEvent{}, errors.New("no events found")
	}

	// sort events by toDate
	sort.Slice(events, func(i, j int) bool {
		toDate := events[i].ToDate
		fromDate := events[j].ToDate
		return toDate < fromDate
	})

	// get the last event
	event := events[len(events)-1]

	return event, nil
}

// determine next date to fetch transactions from
func (h *History) GetNextEvent(
	source string,
	accountNo string,
	accountName string,
	daysToFetch int,
) HistoryEvent {
	logrus.Debugln("Days to fetch", daysToFetch)
	defaultFromDate := core.GetTodayMinusDays(daysToFetch).Format(GetDateFormat())
	defaultToDate := core.GetToday().Format(GetDateFormat())

	// by default we fetch from today - daysToFetch
	nextEvent := HistoryEvent{
		Source:   source,
		Account:  accountNo,
		FromDate: defaultFromDate,
		ToDate:   defaultToDate,
	}

	event, err := h.GetLatestEvent(
		source,
		accountNo,
		accountName,
	)

	if err != nil {
		logrus.Debugln("No events found, using default")
		return nextEvent
	}

	logrus.Debugln("latest event", event)

	// fromDate will be the last toDate
	fromDate := event.ToDate

	// if fromDate is less than daysToFetch days ago, compute it
	daysSinceLastEvent := core.GetDaysBetweenDates(
		core.StringToDate(fromDate, GetDateFormat()),
		core.GetToday(),
	)

	if daysSinceLastEvent < daysToFetch {
		fromDate = core.GetTodayMinusDays(daysSinceLastEvent).Format(GetDateFormat())
		nextEvent.FromDate = fromDate
	}

	return nextEvent
}

// save the event
func (this *History) SaveEvent(source string, accountNo string, fromDate time.Time, toDate time.Time) {
	var event HistoryEvent
	format := GetDateFormat()

	event.Source = source
	event.Account = accountNo
	event.FromDate = fromDate.Format(format)
	event.ToDate = toDate.Format(format)

	this.Events = append(this.Events, event)
}

func (this *History) Save() error {
	// marshal contents into bytes[]

	SaveYamlFile(this, historyFilePath)
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
