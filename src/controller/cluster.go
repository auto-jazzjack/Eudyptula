package controller

import (
	"encoding/json"
	"fmt"
	"go-ka/consumer"
	"net/http"
)

const REVIVE = "/api/v1/cluster/revive"

type Cluster struct {
	manager *consumer.Manager
}

func NewCluster(manager *consumer.Manager) *Cluster {
	return &Cluster{
		manager: manager,
	}
}

/**
Actually request is usless for this phase
*/
func (c *Cluster) ServeHTTP(res http.ResponseWriter, req *http.Request) {

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
