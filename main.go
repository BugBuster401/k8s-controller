package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	mux.HandleFunc("POST /workers", handleCreateWorkers)
	mux.HandleFunc("GET /workers/{name}", handleGetWorkers)
	mux.HandleFunc("DELETE /workers/{name}", handleDeleteWorkers)

	return http.ListenAndServe(":9090", mux)
}

func handleGetWorkers(w http.ResponseWriter, r *http.Request) {
	deploymentName := r.PathValue("name")

	deployment := getDeployment(clientset, "default", deploymentName)

	buf, err := json.Marshal(deployment)
	if err != nil {
		log.Printf("failed encode json: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
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

	res, err := createDeployment(req.TaskNumber)
	if err != nil {
		log.Printf("failed create deployment: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func handleDeleteWorkers(w http.ResponseWriter, r *http.Request) {
	deploymentName := r.PathValue("name")

	err := deleteDeployment(clientset, "default", deploymentName)
	if err != nil {
		log.Printf("failed delete deployment: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
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

func createDeployment(taskNumber int) (string, error) {
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

func getDeployment(clientset *kubernetes.Clientset, namespace, deploymentName string) *appsv1.Deployment {
	ctx := context.Background()

	// Получение конкретного Deployment
	deployment, err := clientset.AppsV1().Deployments(namespace).
		Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Error getting deployment: %v\n", err)
		return nil
	}

	// Вывод информации о Deployment
	fmt.Printf("Deployment Name: %s\n", deployment.Name)
	fmt.Printf("Namespace: %s\n", deployment.Namespace)
	fmt.Printf("Labels: %v\n", deployment.Labels)
	fmt.Printf("Replicas: %d\n", *deployment.Spec.Replicas)
	fmt.Printf("Available Replicas: %d\n", deployment.Status.AvailableReplicas)
	fmt.Printf("Creation Timestamp: %s\n", deployment.CreationTimestamp)

	// // Вывод информации о контейнерах
	// for _, container := range deployment.Spec.Template.Spec.Containers {
	// 	fmt.Printf("Container: %s, Image: %s\n",
	// 		container.Name, container.Image)
	// }

	// return deployment.Spec.Template.Spec.Containers

	return deployment
}

func deleteDeployment(clientset *kubernetes.Clientset, namespace, deploymentName string) error {
	ctx := context.Background()

	// Опции удаления
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	fmt.Printf("Deleting deployment %s in namespace %s...\n",
		deploymentName, namespace)

	err := clientset.AppsV1().Deployments(namespace).
		Delete(ctx, deploymentName, deleteOptions)
	if err != nil {
		fmt.Printf("Error deleting deployment: %v\n", err)
		return err
	}

	fmt.Println("Deployment deleted successfully")
	return nil
}
