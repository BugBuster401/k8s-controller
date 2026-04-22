package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"time"
)

func main() {
	workerID := os.Getenv("HOSTNAME")
	taskNumber := os.Getenv("TASK_NUMBER")
	log.Printf("Worker %s sarted for %s", workerID, taskNumber)

	if rand.IntN(100) < 50 { // 50%
		OOMKilled()
	} else {
		JobLogics()
	}
}

func JobLogics() {
	time.Sleep(time.Second * 10)
	log.Println("task completed!!!")
}

func DeploymentLogics() {
	mux := http.NewServeMux()
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func OOMKilled() {
	fmt.Println("Starting memory allocation test...")

	// trying to allocate 1 GiB (1073741824 byte)
	const size = 1073741824
	bigSlice := make([]byte, size)

	// fill it up so that the memory is actually fixed (touch pages)
	for i := range bigSlice {
		bigSlice[i] = 0xFF
	}

	fmt.Println("Memory allocated, sleeping...")
	time.Sleep(1 * time.Minute) // let's see OOMKilled
}
