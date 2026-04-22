package main

import (
	"log"

	"github.com/BugBuster401/k8s-controller/k8s"
)

type JobOrchestrator struct {
	k8sClient *k8s.K8sClient
}

func NewJobOrchestrator(k8sClient *k8s.K8sClient) *JobOrchestrator {
	return &JobOrchestrator{
		k8sClient: k8sClient,
	}
}

func (o *JobOrchestrator) StartEventProcessor() {
	go func() {
		for event := range o.k8sClient.GetEventChannel() {
			switch event.Status {
			case k8s.JobCompleted:
				log.Printf("%s %s", event.JobName, event.Message)
			case k8s.JobFailed:
				log.Printf("%s job failed with error: %s. Reason: %s", event.JobName, event.Message, event.Error)
			}

			if err := o.k8sClient.DeleteJob(event.JobName); err != nil {
				log.Printf("failed delete job %s: %v", event.JobName, err)
			}
		}
	}()
}
