package sources

import (
	"time"
)

type SourceCommand interface {
	Login(credentials any) error

	// function to download the transactions
	DownloadTransactions(
		accountName string,
		accountNumber string,
		format string,
		fromDate time.Time,
		toDate time.Time,
	) (string, error)

	OpenBrowser() error
}

type SourceRegistry struct {
	// a map of source names to sources
	// this is populated by the init() function in each source file
	// the init() function is called when the package is imported
	// see https://golang.org/doc/effective_go#init
	registry map[string]SourceCommand
}

func (r *SourceRegistry) Register(name string, source SourceCommand) {
	r.registry[name] = source
}

func (r *SourceRegistry) GetSource(name string) SourceCommand {
	return r.registry[name]
}

func NewSourceRegistry() SourceRegistry {
	return SourceRegistry{
		registry: make(map[string]SourceCommand),
	}
}

var registry SourceRegistry = NewSourceRegistry()

func InitRegistry() {
	// register all the sources

	// assign anz as a Source[any] to satisfy the compiler
	registry.Register("anz", NewAnzSource(
		NewSourceParams{
			Domain: "https://www.anz.com.au",
		},
	))
}

func GetRegistry() SourceRegistry {
	return registry
}
