package store

import (
	"os"
	"path/filepath"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/santhosh-tekuri/jsonschema/v5"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// SaveYamlFile saves a yaml file
func SaveYamlFile(
	content interface{},
	filepath string,
) error {
	// serialise the config to yaml
	contents, err := yaml.Marshal(content)
	if core.AssertErrorToNilf("could not marshal config: %w", err) {
		return err
	}
	WriteFile(filepath, contents)

	log.Info("saved: ", filepath)

	return nil
}

func WriteFile(filepath string, contents []byte) (err error) {
	EnsureStoragePath(filepath)

	err = os.WriteFile(filepath, contents, 0644)
	if err != nil {
		return err
	}
	return nil
}

// LoadYamlFile loads a yaml file
func LoadYamlFile[T any](
	path string,
	schema *jsonschema.Schema,
) (T, error) {
	path, err := filepath.Abs(path)
	var fileJson interface{}
	var output T

	if core.AssertErrorToNilf("could not get absolute path: %w", err) {
		return output, err
	}

	content, err := os.ReadFile(path)
	if core.AssertErrorToNilf("could not read file: %w", err) {
		return output, err
	}

	err = yaml.Unmarshal(content, &fileJson)
	if core.AssertErrorToNilf("could not unmarshal file: "+path+" [ %w ]", err) {
		return output, err
	}

	err = schema.Validate(fileJson)
	if core.AssertErrorToNilf("could not validate file: "+path+" [ %s ]", err) {
		return output, err
	}

	err = yaml.Unmarshal(content, &output)
	if core.AssertErrorToNilf("could not unmarshal file: "+path+" [ %w ]", err) {
		return output, err
	}

	log.Info("loaded: ", path)

	return output, nil
}

func EnsureStoragePath(path string) error {
	// get the folder path
	folderPath := filepath.Dir(path)
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
