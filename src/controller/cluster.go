package controller

import (
	"encoding/json"
	"fmt"
	"go-ka/consumer"
	"net/http"
)

const REVIVE = "/api/v1/cluster/revive"
const REWIND = "/api/v1/cluster/rewind"

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

	var result []byte
	var err error
	if req.URL.Path == REVIVE {
		result, err = json.Marshal(c.manager.ExecuteAll())
	} else if req.URL.Path == REWIND {
		c.manager.ExecuteAll()
	}

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
