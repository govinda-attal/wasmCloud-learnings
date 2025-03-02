//go:generate wit-bindgen-go generate --world clusters --out gen ./wit

package main

import (
	clusterapi "github.com/govinda-attal/wasmCloud-learnings/clusters/gen/cloud-platform/clusters/cluster-api"
)

func init() {
	clusterapi.Exports.Get = Get
}

func Get(tier string, cluster string) *clusterapi.Data {
	return nil
}

func main() {}
