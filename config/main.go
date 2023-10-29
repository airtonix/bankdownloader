package config

import (
	"strings"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/santhosh-tekuri/jsonschema/v5"
	log "github.com/sirupsen/logrus"
)

func init() {
	c := jsonschema.NewCompiler()
	c.RegisterExtension("itemsUniqueProperties", itemsUniquePropertiesMeta, itemsUniquePropertiessCompiler{})
	err := c.AddResource("schema.json", strings.NewReader(SchemaJson))
	if core.AssertErrorToNilf("could not add schema.json: %w", err) {
		log.Fatal(err)
	}

	schema = c.MustCompile("schema.json")
}
