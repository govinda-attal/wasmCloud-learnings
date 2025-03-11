package main

import (
	"context"
	"log"
	"strings"

	sdk "go.wasmcloud.dev/provider"
	wrpc "wrpc.io/go"

	"k8s.io/apimachinery/pkg/api/errors"

	// Generated bindings from the wit world
	store "github.com/govinda-attal/wasmCloud-learnings/providers/kv-kube-cfgmaps/bindings/exports/wrpc/keyvalue/store"
)

type Provider struct {
	*sdk.WasmcloudProvider
	cmSrv interface {
		GetValue(ctx context.Context, namespace, name, key string) (string, error)
		Exists(ctx context.Context, namespace, name, key string) (bool, error)
		ListKeys(ctx context.Context, namespace, name string) ([]string, error)
	}
}

var _ store.Handler = (*Provider)(nil)

func (p *Provider) Get(ctx context.Context, bucket string, key string) (*wrpc.Result[[]uint8, store.Error], error) {
	namespace, cm := namespacedConfigMap(bucket)
	log.Println("namespace", namespace, "cm", cm)
	value, err := p.cmSrv.GetValue(ctx, namespace, cm, key)
	if err != nil {
		log.Println("error getting value", err)
		if errors.IsNotFound(err) {
			return nil, store.NewErrorNoSuchStore()
		}
		return nil, store.NewErrorOther(err.Error())
	}
	valBytes := []byte(value)
	log.Println("value", string(valBytes))
	return Ok(valBytes), nil
}

func (p *Provider) Exists(ctx context.Context, bucket string, key string) (*wrpc.Result[bool, store.Error], error) {
	namespace, cm := namespacedConfigMap(bucket)
	exists, err := p.cmSrv.Exists(ctx, namespace, cm, key)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, store.NewErrorNoSuchStore()
		}
		return nil, store.NewErrorOther(err.Error())
	}
	log.Println("exists", exists)
	return Ok(exists), nil
}

func (p *Provider) ListKeys(ctx context.Context, bucket string, cursor *uint64) (*wrpc.Result[store.KeyResponse, store.Error], error) {
	namespace, cm := namespacedConfigMap(bucket)
	keys, err := p.cmSrv.ListKeys(ctx, namespace, cm)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, store.NewErrorNoSuchStore()
		}
		return nil, store.NewErrorOther(err.Error())
	}
	keyRs := store.KeyResponse{
		Keys: keys,
	}
	log.Println("keys", keys)
	return Ok(keyRs), nil
}

func (*Provider) Set(ctx__ context.Context, bucket string, key string, value []uint8) (*wrpc.Result[struct{}, store.Error], error) {
	return nil, store.NewErrorOther("not applicable")
}

func (*Provider) Delete(ctx__ context.Context, bucket string, key string) (*wrpc.Result[struct{}, store.Error], error) {
	return nil, store.NewErrorOther("not applicable")
}

func (p *Provider) handleHealthCheck() string {
	p.Logger.Info("performing health check")
	return "provider is healthy"
}

func (p *Provider) handleShutdown() error {
	p.Logger.Info("shutting down provider")
	return nil
}

func namespacedConfigMap(bucket string) (namespace, cm string) {

	parts := strings.Split(bucket, "/")
	if len(parts) != 2 {
		return "default", bucket
	}
	return parts[0], parts[1]
}


func Ok[T any](v T) *wrpc.Result[T, store.Error] {
	return wrpc.Ok[store.Error](v)
}