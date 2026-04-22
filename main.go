package main

import (
	"log"
	"net/http"

	"github.com/BugBuster401/k8s-controller/k8s"
)

func main() {
	k8sClient, err := k8s.NewK8sClient("test")
	if err != nil {
		log.Fatalf("Failed k8s api connection ")
	}

	orchestrator := NewJobOrchestrator(k8sClient)
	orchestrator.StartEventProcessor()

	handler := NewHandler(k8sClient)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /workers", handler.CreateWorkers)
	mux.HandleFunc("GET /workers/{name}", handler.GetWorkers)
	mux.HandleFunc("DELETE /workers/{name}", handler.DeleteWorkers)
	mux.HandleFunc("POST /jobs", handler.CreateJob)
	mux.HandleFunc("POST /jobs/batch", handler.CreateJobs)

	log.Println("Server start...")

	log.Fatal(http.ListenAndServe(":9090", mux))
}
