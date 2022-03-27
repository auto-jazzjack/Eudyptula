package main

import (
	"go-ka/config"
	//"go-ka/consumer"
	"go-ka/controller"
	"net/http"
)

func main() {

	c := config.NewContainer()
	c.Provide(config.NewProcessConfigs)
	//c.Provide(consumer.NewManager)

	http.HandleFunc(controller.REVIVE, controller.NewCluster().ServeHTTP)
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		return
	}
}
