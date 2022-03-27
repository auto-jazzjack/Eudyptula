package config

import (
	"go.uber.org/dig"
	"sync"
)

var cached *Container
var lock = &sync.Mutex{}

type Container struct {
	c *dig.Container
}

func NewContainer() *Container {

	lock.Lock()
	if cached == nil {
		cached = &Container{
			c: dig.New(),
		}
	}
	lock.Unlock()
	return cached

}

func (con *Container) Provide(constructor interface{}, opts ...dig.ProvideOption) {
	err := con.c.Provide(constructor, opts...)
	if err != nil {
		panic(err)
	}
}

func (con *Container) Invoke(function interface{}, opts ...dig.InvokeOption) {
	err := con.c.Invoke(function, opts...)
	if err != nil {
		panic(err)
	}
}
