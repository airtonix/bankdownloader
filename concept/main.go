package main

import (
	"errors"
	"fmt"

	"github.com/airtonix/bank-downloaders/core"
	"gopkg.in/yaml.v3"
)

var yaml_data = []byte(`
name: something
jobs:
- source: anz
  credentials:
    username: user
    password: pass

- source: commbank
  credentials:
    ident: user
    password: pass
    code: code
`)

type ConfigFile struct {
	Name string          `json:"name"`
	Jobs []ConfigFileJob `json:"jobs"`
}

type ConfigFileJob struct {
	SourceName        string `json:"source"`
	SourceCredentials any    `json:"credentials"`
	Source            Source // the actual source once we unmarshal
}

func (job *ConfigFileJob) UnmarshalYAML(value *yaml.Node) error {
	var raw interface{}
	if err := value.Decode(&raw); err != nil {
		return err
	}
	sourceName := raw.(map[string]interface{})["source"].(string)
	job.SourceName = sourceName

	source, err := GetSourceFactory(sourceName)
	core.AssertErrorToNilf("could not get source factory", err)

	job.Source = source
	job.Source.Load(value)

	return nil
}

// Source
type SourceProps struct {
	Name string
}
type Source interface {
	Load(node *yaml.Node) error
	Render() error
}

// ANZ
type AnzConfig struct {
	Credentials struct {
		Username string
		Password string
	}
}
type AnzSource struct {
	Config AnzConfig
	*SourceProps
}

func (source *AnzSource) Load(node *yaml.Node) error {
	var config AnzConfig
	if err := node.Decode(&config); err != nil {
		return err
	}
	source.Config = config
	return nil
}
func (source *AnzSource) Render() error {
	fmt.Println(source.Config.Credentials)
	return nil
}

// Commbank
type CommbankConfig struct {
	Credentials struct {
		Ident    string
		Password string
		Code     string
	}
}
type CommbankSource struct {
	Config CommbankConfig
	*SourceProps
}

func (source *CommbankSource) Load(node *yaml.Node) error {
	var config CommbankConfig
	if err := node.Decode(&config); err != nil {
		return err
	}
	source.Config = config
	return nil
}
func (source *CommbankSource) Render() error {
	fmt.Println(source.Config.Credentials)
	return nil
}

func GetSourceFactory(source string) (Source, error) {
	switch source {
	case "anz":
		return &AnzSource{}, nil
	case "commbank":
		return &CommbankSource{}, nil
	default:
		return nil, errors.New("unsupported source")
	}
}

func main() {
	var config ConfigFile
	if err := yaml.Unmarshal(yaml_data, &config); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", config)

	// print the type of job 0
	fmt.Printf("%T\n", config.Jobs[0].Source)
	fmt.Printf("%T\n", config.Jobs[0].Source.Render())

	// print the type of job 1
	fmt.Printf("%T\n", config.Jobs[1].Source)
	fmt.Printf("%T\n", config.Jobs[1].Source.Render())

}
