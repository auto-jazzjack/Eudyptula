package config

import (
	"fmt"
	"go-ka/logic"
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
	"reflect"
)

type ProcessorConfigs[V any] struct {
	Processors map[string]ProcessorConfig[V] `yaml:"Processors"`
}

type ProcessorConfig[V any] struct {
	BoostrapServer string            `yaml:"BoostrapServer"`
	GroupId        string            `yaml:"GroupId"`
	Offset         string            `yaml:"Offset"`
	Topic          string            `yaml:"Topic"`
	Concurrency    int               `yaml:"Concurrency"`
	PollTimeout    int               `yaml:"PollTimeout"`
	Target         LogicContainer[V] `yaml:"Target"`
}

type LogicContainer[V any] struct {
	logic logic.Logic[V]
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

func (target *LogicContainer[V]) UnmarshalYAML(value *yaml.Node) error {

	//reflect.New(value.Value)
	//v = string(target)
	a := reflect.ValueOf(value.Value)
	fmt.Println(value.Value)
	fmt.Println(a)
	target.logic = logic.NewPrinter[V]()
	return nil
}
