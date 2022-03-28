package subject

import (
	"encoding/json"
	"fmt"
	"go-ka/consumer"
)

type LocalSubject struct {
	Data map[string]string
}

func NewLocalSubject[LocalSubject]() *consumer.Executor[LocalSubject] {
	l := &LocalSubject[LocalSubject]{}
	return l
}

func (LocalSubject[LocalSubject]) DeSerialize(data []byte, target LocalSubject) interface{} {
	err := json.Unmarshal(data, &target)
	if err != nil {
		return nil
	}
	return target
	//TODO implement me
	//	panic("implement me")
}

func (LocalSubject[LocalSubject]) DoAction(data LocalSubject) {
	fmt.Println(data)
}

func (LocalSubject[LocalSubject]) DefaultValue() interface{} {
	return &LocalSubject{}
}
