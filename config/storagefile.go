package config

import (
	"os"
	"path/filepath"

	"github.com/airtonix/bank-downloaders/core"
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
	EnsureStoragePath(filepath)

	err = os.WriteFile(filepath, contents, 0644)
	if err != nil {
		return err
	}

	return nil
}

// LoadYamlFile loads a yaml file
func LoadYamlFile(path string) ([]byte, error) {
	path, err := filepath.Abs(path)
	if core.AssertErrorToNilf("could not get absolute path: %w", err) {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if core.AssertErrorToNilf("could not read file: %w", err) {
		return nil, err
	}

	log.Info("Using file: ", path)
	return content, err
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
