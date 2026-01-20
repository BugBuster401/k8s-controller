package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	client *K8sClient
}

func NewHandler(client *K8sClient) *Handler {
	return &Handler{client: client}
}

func (h *Handler) GetWorkers(w http.ResponseWriter, r *http.Request) {
	deploymentName := r.PathValue("name")

	deployment, err := h.client.GetDeployment(deploymentName)
	if err != nil {
		log.Printf("failed get deployment: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	TaskNumber string `json:"task_number"`
}

func (h *Handler) CreateWorkers(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("failed decode json: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.client.CreateDeployment(req.TaskNumber, "k8s-worker"); err != nil {
		log.Printf("failed create deployment: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteWorkers(w http.ResponseWriter, r *http.Request) {
	deploymentName := r.PathValue("name")

	err := h.client.DeleteDeployment(deploymentName)
	if err != nil {
		log.Printf("failed delete deployment: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
