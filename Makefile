.PHONY: docker-up docker-up-all docker-stop docker-down help

COMPOSE_FILE = deployment/compose.yaml

docker-up: ## Start only background docker containers (DB, Redis, Jaeger) in detached mode
	docker compose -f $(COMPOSE_FILE) up -d 

docker-up-all: ## Start ALL docker containers including the API app in detached mode
	docker compose -f $(COMPOSE_FILE) --profile all up -d --build

docker-stop: ## Stop docker containers without removing them
	docker compose -f $(COMPOSE_FILE) --profile all stop

docker-down: ## Stop and remove docker containers
	docker compose -f $(COMPOSE_FILE) --profile all down

help: ## Display all available commands
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
