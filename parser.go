package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// YamlDevice represents a device in the yaml configuration file
type YamlDevice struct {
	SysName      string `yaml:"sysname"`
	OnConnect    string `yaml:"on_connect"`
	OnDisconnect string `yaml:"on_disconnect"`
}

func loadYamlDevices(filename string) (*map[string]map[string]YamlDevice, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		ErrLog.Println("Could not read file ", filename, err)
		return nil, err
	}

	result := make(map[string]map[string]YamlDevice)
	if err := yaml.Unmarshal(yamlFile, &result); err != nil {
		ErrLog.Println("Could not load file ", filename, err)
		return nil, err
	}
	return &result, nil
}
