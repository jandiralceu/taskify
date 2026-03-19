.PHONY: docker-up docker-up-all docker-stop docker-down help lint test

COMPOSE_FILE = deployment/compose.yaml

docker-up: ## Start only background docker containers (DB, Redis, Jaeger) in detached mode
	docker compose -f $(COMPOSE_FILE) up -d 

docker-up-all: ## Start ALL docker containers including the API app in detached mode
	docker compose -f $(COMPOSE_FILE) --profile all up -d --build

docker-stop: ## Stop docker containers without removing them
	docker compose -f $(COMPOSE_FILE) --profile all stop

docker-down: ## Stop and remove docker containers
	docker compose -f $(COMPOSE_FILE) --profile all down

lint: ## Run linter across both backend and frontend
	@echo "Linting backend..."
	@cd backend && make lint
	@echo "Linting frontend..."
	@cd frontend && npm run lint

test: ## Run all tests across both backend and frontend
	@echo "Testing backend..."
	@cd backend && make test
	@echo "Testing frontend..."
	@cd frontend && npm run test:unit -- --run

generate-keys: ## Generate RSA key pair for JWT (proxied to backend)
	@cd backend && make generate-keys

db-dump: ## Create a data-only dump of the database to deployment/seed.sql
	@echo "Dumping database data..."
	@docker exec -e PGPASSWORD=nk6eloPZ1FQk8Ku1 taskify-postgres pg_dump -U jandiralceu -d taskify --data-only --column-inserts --no-owner --no-privileges > deployment/seed.sql

db-restore: ## Restore the database data from deployment/seed.sql
	@echo "Restoring database data..."
	@docker exec -i -e PGPASSWORD=nk6eloPZ1FQk8Ku1 taskify-postgres psql -U jandiralceu -d taskify < deployment/seed.sql

help: ## Display all available commands
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
