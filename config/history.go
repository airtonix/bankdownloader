package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"dario.cat/mergo"
	"github.com/airtonix/bank-downloaders/core"
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

func (h *History) Save() error {
	// marshal contents into bytes[]

	SaveYamlFile(h, historyFileName)
	return nil
}

// History Singleton
var history History
var historyFileName string
var defaultHistory = []History{}

func LoadHistory(historyFile string) {
	historyFilename := "history.yaml"
	historyFilepath := core.GetUserFilePath(historyFilename)

	// envvar runtime override
	if envhistoryFile := os.Getenv("BANKSCRAPER_CONFIG"); envhistoryFile != "" {
		NewHistory(envhistoryFile)

		// args filename override
	} else if historyFile != "" {
		NewHistory(historyFile)

		// config file in current directory
	} else if core.FileExists(historyFilename) {
		NewHistory(historyFilename)

		// config file in XDG directory
	} else if core.FileExists(historyFilepath) {
		NewHistory(historyFilepath)

	} else {
		EnsureStoragePath(historyFilepath)
	}
}

func NewHistory(historyFilePath string) (History, error) {
	var historyJson interface{}
	var err error

	content, err := LoadYamlFile(historyFilePath)
	if core.AssertErrorToNilf("could not load history file: %w", err) {
		return history, err
	}

	err = yaml.Unmarshal(content, &historyJson)
	if core.AssertErrorToNilf("could not unmarshal history file: %w", err) {
		return history, err
	}

	err = schema.Validate(historyJson)
	if core.AssertErrorToNilf("could not validate history file: %w", err) {
		return history, errors.New(fmt.Sprintf("Invalid configuration\n%#v", err))
	}

	err = yaml.Unmarshal(content, &history)
	if core.AssertErrorToNilf("could not unmarshal history file: %w", err) {
		return history, err
	}

	err = mergo.Merge(&history, defaultHistory, mergo.WithOverrideEmptySlice)
	if core.AssertErrorToNilf("could not merge history file: %w", err) {
		return history, err
	}

	historyFileName = historyFilePath
	return history, nil
}

func GetHistory() History {
	return history
}
