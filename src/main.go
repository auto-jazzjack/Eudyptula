package main

import (
	"go-ka/config"
	"go-ka/consumer"
	"go-ka/controller"
	"net/http"
)

func main() {

	c := config.NewContainer()
	c.Provide(config.NewProcessConfigs)
	c.Provide(consumer.NewManager[any])
	c.Provide(controller.NewCluster[any])

	c.Invoke(func(cluster *controller.Cluster[any]) {
		http.HandleFunc(controller.REVIVE, cluster.ServeHTTP)
	})

	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		return
	}
}
