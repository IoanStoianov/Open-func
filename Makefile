# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

MAINBIN=cmd/open-func/main.go

CONTAINERNAME=open-func

run:
	@sh -c "trap '$(GORUN) $(MAINBIN)' EXIT"
	
build:
	@eval $$(minikube -p minikube docker-env); docker build . -t $(CONTAINERNAME)

test:
	@$(GOTEST) ./...