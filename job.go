package main

import (
	"context"
	"fmt"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *K8sClient) CreateJob(taskNumber, jobName string) error {
	// Environment variables
	envVars := []corev1.EnvVar{
		{
			Name:  "TASK_NUMBER",
			Value: taskNumber,
		},
	}

	// Defining Deployment
	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: c.namespace,
		},
		Spec: batchv1.JobSpec{
			Completions:  int32Ptr(2),
			Parallelism:  int32Ptr(2),
			BackoffLimit: int32Ptr(0),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "k8s-worker", "number": taskNumber},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  "k8s-worker",
							Image: "k8s-worker:1.0",
							Env:   envVars, // Passing environment variables
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("100m"),
									corev1.ResourceMemory: resource.MustParse("128Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("256Mi"),
								},
							},
						},
					},
				},
			},
		},
	}

	jobClient := c.clientset.BatchV1().Jobs(c.namespace)
	job, err := jobClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	go c.watchJob(job.Name)

	return nil
}

func (c *K8sClient) watchJob(jobName string) {
	// Create a Watcher to track changes to a specific Job
	watch, err := c.clientset.BatchV1().Jobs(c.namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", jobName),
	})
	if err != nil {
		log.Printf("failed to create watch: %v", err)
		return
	}
	defer watch.Stop()

	log.Printf("starting to monitor Job: %s", jobName)

	for event := range watch.ResultChan() {
		job, ok := event.Object.(*batchv1.Job)
		if !ok {
			continue
		}

		// Checking the Job Completion Conditions
		for _, condition := range job.Status.Conditions {
			if condition.Status == corev1.ConditionTrue {
				switch condition.Type {
				case batchv1.JobComplete:
					log.Printf("Job %s completed successfully", jobName)
					if err := c.DeleteJob(jobName); err != nil {
						log.Printf("failed delete job %s: %v", jobName, err)
					}
					return // Everything is fine, the application can continue to work.

				case batchv1.JobFailed:
					// Handle job errors
					log.Printf("job %s failed with error: %s. Reason: %s", jobName, condition.Message, condition.Reason)
					if err := c.DeleteJob(jobName); err != nil {
						log.Printf("failed delete job %s: %v", jobName, err)
					}

					return
				}
			}
		}
	}
}

func (c *K8sClient) DeleteJob(jobName string) error {
	ctx := context.Background()

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	log.Printf("Deleting job %s in namespace %s...\n",
		jobName, c.namespace)

	err := c.clientset.BatchV1().Jobs(c.namespace).
		Delete(ctx, jobName, deleteOptions)
	if err != nil {
		return err
	}

	log.Println("Job deleted successfully")
	return nil
}
