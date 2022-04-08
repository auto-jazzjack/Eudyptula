package logic

import (
	"fmt"
)

type Printer[V string] struct {
}

func NewPrinter[V string]() Logic[string] {
	return Printer[string]{}
}
func (p Printer[string]) Deserialize(bytes []byte) *string {

	v := string(bytes)
	return &v
	//return string(bytes)
}

func (p Printer[string]) DoAction(v string) {
	fmt.Println(v)
}

func (p Printer[string]) DefaultValue() *string {
	return nil
}
