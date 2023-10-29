package config

import (
	_ "embed"

	"encoding/json"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed schema.json
var SchemaJson string
var schema *jsonschema.Schema

func GetSchema() any {
	var schemaObject any
	err := json.Unmarshal([]byte(SchemaJson), &schemaObject)
	core.AssertErrorToNilf("could not unmarshal schema.json: %w", err)

	return schemaObject
}

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
