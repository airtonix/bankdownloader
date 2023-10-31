package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/schemas"
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
	Events []HistoryEvent `yaml:"history"`
}

func (h *History) GetEvents(
	source string,
	accountNo string,
	accountName string,
) (time.Time, time.Time, error) {
	var fromDate time.Time
	var toDate time.Time
	format := GetDateFormat()

	for _, event := range h.Events {
		if event.Source == source && event.Account == accountNo {
			fromDate := core.StringToDate(event.FromDate, format)
			toDate := core.StringToDate(event.ToDate, format)
			return fromDate, toDate, nil
		}
	}
	errorMsg := errors.New(
		fmt.Sprintf("could not find history for source: %s, account: %s",
			source,
			accountName,
		),
	)

	return fromDate, toDate, errorMsg
}

// determine next date to fetch transactions from
func (h *History) GetNextDate(
	source string,
	accountNo string,
	accountName string,
	daysToFetch int,
) (time.Time, error) {
	_, toDate, err := h.GetEvents(
		source,
		accountNo,
		accountName,
	)

	if core.AssertErrorToNilf("could not get history: %w", err) || toDate.IsZero() {
		return core.GetTodayMinusDays(daysToFetch), nil
	}

	return core.GetDaysAgo(toDate, daysToFetch), nil
}

// save the event
func (h *History) SaveEvent(
	source string,
	accountNo string,
	accountName string,
	fromDate time.Time,
	toDate time.Time,
) error {
	var event HistoryEvent
	format := GetDateFormat()

	event.Source = source
	event.Account = accountNo
	event.FromDate = fromDate.Format(format)
	event.ToDate = toDate.Format(format)

	h.Events = append(h.Events, event)

	return nil
}

func (this *History) Save() error {
	// marshal contents into bytes[]

	SaveYamlFile(this, historyFilePath)
	return nil
}

// History Singleton
var history History
var historyFilePath string
var defaultHistory = []History{}
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

	if !core.FileExists(filepath) {
		log.Info("creating default history: ", filepath)
		CreateDefaultHistory(filepath)
	}

	historyObject, err := LoadYamlFile[History](
		filepath,
		schemas.GetHistorySchema(),
	)
	if core.AssertErrorToNilf("could not load history file: %w", err) {
		return history, err
	}

	// store the history as a singleton
	history = historyObject
	historyFilePath = filepath

	// also return it
	return history, nil
}

func CreateDefaultHistory(historyFilePath string) History {
	var defaultHistory History

	content, err := yaml.Marshal(defaultConfigTree)
	WriteFile("history.yaml", content)

	if core.AssertErrorToNilf("could not marshal default history: %w", err) {
		return defaultHistory
	}
	return defaultHistory
}

func GetHistory() History {
	return history
}
