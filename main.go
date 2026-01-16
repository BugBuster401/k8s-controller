package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /workers", handleGetWorkers)
	mux.HandleFunc("POST /workers", handleCreateWorkers)
	mux.HandleFunc("DELETE /workers/{name}", handleDeleteWorkers)

	return http.ListenAndServe(":9090", mux)
}

func handleGetWorkers(w http.ResponseWriter, r *http.Request) {

}

type CreateWorkerRequest struct {
	TaskNumber int `json:"task_number"`
}

func handleCreateWorkers(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("failed decode json: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := createDeploy(req.TaskNumber)
	if err != nil {
		log.Printf("failed create deployment: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func handleDeleteWorkers(w http.ResponseWriter, r *http.Request) {

}

var clientset *kubernetes.Clientset

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error getting cluster configuration: %v", err)
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	log.Fatal(run())
}

func int32Ptr(i int32) *int32 { return &i }

func createDeploy(taskNumber int) (string, error) {
	// Environment variables
	envVars := []corev1.EnvVar{
		{
			Name:  "TASK_NUMBER",
			Value: strconv.Itoa(taskNumber),
		},
	}

	// Defining Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "controller-k8s",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1), // Two consumers
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "controller-k8s", "number": "1"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "controller-k8s", "number": "1"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "worker-k8s",
							Image: "worker-k8s:1.0",
							Env:   envVars, // Passing environment variables
						},
					},
				},
			},
		},
	}

	// Deploy
	deploymentsClient := clientset.AppsV1().Deployments("default")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	return result.GetName(), nil
}

func deleteDeploy()
