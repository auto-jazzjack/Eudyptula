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
	c.Provide(controller.NewCluster)

	c.Invoke(func(cluster *controller.Cluster) {
		http.HandleFunc(controller.REVIVE, cluster.ServeHTTP)
	})

	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		return
	}
}
