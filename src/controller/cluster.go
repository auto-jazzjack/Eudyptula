package controller

import "net/http"

const REVIVE = "/api/v1/cluster/revive"

type Cluster struct {
}

func NewCluster() *Cluster {
	return &Cluster{}
}

func (c *Cluster) ServeHTTP(res http.ResponseWriter, req *http.Request) {

}
