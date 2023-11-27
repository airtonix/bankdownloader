package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// it should test that the example config loads and therefore the config shcema is working
func TestCreateCompiler(t *testing.T) {
	compiler := NewSchemaCompiler()
	assert.NotNil(t, compiler)
}

// test that the config schema compiles
func TestRegisterConfigSchema(t *testing.T) {
	compiler := NewSchemaCompiler()
	err := compiler.RegisterConfigSchema()
	assert.Nil(t, err)
}

// test that the history schema compiles
func TestRegisterHistorySchema(t *testing.T) {
	compiler := NewSchemaCompiler()
	err := compiler.RegisterHistorySchema()
	assert.Nil(t, err)
}
