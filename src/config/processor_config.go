package config

import (
	"fmt"
	yamlToJson "github.com/ghodss/yaml"
	"io/ioutil"
	"reflect"
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
	Target         *reflect.StructField
}

func NewProcessConfigs() *ProcessorConfigs {
	yamlFile, err := ioutil.ReadFile("./src/application.yaml")
	if err != nil {
		fmt.Println(err)
		panic("yamlFile.Get err")
	}

	var v = &ProcessorConfigs{}

	//var json, err1 = yamlToJson.YAMLToJSON(yamlFile)
	//vv := reflect.TypeOf("sample.simple_map").Elem()
	//fmt.Println(vv)
	fmt.Println(string(yamlFile))
	err2 := yamlToJson.Unmarshal(yamlFile, v)
	if err2 != nil {
		panic(err2)
	}

	return v
}
