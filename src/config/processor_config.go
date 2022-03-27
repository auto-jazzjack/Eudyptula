package config

import (
	"fmt"
	yamlToJson "github.com/ghodss/yaml"
	"io/ioutil"
)

type ProcessorConfigs struct {
	Processors map[string]ProcessorConfig
}

type ProcessorConfig struct {
	BoostrapServer string
	GroupId        string
	Offset         string
	Topic          string
	Concurrency    int
	PollTimeout    int
}

func NewProcessConfigs() *ProcessorConfigs {
	yamlFile, err := ioutil.ReadFile("./application.yaml")
	if err != nil {
		panic("yamlFile.Get err")
	}

	var v = &ProcessorConfigs{}

	//var json, err1 = yamlToJson.YAMLToJSON(yamlFile)
	err2 := yamlToJson.Unmarshal(yamlFile, v)
	if err2 != nil {
		return nil
	}

	fmt.Println(v)
	return v
}
