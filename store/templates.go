package store

import (
	"bytes"
	"html/template"
	"time"

	"github.com/gosimple/slug"
)

// FilenameTemplateContext is the context used to render a filename template
// the configuiation for the template is stored in the config file under the
// key `outputTemplate`
type FilenameTemplateContext struct {
	Source     string
	SourceSlug string
	Account    FilenameTemplateContextAccount
	DateRange  FilenameTemplateContextDateRange
}
type FilenameTemplateContextAccount struct {
	Name       string
	NameSlug   string
	Number     string
	NumberSlug string
}
type FilenameTemplateContextDateRange struct {
	From     time.Time
	FromUnix int64
	FromSlug string
	To       time.Time
	ToUnix   int64
	ToSlug   string
}

func NewFilenameTemplateContext(
	source string,
	accountName string,
	accountNumber string,
	fromDate time.Time,
	toDate time.Time,
) *FilenameTemplateContext {
	return &FilenameTemplateContext{
		Source:     source,
		SourceSlug: slug.Make(source),
		Account: FilenameTemplateContextAccount{
			Name:       accountName,
			NameSlug:   slug.Make(accountName),
			Number:     accountNumber,
			NumberSlug: slug.Make(accountNumber),
		},
		DateRange: FilenameTemplateContextDateRange{
			From:     fromDate,
			FromUnix: fromDate.Unix(),
			FromSlug: slug.Make(fromDate.Format(time.RFC3339)),
			To:       toDate,
			ToUnix:   toDate.Unix(),
			ToSlug:   slug.Make(toDate.Format(time.RFC3339)),
		},
	}
}

func (f *FilenameTemplateContext) Render(template string) string {
	return ""
}

type FilenameTemplate struct {
	Template *template.Template
}

func NewFilenameTemplate(content string) *FilenameTemplate {
	tmpl, err := template.New("filename").Parse(content)
	if err != nil {
		panic(err)
	}

	return &FilenameTemplate{
		Template: tmpl,
	}
}

func (f *FilenameTemplate) Render(context *FilenameTemplateContext) string {
	var doc bytes.Buffer
	err := f.Template.Execute(&doc, context)
	if err != nil {
		panic(err)
	}
	return doc.String()
}
