package config

import (
	"fmt"
	"go-ka/logic"
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
)

var printerMapping = map[string]logic.Logic[any]{
	"logic.Printer": logic.Logic[any](logic.NewPrinter()),
}

type ProcessorConfigs[V any] struct {
	Processors map[string]ProcessorConfig[V] `yaml:"Processors"`
}

type ProcessorConfig[V any] struct {
	BoostrapServer string         `yaml:"BoostrapServer"`
	GroupId        string         `yaml:"GroupId"`
	Offset         string         `yaml:"Offset"`
	Topic          string         `yaml:"Topic"`
	Concurrency    int            `yaml:"Concurrency"`
	PollTimeout    int            `yaml:"PollTimeout"`
	LogicContainer LogicContainer `yaml:"LogicContainer"`
	FetchSize      int32          `yaml:"FetchSize"`
}

type LogicContainer struct {
	Logic logic.Logic[any]
}

func NewProcessConfigs[V any]() *ProcessorConfigs[V] {
	yamlFile, err := ioutil.ReadFile("./src/application.yaml")
	if err != nil {
		fmt.Println(err)
		panic("yamlFile.Get err")
	}

	var v = &ProcessorConfigs[V]{}

	//var json, err1 = yamlToJson.YAMLToJSON(yamlFile)
	//vv := reflect.TypeOf("sample.simple_map").Elem()
	//fmt.Println(vv)
	fmt.Println(string(yamlFile))
	err2 := yaml.Unmarshal(yamlFile, v)
	if err2 != nil {
		panic(err2)
	}

	return v
}

func (target *LogicContainer) UnmarshalYAML(value *yaml.Node) error {

	//reflect.New(value.Value)
	//v = string(target)
	//a := reflect.ValueOf(value.Value)
	target.Logic = logic.Logic[any](logic.Logic[any](printerMapping[value.Value]))
	if target.Logic == nil {
		panic("No such LogicContainer")
	}
	return nil
}
