package store

type HistoryStrategy interface {
	Strategy() strategy
	ToString() string
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
func (h strategy) ToString() string {
	switch h {
	case DaysAgo:
		return "days-ago"
	case SinceLastDownload:
		return "since-last-download"
	default:
		return "days-ago"
	}
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
