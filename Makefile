.PHONY: build test lint run clean docker-build docker-run help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
BINARY_NAME=server
BINARY_DIR=bin
BINARY_PATH=$(BINARY_DIR)/$(BINARY_NAME)

# Docker parameters
DOCKER_IMAGE=country-iso-matcher
DOCKER_TAG=latest

all: test build

build: ## Build the binary
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_PATH) -v ./src/cmd/server
	@echo "Build complete: $(BINARY_PATH)"

test: ## Run tests
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests and show coverage
	$(GOCMD) tool cover -html=coverage.out

benchmark: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem ./...

lint: ## Run linter
	golangci-lint run --timeout=5m

fmt: ## Format code
	$(GOFMT) ./...

run: build ## Build and run the application
	./$(BINARY_PATH)

dev: ## Run with config file for development
	./$(BINARY_PATH) -config configs/config.example.yaml

clean: ## Clean build files
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)
	rm -f coverage.out

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

# Docker commands
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container
	docker run -p 3030:3030 \
		-v $(PWD)/data:/app/data \
		-v $(PWD)/web:/app/web \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

docker-compose-up: ## Run with docker-compose
	docker-compose up --build

docker-compose-down: ## Stop docker-compose
	docker-compose down

docker-compose-logs: ## View docker-compose logs
	docker-compose logs -f app

# Kubernetes commands
k8s-apply: ## Apply Kubernetes manifests
	kubectl apply -f k8s/

k8s-delete: ## Delete Kubernetes resources
	kubectl delete -f k8s/

k8s-logs: ## View Kubernetes logs
	kubectl logs -f -l app=country-iso-matcher

# Data management
generate-data: ## Generate CSV data from memory loader
	@echo "CSV data files already exist in data/ directory"
	@echo "Edit data/countries.csv and data/aliases.csv to customize"

# Development helpers
watch: ## Watch for file changes and rebuild
	@which air > /dev/null || $(GOGET) -u github.com/cosmtrek/air
	air

install-tools: ## Install development tools
	$(GOGET) -u github.com/cosmtrek/air
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint

# Quick start
quickstart: build ## Quick start with default settings
	@echo "Starting Country ISO Matcher..."
	@echo "API: http://localhost:3030/api/convert?country=Germany"
	@echo "GUI: http://localhost:3030/admin"
	@echo "Health: http://localhost:3030/health"
	@echo "Metrics: http://localhost:3030/metrics"
	./$(BINARY_PATH)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $1, $2}'
