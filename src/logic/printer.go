package logic

import (
	"fmt"
	"reflect"
)

type Printer struct {
}

func NewPrinter() Logic[any] {
	return Printer{}
}
func (p Printer) Deserialize(bytes []byte) *any {

	var v interface{}
	v = string(bytes)
	return &v
	//return string(bytes)
}

func (p Printer) DoAction(v any) {
	if reflect.ValueOf(v).Kind() == reflect.Ptr {
		var value = v.(*interface{})
		fmt.Printf("Action %s\n", *value)
	} else {
		fmt.Printf("Action %s\n", v)
	}

}

func (p Printer) DefaultValue() *any {
	return nil
}
