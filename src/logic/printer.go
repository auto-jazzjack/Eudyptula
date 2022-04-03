package logic

import (
	"encoding/json"
	"fmt"
	"go-ka/sample"
)

type Printer[V any] struct {
}

func NewPrinter[V any]() *Printer[V] {
	return &Printer[V]{}
}
func (p *Printer[V]) Deserialize(bytes []byte) any {
	retv := p.DefaultValue()
	err := json.Unmarshal(bytes, &retv)
	if err != nil {
		return nil
	}
	return retv
}

func (p *Printer[V]) DoAction(v any) {
	fmt.Println(v)
}

func (p *Printer[V]) DefaultValue() any {
	return sample.Sample{}
}
