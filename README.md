# **Country ISO Matcher**

A high-performance Go web service that converts country names, including common aliases and variations, into their official name and **ISO 3166-1 alpha-2** code. It's designed for simplicity, performance, and easy integration, featuring built-in observability with structured logging and Prometheus metrics.

## **✨ Features**

* **Intelligent Name Matching**: Handles variations in casing, accents (diacritics), and surrounding whitespace.
* **Multi-lingual Alias Support**: Recognizes common country names and aliases in multiple languages (e.g., "Germany", "Deutschland", "Allemagne").
* **High Performance**: Utilizes an in-memory data store for fast, sub-millisecond lookups.
* **RESTful API**: Provides simple and predictable API endpoints for conversion, health checks, and stats.
* **Observability Ready**:
    * Exports detailed Prometheus metrics on the /metrics endpoint.
    * Offers a /stats endpoint for a high-level overview of application performance.
    * Structured JSON logging for efficient log parsing and analysis.
* **Containerized**: Comes with a multi-stage Dockerfile, docker-compose.yml, and Kubernetes manifests for easy deployment.
* **Robust & Resilient**: Implements graceful shutdown and a panic recovery middleware.

## **🚀 Getting Started**

### **Prerequisites**

* [Go](https://golang.org/dl/) version 1.23 or later
* [Docker](https://www.docker.com/get-started) & [Docker Compose](https://docs.docker.com/compose/install/)
* [Make](https://www.gnu.org/software/make/)

### **Local Installation**

1. **Clone the repository:**  
   git clone \[https://github.com/baditaflorin/country-iso-matcher.git\](https://github.com/baditaflorin/country-iso-matcher.git)  
   cd country-iso-matcher

2. Set up environment variables:  
   Copy the example environment file. No changes are needed to run locally.  
   cp .env.example .env

3. **Install dependencies:**  
   make deps

4. **Run the application:**  
   make run

The server will start on http://localhost:3030.

## **🛠️ Usage**

### **API Endpoints**

The service exposes the following endpoints:

#### **GET /api/convert**

Converts a country name to its ISO code and official name.

* **Query Parameter**: country (string, required) \- The name of the country to look up.
* **Example Request:**  
  curl "http://localhost:3030/api/convert?country=Deutschland"

* **Success Response (200 OK):**  
  {  
  "query": "Deutschland",  
  "officialName": "Germany",  
  "isoCode": "DE"  
  }

* **Not Found Response (404 Not Found):**  
  {  
  "error": "Country not found: Atlantis",  
  "query": "Atlantis"  
  }

* **Validation Error Response (400 Bad Request):**  
  {  
  "error": "Country query parameter is required",  
  "query": ""  
  }

#### **GET /health**

Provides a simple health check for monitoring systems.

* **Example Request:**  
  curl "http://localhost:3030/health"

* **Success Response (200 OK):**  
  {  
  "status": "healthy",  
  "service": "country-iso-matcher"  
  }

#### **GET /stats**

Returns aggregated statistics derived from the Prometheus metrics, including request counts and the top 10 most popular country lookups.

* **Example Request:**  
  curl "http://localhost:3030/stats"

* **Example Response (200 OK):**  
  {  
  "total\_requests": 5,  
  "success\_count": 3,  
  "not\_found\_count": 1,  
  "error\_count": 0,  
  "validation\_error\_count": 1,  
  "success\_rate": 0.6,  
  "failure\_rate": 0.4,  
  "popular\_countries": \[  
  {  
  "code": "RO",  
  "name": "Romania",  
  "count": 2  
  },  
  {  
  "code": "DE",  
  "name": "Germany",  
  "count": 1  
  }  
  \]  
  }

#### **GET /metrics**

Exposes detailed application and request metrics in Prometheus format.

### **Makefile Commands**

A Makefile is included for common development tasks.

* make build: Builds the Go binary.
* make run: Builds and runs the application locally.
* make test: Runs all unit tests.
* make test-coverage: Runs tests and opens an HTML coverage report.
* make lint: Lints the codebase using golangci-lint.
* make clean: Removes build artifacts.
* make deps: Tidies and downloads Go modules.
* make docker-build: Builds the Docker image.
* make docker-run: Runs the application inside a Docker container.
* make docker-compose-up: Starts the application using Docker Compose.
* make docker-compose-down: Stops the application running via Docker Compose.
* make help: Shows a list of all available commands.

## **🐳 Docker & Kubernetes**

### **Docker**

Build and run the container directly:

\# Build the image  
make docker-build

\# Run the container  
make docker-run

### **Docker Compose**

The simplest way to run the application with its production configuration is using Docker Compose.

\# Build and start the service in the background  
docker-compose up \--build \-d

The service will be available at http://localhost:3030.

### **Kubernetes**

Basic Kubernetes manifests for deploying the service are available in the /k8s directory. These can be used as a starting point for a production deployment.

\# Apply the deployment and service  
kubectl apply \-f k8s/deployment.yaml

\# Check the status  
kubectl get pods \-l app=country-iso-matcher

## **⚙️ Configuration**

The application is configured via environment variables.

| Variable | Description | Default |
| :---- | :---- | :---- |
| PORT | Port for the HTTP server to use. | 3030 |
| ENV | Application environment. | development |
| READ\_TIMEOUT | HTTP server read timeout (sec). | 10 |
| WRITE\_TIMEOUT | HTTP server write timeout (sec). | 10 |
| LOG\_LEVEL | Logging level (info, debug). | info |
| LOG\_FORMAT | Logging format (text, json). | json |

## **🧪 Testing & Benchmarks**

* **Run Unit Tests**:  
  make test

* **Run Benchmarks**:  
  go test \-bench=. ./benchmarks/...

## **📂 Project Structure**

.  
├── benchmarks/         \# Go benchmark tests  
├── k8s/                \# Kubernetes deployment manifests  
├── src/  
│   ├── cmd/server/     \# Application entrypoint (main.go)  
│   ├── internal/       \# Internal application logic (not for export)  
│   │   ├── config/     \# Environment configuration  
│   │   ├── domain/     \# Core data structures and errors  
│   │   ├── factory/    \# Dependency injection and object creation  
│   │   ├── handler/    \# HTTP handlers and middleware  
│   │   ├── metrics/    \# Prometheus metric definitions  
│   │   ├── repository/ \# Data access layer (in-memory)  
│   │   ├── server/     \# HTTP server setup  
│   │   └── service/    \# Business logic  
│   └── pkg/            \# Reusable packages (e.g., normalizer)  
├── .golangci.yml       \# Linter configuration  
├── Dockerfile          \# Multi-stage Docker build file  
├── docker-compose.yml  \# Docker compose configuration  
├── go.mod              \# Go module definition  
└── Makefile            \# Make commands for development

## **📜 License**

This project is licensed under the MIT License. See the [LICENSE](https://www.google.com/search?q=LICENSE) file for details.
