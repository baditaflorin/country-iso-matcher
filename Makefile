.PHONY: build test lint run clean docker-build docker-run help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=country-iso-matcher
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker parameters
DOCKER_IMAGE=country-iso-service-go
DOCKER_TAG=latest

all: test build

build: ## Build the binary
	$(GOBUILD) -o $(BINARY_NAME) -v ./src/cmd/server

test: ## Run tests
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests and show coverage
	$(GOCMD) tool cover -html=coverage.out

lint: ## Run linter
	golangci-lint run

run: ## Run the application
	$(GOBUILD) -o $(BINARY_NAME) -v ./src/cmd/server && ./$(BINARY_NAME)

clean: ## Clean build files
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container
	docker run -p 3030:3030 $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-compose-up: ## Run with docker-compose
	docker-compose up --build

docker-compose-down: ## Stop docker-compose
	docker-compose down

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $1, $2}'