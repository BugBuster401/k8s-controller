package main

import (
	"log"
	"net/http"
)

func main() {
	k8sClient, err := NewK8sClient("test")
	if err != nil {
		log.Fatalf("Failed k8s api connection ")
	}
	handler := NewHandler(k8sClient)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /workers", handler.CreateWorkers)
	mux.HandleFunc("GET /workers/{name}", handler.GetWorkers)
	mux.HandleFunc("DELETE /workers/{name}", handler.DeleteWorkers)
	mux.HandleFunc("POST /jobs", handler.CreateJob)

	log.Println("Server start...")

	log.Fatal(http.ListenAndServe(":9090", mux))
}
