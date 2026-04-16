package main

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sClient struct {
	clientset *kubernetes.Clientset
	namespace string
}

func NewK8sClient(namespace string) (*K8sClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error getting cluster configuration: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, err
	}

	log.Printf("Kubernetes API версия: %s\n", version.String())

	return &K8sClient{
		clientset: clientset,
		namespace: namespace,
	}, nil
}

func int32Ptr(i int32) *int32 { return &i }
