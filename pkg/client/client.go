package client

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

// NewMeticsClientSets ...
func NewMeticsClientSets(conf *rest.Config) (*kubernetes.Clientset, *metricsv.Clientset, error) {

	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return nil, nil, err
	}

	metricsClientset, err := metricsv.NewForConfig(conf)
	if err != nil {
		return nil, nil, err
	}

	return clientset, metricsClientset, nil

}
