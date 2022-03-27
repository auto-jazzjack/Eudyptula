package config

import (
	"go.uber.org/dig"
)

var c = dig.New()

type Container struct {
	c *dig.Container
}

func NewContainer() *Container {
	return &Container{
		c: dig.New(),
	}
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
