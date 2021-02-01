# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

MAINBIN=cmd/open-func/main.go

run:
	@sh -c "trap '$(GORUN) $(MAINBIN)' EXIT"