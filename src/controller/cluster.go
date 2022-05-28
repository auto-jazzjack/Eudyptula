package controller

import (
	"encoding/json"
	"fmt"
	"go-ka/consumer"
	"net/http"
	"strings"
)

const REVIVE = "/api/v1/cluster/revive"
const REWIND = "/api/v1/cluster/rewind/"

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
	} else if strings.HasPrefix(req.URL.Path, REWIND) {

		date := req.URL.Query().Get("date")
		processor := req.URL.Query().Get("processor")
		tmp, err2 := c.manager.Rewind(processor, date)
		if err2 != nil {
			res.WriteHeader(500)
			res.Write([]byte(fmt.Sprint(err2)))
		}
		result, err = json.Marshal(tmp)
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
