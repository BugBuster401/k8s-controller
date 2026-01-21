package main

import (
	"log"
	"net/http"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error getting cluster configuration: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	k8sClient := NewK8sClient(clientset, "test")
	handler := NewHandler(k8sClient)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /workers", handler.CreateWorkers)
	mux.HandleFunc("GET /workers/{name}", handler.GetWorkers)
	mux.HandleFunc("DELETE /workers/{name}", handler.DeleteWorkers)

	log.Println("Server start...")

	log.Fatal(http.ListenAndServe(":9090", mux))
}
