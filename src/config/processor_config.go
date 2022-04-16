package config

import (
	"errors"
	"fmt"
	"go-ka/logic"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v3"
)

var printerMapping = map[string]logic.Logic[any]{
	"logic.Printer": logic.Logic[any](logic.NewPrinter()),
}

type ProcessorConfigs[V any] struct {
	Processors map[string]ProcessorConfig[V] `yaml:"Processors"`
	Zookeeper  []string                      `yaml:"Zookeeper"`
}

type ProcessorConfig[V any] struct {
	BoostrapServer string         `yaml:"BoostrapServer"`
	GroupId        string         `yaml:"GroupId"`
	Offset         string         `yaml:"Offset"`
	Topic          string         `yaml:"Topic"`
	Concurrency    int32          `yaml:"Concurrency"`
	PollTimeout    int            `yaml:"PollTimeout"`
	LogicContainer LogicContainer `yaml:"LogicContainer"`
	FetchSize      int32          `yaml:"FetchSize"`
	UserName       string         `yaml:"UserName"`
	Password       string         `yaml:"Password"`
	Algorithm      string         `yaml:"Algorithm"`
}

type LogicContainer struct {
	Logic logic.Logic[any]
}

func NewProcessConfigs[V any]() *ProcessorConfigs[V] {
	path, _ := os.Getwd()
	yamlFile, err := ioutil.ReadFile(path + "/application.yaml")
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

func validate(p ProcessorConfig[any]) error {

	if p.BoostrapServer == "" {
		return errors.New("Bootstrap server should be set")
	}
	if p.GroupId == "" {
		return errors.New("GroupId server should be set")
	}
	if p.Concurrency <= 0 {
		return errors.New("Concurrency server should be postive")
	}
	if p.PollTimeout <= 0 {
		return errors.New("PollTimeout server should be postive")
	}
	if p.UserName != "" || p.Password != "" || p.Algorithm != "" {
		if p.UserName == "" {
			return errors.New("GroupId server should be set")
		}
		if p.Password == "" {
			return errors.New("return server should be set")
		}
		if p.Algorithm == "" {
			return errors.New("Algorithm server should be set")
		}
	}
	return nil
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
