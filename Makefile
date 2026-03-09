.PHONY: help build test lint clean docker-up docker-down proto fmt vet

SERVICES := auth-service config-service user-service franchise-service customer-service \
	address-service rating-service carrier-service shipment-service document-service \
	notification-service ecommerce-service invoice-service reporting-service \
	payment-service freight-service live-rating-service

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

## --- Go Backend ---

build: ## Build all services
	@for svc in $(SERVICES); do \
		echo "Building $$svc..."; \
		(cd services/$$svc && go build -o ../../bin/$$svc ./cmd/server/) || exit 1; \
	done
	@echo "All services built successfully"

test: ## Run all tests
	go test ./pkg/... ./services/...

test-pkg: ## Run shared package tests only
	go test -v ./pkg/...

test-svc: ## Run tests for a specific service (usage: make test-svc SVC=auth-service)
	go test -v ./services/$(SVC)/...

lint: ## Run golangci-lint
	golangci-lint run ./pkg/... ./services/...

fmt: ## Format Go code
	gofmt -w pkg/ services/

vet: ## Run go vet
	go vet ./pkg/... ./services/...

clean: ## Remove build artifacts
	rm -rf bin/
	@for svc in $(SERVICES); do \
		rm -f services/$$svc/server; \
	done

## --- Infrastructure ---

docker-up: ## Start all infrastructure services
	docker-compose -f deployments/docker-compose.yml up -d

docker-down: ## Stop all infrastructure services
	docker-compose -f deployments/docker-compose.yml down

docker-logs: ## Tail infrastructure logs
	docker-compose -f deployments/docker-compose.yml logs -f

docker-reset: ## Reset all infrastructure (WARNING: deletes data)
	docker-compose -f deployments/docker-compose.yml down -v

## --- Protobuf ---

proto: ## Generate protobuf Go code
	./scripts/generate-proto.sh

## --- Frontend ---

frontend-install: ## Install frontend dependencies
	cd frontend && pnpm install

frontend-dev: ## Start frontend dev servers
	cd frontend && pnpm dev

frontend-build: ## Build frontend
	cd frontend && pnpm build

frontend-test: ## Run frontend tests
	cd frontend && pnpm test

## --- Development ---

dev: docker-up ## Start full dev environment
	@echo "Infrastructure started. Run 'make run-svc SVC=auth-service' to start a service."

run-svc: ## Run a specific service (usage: make run-svc SVC=auth-service)
	cd services/$(SVC) && go run ./cmd/server/

run-all: ## Run all services (background)
	@for svc in $(SERVICES); do \
		echo "Starting $$svc..."; \
		(cd services/$$svc && go run ./cmd/server/ &); \
	done
	@echo "All services starting..."
