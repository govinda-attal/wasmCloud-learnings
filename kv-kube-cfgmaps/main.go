//go:generate wit-bindgen-wrpc go --out-dir bindings --package github.com/govinda-attal/wasmCloud-learnings/providers/kv-kube-cfgmaps/bindings wit

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	server "github.com/govinda-attal/wasmCloud-learnings/providers/kv-kube-cfgmaps/bindings"
	"github.com/govinda-attal/wasmCloud-learnings/providers/kv-kube-cfgmaps/cfgmap"
	"go.wasmcloud.dev/provider"
)

func main() {

	envs := os.Environ()
	fmt.Println("len of envs", len(envs))
	for _, env := range envs {
		fmt.Println(env)
	}

	filepath := "/var/run/secrets/kubernetes.io/serviceaccount/token"
	if _, err := os.Stat(filepath); err == nil {
		fmt.Printf("File '%s' exists\n", filepath)
	} else if os.IsNotExist(err) {
		fmt.Printf("File '%s' does not exist\n", filepath)
	} else {
		fmt.Printf("Error checking file: %v\n", err)
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	os.Setenv("KUBERNETES_SERVICE_HOST", "10.96.0.1")
	os.Setenv("KUBERNETES_SERVICE_PORT", "443")
	
	cmSrv, err := cfgmap.NewService()
	if err != nil {
		panic(err)
	}
	log.Println("registered cfgmap service")
	// Initialize the provider with callbacks to track linked components
	prv := Provider{
		cmSrv: cmSrv,
	}
	wcPrv, err := provider.New(
		provider.HealthCheck(prv.handleHealthCheck),
		provider.Shutdown(prv.handleShutdown),
	)
	if err != nil {
		panic(err)
	}

	prv.WasmcloudProvider = wcPrv

	providerCh := make(chan error, 1)
	signalCh := make(chan os.Signal, 1)

	log.Println("Starting provider server")

	stopFunc, err := server.Serve(prv.RPCClient, &prv)
	if err != nil {
		wcPrv.Shutdown()
		panic(err)
	}

	go func() {
		err := wcPrv.Start()
		providerCh <- err
	}()

	signal.Notify(signalCh, syscall.SIGINT)

	select {
	case err = <-providerCh:
		stopFunc()
		panic(err)
	case <-signalCh:
		wcPrv.Shutdown()
		stopFunc()
	}

	return nil
}
