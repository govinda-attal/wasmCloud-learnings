//go:generate wit-bindgen-go generate --world clusters --out gen ./wit

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	//"log/slog"

	"go.bytecodealliance.org/cm"
	_ "go.bytecodealliance.org/x/cabi"

	//"go.wasmcloud.dev/component/log/wasilog"
	"go.wasmcloud.dev/component/net/wasihttp"

	clusterapi "github.com/govinda-attal/wasmCloud-learnings/clusters/gen/cloud-platform/clusters/cluster-api"
	"github.com/govinda-attal/wasmCloud-learnings/clusters/gen/wrpc/keyvalue/store"
)

func init() {
	clusterapi.Exports.GetClusterInfo = GetClusterInfo
	wasihttp.HandleFunc(handleRequest)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	rs := GetClusterInfo()
	if rs.IsErr() {
		http.Error(w, rs.Err().String(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rs.OK())
}

func GetClusterInfo() getClusterInfoResult {
	//logger := slog.New(wasilog.DefaultOptions().NewHandler())
	//logger.Info("request received")

	tier := "dev"
	cluster := "cluster1"
	data, err := getClusterInfo(tier, cluster)
	if err != nil {
		return cm.Err[getClusterInfoResult](clusterapi.ErrorOther(err.Error()))
	}
	callRs := cm.OK[getClusterInfoResult](data)
	return callRs
}

func getClusterInfo(tier, cluster string) (clusterapi.Data, error) {
	bucket := fmt.Sprintf("default/%s-clusters", tier)
	rs := store.Get(bucket, cluster)
	if rs.IsErr() {
		return clusterapi.Data{}, fmt.Errorf("error getting cluster info: %s", rs.Err())
	}
	rawData := rs.OK().Some().Slice()
	var data clusterapi.Data
	if err := json.Unmarshal(rawData, &data); err != nil {
		return clusterapi.Data{}, err
	}
	return data, nil
}

func main() {}

type getClusterInfoResult = cm.Result[clusterapi.DataShape, clusterapi.Data, clusterapi.Error]

