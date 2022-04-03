package controller

import (
	"encoding/json"
	"fmt"
	"go-ka/consumer"
	"net/http"
)

const REVIVE = "/api/v1/cluster/revive"

type Cluster[V any] struct {
	manager *consumer.Manager[V]
}

func NewCluster[V any](manager *consumer.Manager[V]) *Cluster[V] {
	return &Cluster[V]{
		manager: manager,
	}
}

/**
Actually request is usless for this phase
*/
func (c *Cluster[V]) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	result, err := json.Marshal(c.manager.ExecuteAll())
	if err != nil {
		res.WriteHeader(500)
		_, err := res.Write([]byte(fmt.Sprint(err)))
		if err != nil {
			return
		}
	} else {
		res.WriteHeader(200)
		_, err := res.Write(result)
		if err != nil {
			return
		}
	}

}
