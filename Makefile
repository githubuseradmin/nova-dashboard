# nova — developer tasks. Run `make` or `make help` for the list.
.DEFAULT_GOAL := help
.PHONY: help dev build test tidy fmt vet db-up db-down db-logs web-install web-dev web-build

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

## --- Backend (Go) ---
dev: ## Run the API server (needs Postgres up)
	go run ./cmd/server

build: ## Build the server binary into ./bin
	go build -o bin/server ./cmd/server

test: ## Run all Go tests
	go test ./...

tidy: ## Sync go.mod / go.sum
	go mod tidy

fmt: ## Format Go code
	go fmt ./...

vet: ## Static checks
	go vet ./...

## --- Database (Docker) ---
db-up: ## Start Postgres
	docker compose up -d

db-down: ## Stop Postgres
	docker compose down

db-logs: ## Tail Postgres logs
	docker compose logs -f postgres

## --- Frontend (Svelte) ---
web-install: ## Install frontend deps
	cd web/app && npm install

web-dev: ## Run the Svelte dev server
	cd web/app && npm run dev

web-build: ## Build the frontend for production
	cd web/app && npm run build
