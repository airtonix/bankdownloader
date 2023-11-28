package processors

import (
	"errors"
	"time"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/store"
	"github.com/airtonix/bank-downloaders/store/credentials"
)

type IProcessor interface {
	// function to login to the source
	Login() error

	// function to download the transactions
	DownloadTransactions(
		accountName string,
		accountNumber string,
		fromDate time.Time,
		toDate time.Time,
	) (string, error)
}

func GetProcecssorFactory(
	processorName store.SourceType,
	config store.SourceConfig,
	credentials credentials.Credentials,
	automation *core.Automation,
) (IProcessor, error) {
	var processor IProcessor
	var err error

	switch processorName {
	case store.AnzSourceType:
		processor = NewAnzProcessor(
			config,
			credentials.UsernameAndPassword,
			automation,
		)
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
