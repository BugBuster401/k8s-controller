package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	workerID := os.Getenv("HOSTNAME")
	taskNumber := os.Getenv("TASK_NUMBER")
	log.Printf("Worker %s sarted for %s", workerID, taskNumber)

	mux := http.NewServeMux()

	OOMKilled()

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func OOMKilled() {
	fmt.Println("Starting memory allocation test...")

	// Пытаемся выделить 1 GiB (1073741824 байт)
	const size = 1073741824
	bigSlice := make([]byte, size)

	// Заполняем, чтобы память реально зафиксировалась (touch pages)
	for i := range bigSlice {
		bigSlice[i] = 0xFF
	}

	fmt.Println("Memory allocated, sleeping...")
	time.Sleep(1 * time.Minute) // Даем время увидеть OOMKilled
}
