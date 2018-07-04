package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Parse(filepath string) (*Configuration, error) {
	fileContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return parseFromBytes(fileContent)
}

func parseFromBytes(bytes []byte) (*Configuration, error) {
	var conf Configuration

	err := yaml.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
