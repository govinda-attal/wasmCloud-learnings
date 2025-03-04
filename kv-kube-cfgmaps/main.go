//go:generate wit-bindgen-wrpc go --out-dir bindings --package github.com/govinda-attal/wasmCloud-learnings/providers/kv-kube-cfgmaps/bindings wit

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	server "github.com/govinda-attal/wasmCloud-learnings/providers/kv-kube-cfgmaps/bindings"
	"github.com/govinda-attal/wasmCloud-learnings/providers/kv-kube-cfgmaps/cfgmap"
	"go.wasmcloud.dev/provider"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
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
