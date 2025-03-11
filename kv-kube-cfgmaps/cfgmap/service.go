package cfgmap

import (
	"context"
	"fmt"
	"log"
	"maps"
	"slices"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type Service struct {
	core clientcorev1.CoreV1Interface
}

func NewService() (*Service, error) {
	core, err := kubeClient()
	if err != nil {
		return nil, err
	}

	return &Service{
		core: core,
	}, nil
}

func (s *Service) GetValue(ctx context.Context, namespace, name, key string) (string, error) {
	var value string
	log.Println("namespace", namespace, "name", name, "key", key)
	cm, err := s.core.ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Println("error getting configmap", err)
		return value, err
	}

	value, ok := cm.Data[key]
	if !ok {
		return value, fmt.Errorf("key not found")
	}
	return value, nil
}

func (s *Service) ListKeys(ctx context.Context, namespace, name string) ([]string, error) {
	cm, err := s.core.ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return slices.Collect(maps.Keys(cm.Data)), nil
}

func (s *Service) Exists(ctx context.Context, namespace, name, key string) (bool, error) {
	if _, err := s.GetValue(ctx, namespace, name, key); err != nil {
		return false, err
	}

	return true, nil
}

func kubeClient() (clientcorev1.CoreV1Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	//loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	//config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, nil).ClientConfig()
	//if err != nil {
	//	return nil, err
	//}

	kubeClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return kubeClientset.CoreV1(), nil
}
