package main

import (
	"controller"
	"net/http"
)

func main() {
	http.HandleFunc(controller.REVIVE, controller.NewCluster().ServeHTTP)
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		return
	}
}
