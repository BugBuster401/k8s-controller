IMAGE_NAME = k8s-controller
TAG = 1.0
WORKER_IMAGE_NAME = k8s-worker
WORKER_TAG = 1.0

build:
		docker build -t $(IMAGE_NAME):$(TAG) .

build-worker:
		docker build -f worker/Dockerfile -t $(WORKER_IMAGE_NAME):$(WORKER_TAG) worker/.

push: 
		minikube image load $(IMAGE_NAME):$(TAG)

push-worker: 
		minikube image load $(WORKER_IMAGE_NAME):$(WORKER_TAG)

release: build build-worker push push-worker

.PHONY: build build-worker push push-worker release