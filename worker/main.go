package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	workerID := os.Getenv("WORKER_ID")
	taskNumber := os.Getenv("TASK_NUMBER")
	log.Printf("Worker %s sarted for %s", workerID, taskNumber)

	mux := http.NewServeMux()

	log.Fatal(http.ListenAndServe(":8080", mux))
}
