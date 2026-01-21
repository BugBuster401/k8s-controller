package main

import (
	"context"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sClient struct {
	clientset *kubernetes.Clientset
	namespace string
}

func NewK8sClient(clientset *kubernetes.Clientset, namespace string) (*K8sClient, error) {
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

func (c *K8sClient) CreateDeployment(taskNumber, deploymentName string) error {
	// Environment variables
	envVars := []corev1.EnvVar{
		{
			Name:  "TASK_NUMBER",
			Value: taskNumber,
		},
	}

	// Defining Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: c.namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2), // Two consumers
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "k8s-worker", "number": taskNumber},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "k8s-worker", "number": taskNumber},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "k8s-worker",
							Image: "k8s-worker:1.0",
							Env:   envVars, // Passing environment variables
						},
					},
				},
			},
		},
	}

	// Create Deployment
	deploymentsClient := c.clientset.AppsV1().Deployments(c.namespace)
	_, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func int32Ptr(i int32) *int32 { return &i }

func (c *K8sClient) GetDeployment(deploymentName string) (*appsv1.Deployment, error) {
	ctx := context.Background()

	deployment, err := c.clientset.AppsV1().Deployments(c.namespace).
		Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (c *K8sClient) DeleteDeployment(deploymentName string) error {
	ctx := context.Background()

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	log.Printf("Deleting deployment %s in namespace %s...\n",
		deploymentName, c.namespace)

	err := c.clientset.AppsV1().Deployments(c.namespace).
		Delete(ctx, deploymentName, deleteOptions)
	if err != nil {
		return err
	}

	log.Println("Deployment deleted successfully")
	return nil
}
