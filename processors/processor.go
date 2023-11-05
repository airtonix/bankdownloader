package processors

import (
	"errors"
	"time"

	"github.com/airtonix/bank-downloaders/core"
	"gopkg.in/yaml.v3"
)

type IProcessor interface {
	GetName() string

	GetDaysToFetch() int

	GetFormat() string

	Render() error

	// function to login to the source
	Login(automation *core.Automation) error

	// function to download the transactions
	DownloadTransactions(
		accountName string,
		accountNumber string,
		fromDate time.Time,
		toDate time.Time,
		automation *core.Automation,
	) (string, error)
}

func GetProcecssorFactory(
	processorName string,
	config map[string]interface{},
) (IProcessor, error) {
	var processor IProcessor
	var err error

	switch processorName {
	case "anz":
		processor, err = NewAnzParsedProcessor(config)
		if err != nil {
			return nil, err
		}
		return processor, nil
	// case "commbank":
	// 	return &CommbankSource{}, nil
	// case "banksa":
	// 	return &BankSaSource{}, nil
	// case "ingorangeau":
	// 	return &IngOrangeAuSource{}, nil
	// case "westpac":
	// 	return &WestpacSource{}, nil
	// case "nab":
	// 	return &NABSource{}, nil
	default:
		return nil, errors.New("unsupported processor")
	}
}

type Processor struct {
	Name string // name of the source
}

func (processor *Processor) GetName() string {
	return processor.Name
}

type ProcessorConfig struct {
	Format      string `json:"format" yaml:"format"`           // format to download transactions in
	DaysToFetch int    `json:"daysToFetch" yaml:"daysToFetch"` // the number of days to fetch transactions for
}

func (config *ProcessorConfig) UnmarshalYAML(node *yaml.Node) error {
	var raw interface{}
	if err := node.Decode(&raw); err != nil {
		return err
	}

	config.Format = raw.(map[string]interface{})["format"].(string)
	config.DaysToFetch = raw.(map[string]interface{})["daysToFetch"].(int)

	return nil
}

func NewProcessorConfig(config map[string]interface{}) *ProcessorConfig {

	var processorConfig ProcessorConfig

	// pull the values from config and test they exist before casting them
	if config["daysToFetch"] != nil {
		processorConfig.DaysToFetch = config["daysToFetch"].(int)
	}

	if config["format"] != nil {
		processorConfig.Format = config["format"].(string)
	}

	return &processorConfig
}
