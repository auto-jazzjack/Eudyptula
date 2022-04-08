package logic

import (
	"fmt"
)

type Printer struct {
}

func NewPrinter() Logic[string] {
	return Printer{}
}
func (p Printer) Deserialize(bytes []byte) *string {

	v := string(bytes)
	return &v
	//return string(bytes)
}

func (p Printer) DoAction(v string) {
	fmt.Println(v)
}

func (p Printer) DefaultValue() *string {
	return nil
}
