
# Image URL to use all building/pushing image targets
IMG ?= quay.io/cargaona/image-cloner-controller:webhook
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/manifests
run: generate fmt vet manifests
	go run ./main.go

# Deploy controller in the configured Kubernetes cluster in ~/.kube/manifests

install:
	kubectl apply -f .manifests/deployment.yaml
	kubectl apply -f .manifests/rbac.yaml

uninstall:
	kubectl delete -f .manifests/

deploy:
	kubectl apply -f .manifests/deployment.yaml

delete-deploy:
	kubectl delete -f .manifests/deployment.yaml

install-webhooks:
	kubectl apply -f .manifests/webhooks.yaml

delete-webhooks:
	kubectl delete -f .manifests/webhooks.yaml
# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

registry-secret:
	kubectl create secret generic registrypullsecret --from-file=.dockerconfigjson=$(PWD)/.manifests/imagePullRegistry/configk8s.json --type=kubernetes.io/dockerconfigjson
	kubectl patch serviceaccount default -p '{"imagePullSecrets": [{"name": "registrypullsecret"}]}'

# Build the docker image
docker-build:
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

docker-prune:
	docker system prune --force --all


redeploy: docker-build docker-push docker-prune delete-deploy install

