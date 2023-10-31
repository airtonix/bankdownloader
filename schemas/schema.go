package schemas

import (
	_ "embed"
	"strings"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/santhosh-tekuri/jsonschema/v5"
	log "github.com/sirupsen/logrus"
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

//go:embed config.json
var ConfigSchemaJson string
var configSchema *jsonschema.Schema

//go:embed history.json
var HistorySchemaJson string
var historySchema *jsonschema.Schema

func Initialize() {
	var err error
	c := jsonschema.NewCompiler()
	c.RegisterExtension(
		"itemsUniqueProperties",
		itemsUniquePropertiesMeta,
		itemsUniquePropertiessCompiler{},
	)

	err = c.AddResource("config.json", strings.NewReader(ConfigSchemaJson))
	if core.AssertErrorToNilf("could not add schemas/config.json: %w", err) {
		log.Fatal(err)
		return
	}
	configSchema = c.MustCompile("config.json")

	err = c.AddResource("history.json", strings.NewReader(HistorySchemaJson))
	if core.AssertErrorToNilf("could not add schemas/history.json: %w", err) {
		log.Fatal(err)
		return
	}
	historySchema = c.MustCompile("history.json")

}

func GetConfigSchema() *jsonschema.Schema {
	return configSchema
}

func GetHistorySchema() *jsonschema.Schema {
	return historySchema
}
