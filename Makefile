ROOT := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
DIST = $(ROOT)/bin/app
GO = go
GOLANGCI_LINT := golangci-lint

DOCKER_COMPOSE := $(or $(DOCKER_COMPOSE),$(DOCKER_COMPOSE),docker-compose)
DOCKER := docker
DOCKER_REPO := amikai
TAG := 1.0

d.lint:
	$(DOCKER) run --mount type=bind,source=$(ROOT)/,target=/workspace -w /workspace golangci/golangci-lint golangci-lint run ./...

lint:
	$(GOLANGCI_LINT) run ./...

build:
	cd $(ROOT) && $(GO) build -o $(DIST)/cmd ./cmd/main.go

image:
	$(DOCKER) build -f $(ROOT)/Dockerfile -t $(DOCKER_REPO)/line-bot-prac:$(TAG) $(ROOT)

dc.run: d.build
	$(DOCKER_COMPOSE) up

clean:
	rm -rf $(DIST)

