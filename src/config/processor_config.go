package config

import (
	"fmt"
	yamlToJson "gopkg.in/yaml.v2"
	"io/ioutil"
)

type ProcessorConfigs struct {
	Processors map[string]ProcessorConfig `yaml:"Processors"`
}

type ReflectionString string

func NewReflectionString(v string) *ReflectionString {
	var retv ReflectionString
	retv = ReflectionString(v)
	return &retv
}

type ProcessorConfig struct {
	BoostrapServer string           `yaml:"BoostrapServer"`
	GroupId        string           `yaml:"GroupId"`
	Offset         string           `yaml:"Offset"`
	Topic          string           `yaml:"Topic"`
	Concurrency    int              `yaml:"Concurrency"`
	PollTimeout    int              `yaml:"PollTimeout"`
	Target         ReflectionString `yaml:"Target"`
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

func (ut *ReflectionString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	v = string(*ut)
	fmt.Println(v)
	return nil
}
