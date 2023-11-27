package store

import (
	_ "embed"
	"strings"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/sirupsen/logrus"
)

var itemsUniquePropertiesMeta = jsonschema.MustCompileString("itemsUniqueProperties.json", `{
	"properties": {
	  "itemsUniqueProperties": {
		"type": "array",
		"items": {
		  "type": "string"
		},
		"minItems": 1
	  }
	}
  }`)

type itemsUniquePropertiesSchema []string
type itemsUniquePropertiessCompiler struct{}

func (itemsUniquePropertiessCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {

	if items, ok := m["itemsUniqueProperties"]; ok {
		itemsInterface := items.([]interface{})
		itemsString := make([]string, len(itemsInterface))
		for i, v := range itemsInterface {
			itemsString[i] = v.(string)
		}
		return itemsUniquePropertiesSchema(itemsString), nil
	}

	return nil, nil
}

func (s itemsUniquePropertiesSchema) Validate(ctx jsonschema.ValidationContext, v interface{}) error {
	for _, uniqueProperty := range s {
		items := v.([]interface{})
		seen := make(map[string]bool)
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			if _, ok := itemMap[uniqueProperty]; ok {
				value := itemMap[uniqueProperty].(string)
				if seen[value] {
					return ctx.Error("itemsUniqueProperty", "duplicate %s %s", uniqueProperty, value)
				}
				seen[value] = true
			}
		}
	}
	return nil
}

//go:embed config-schema.json
var ConfigSchemaJson string

//go:embed history-schema.json
var HistorySchemaJson string

var schema *SchemaCompiler

type SchemaCompiler struct {
	compiler      *jsonschema.Compiler
	configSchema  *jsonschema.Schema
	historySchema *jsonschema.Schema
}

func NewSchemaCompiler() *SchemaCompiler {

	c := jsonschema.NewCompiler()
	c.RegisterExtension(
		"itemsUniqueProperties",
		itemsUniquePropertiesMeta,
		itemsUniquePropertiessCompiler{},
	)

	compiler := &SchemaCompiler{
		compiler:      c,
		configSchema:  nil,
		historySchema: nil,
	}
	return compiler
}
func (compiler *SchemaCompiler) RegisterConfigSchema() error {
	var err error
	compiler.configSchema, err = compiler.RegisterSchema(
		"config-schema.json",
		ConfigSchemaJson,
	)
	if err != nil {
		return err
	}
	return nil

}

func (compiler *SchemaCompiler) RegisterHistorySchema() error {
	var err error

	compiler.historySchema, err = compiler.RegisterSchema(
		"history-schema.json",
		HistorySchemaJson,
	)
	if err != nil {
		return err
	}

	return nil
}

// Register a schema with the compiler
func (s *SchemaCompiler) RegisterSchema(
	name string,
	source string,
) (*jsonschema.Schema, error) {
	var output *jsonschema.Schema

	err := s.compiler.AddResource(
		name,
		strings.NewReader(source),
	)

	if core.AssertErrorToNilf("could not add config-schema.json: %w", err) {
		logrus.Fatal(err)
		return nil, err
	}

	output = s.compiler.MustCompile("config-schema.json")

	return output, nil
}

// Public API to bring the schema compiler to life
func InitialiseSchemas() error {
	var err error
	compiler := NewSchemaCompiler()

	err = compiler.RegisterConfigSchema()
	if err != nil {
		return err
	}

	err = compiler.RegisterHistorySchema()
	if err != nil {
		return err
	}

	return nil
}

// Public API to get the schema compiler
func GetSchema() *SchemaCompiler {
	return schema
}
