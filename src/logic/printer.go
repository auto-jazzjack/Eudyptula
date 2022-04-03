package logic

import (
	"encoding/json"
	"fmt"
)

type Printer[V any] struct {
}

func NewPrinter[V any]() Logic[V] {
	return Printer[V]{}
}
func (p Printer[V]) Deserialize(bytes []byte) *V {
	retv := p.DefaultValue()
	err := json.Unmarshal(bytes, &retv)
	if err != nil {
		return nil
	}
	return retv
}

func (p Printer[V]) DoAction(v V) {
	fmt.Println(v)
}

func (p Printer[V]) DefaultValue() *V {
	return nil
}
