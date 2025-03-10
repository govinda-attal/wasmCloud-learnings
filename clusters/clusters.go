//go:generate wit-bindgen-go generate --world clusters --out gen ./wit

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"go.bytecodealliance.org/cm"
	_ "go.bytecodealliance.org/x/cabi"

	"go.wasmcloud.dev/component/log/wasilog"

	clusterapi "github.com/govinda-attal/wasmCloud-learnings/clusters/gen/cloud-platform/clusters/cluster-api"
	"github.com/govinda-attal/wasmCloud-learnings/clusters/gen/wrpc/keyvalue/store"
)

func init() {
	clusterapi.Exports.GetClusterInfo = GetClusterInfo
}

func GetClusterInfo() getClusterInfoResult {
	logger := slog.New(wasilog.DefaultOptions().NewHandler())
	logger.Info("request received")

	tier := "dev"
	cluster := "cluster1"
	bucket := fmt.Sprintf("default/%s-clusters", tier)
	rs := store.Get(bucket, cluster)
	if rs.IsErr() {
		return cm.Err[getClusterInfoResult](clusterapi.ErrorOther(rs.Err().String()))
	}
	rawData := rs.OK().Some().Slice()
	var data clusterapi.Data
	if err := json.Unmarshal(rawData, &data); err != nil {
		return cm.Err[getClusterInfoResult](clusterapi.ErrorOther(err.Error()))
	}
	callRs := cm.OK[getClusterInfoResult](data)
	return callRs
}

func main() {}

type getClusterInfoResult = cm.Result[clusterapi.DataShape, clusterapi.Data, clusterapi.Error]
