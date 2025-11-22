.PHONY: help docker-up docker-down lint test run all clean start

GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m

DOCKER_COMPOSE := docker compose
GO := go
GOLANGCI_LINT := golangci-lint
APP_NAME := Service-for-assigning-reviewers-for-Pull-Requests

docker-up:
	@echo "$(YELLOW)Starting containers...$(NC)"
	@$(DOCKER_COMPOSE) up -d --build
	@echo "$(YELLOW)Waiting for database to be ready...$(NC)"
	@timeout=30; \
	username=$$(grep -E '^POSTGRES_UNSAFE_USERNAME=' .env 2>/dev/null | cut -d '=' -f2 || echo "myuser"); \
	while [ $$timeout -gt 0 ]; do \
		if $(DOCKER_COMPOSE) exec -T database pg_isready -U $$username -d prdb > /dev/null 2>&1; then \
			echo "$(GREEN)✓ Successful container execution - PostgreSQL is ready!$(NC)"; \
			break; \
		fi; \
		sleep 1; \
		timeout=$$((timeout-1)); \
	done; \
	if [ $$timeout -eq 0 ]; then \
		echo "$(RED)✗ Database failed to start within 30 seconds$(NC)"; \
		exit 1; \
	fi

lint:
	@echo "$(YELLOW)Running linters...$(NC)"
	@$(GOLANGCI_LINT) run --timeout 5m --fix
	@echo "$(GREEN)✓ Successful linter execution$(NC)"

unit-test:
	@echo "$(YELLOW)Running unit tests...$(NC)"
	@$(GO) test ./internal/service/... ./util/... -coverprofile=coverage.out && $(GO) tool cover -func=coverage.out
	@echo "$(GREEN)✓ Successful tests execution - all tests passed$(NC)"

integration-test:
	@echo "$(YELLOW)Running integration tests...$(NC)"
	@$(GO) test -v ./tests/integration
	@echo "$(GREEN)✓ Successful integration tests execution - all tests passed$(NC)"

build:
	@$(GO) build -o app ./cmd/app

run:
	@echo "$(YELLOW)Starting application...$(NC)"
	@echo "$(GREEN)✓ Successful application startup$(NC)"
	./app

clean:
	@echo "$(YELLOW)Stopping and removing containers...$(NC)"
	@$(DOCKER_COMPOSE) down -v

start: docker-up lint unit-test integration-test build run clean