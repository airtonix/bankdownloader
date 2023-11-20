package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnixDateFormat(t *testing.T) {

	var (
		processorName = "anz"
		accountName   = "My Account"
		accountNumber = "123456789"
		fromDate      = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		toDate        = time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)
		template      = "{{ .SourceSlug }}-{{ .Account.NameSlug }}-{{ .Account.NumberSlug }}-{{ .DateRange.FromUnix  }}-{{ .DateRange.ToUnix }}.csv"
	)

	filenameContext := NewFilenameTemplateContext(
		processorName,
		accountName,
		accountNumber,
		fromDate,
		toDate,
	)

	filenameTemplate := NewFilenameTemplate(template)

	assert.Equal(t,
		"anz-my-account-123456789-1577836800-1580428800.csv",
		filenameTemplate.Render(filenameContext),
		"should render the filename correctly",
	)
}
func TestSlugDateFormat(t *testing.T) {

	var (
		processorName = "anz"
		accountName   = "My Account"
		accountNumber = "123456789"
		fromDate      = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		toDate        = time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)
		template      = "{{ .SourceSlug }}_{{ .Account.NameSlug }}_{{ .Account.NumberSlug }}_{{ .DateRange.FromSlug  }}_{{ .DateRange.ToSlug }}.csv"
	)

	filenameContext := NewFilenameTemplateContext(
		processorName,
		accountName,
		accountNumber,
		fromDate,
		toDate,
	)

	filenameTemplate := NewFilenameTemplate(template)

	assert.Equal(t,
    "anz_my-account_123456789_2020-01-01t00-00-00z_2020-01-31t00-00-00z.csv",
		filenameTemplate.Render(filenameContext),
		"should render the filename correctly",
	)
}
