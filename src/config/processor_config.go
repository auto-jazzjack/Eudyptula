package config

import (
	"fmt"
	yamlToJson "gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
)

type ProcessorConfigs struct {
	Processors map[string]ProcessorConfig `yaml:"Processors"`
}

type ProcessorConfig struct {
	BoostrapServer string `yaml:"BoostrapServer"`
	GroupId        string `yaml:"GroupId"`
	Offset         string `yaml:"Offset"`
	Topic          string `yaml:"Topic"`
	Concurrency    int    `yaml:"Concurrency"`
	PollTimeout    int    `yaml:"PollTimeout"`
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
