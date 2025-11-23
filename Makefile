.PHONY: help docker-up docker-down lint test run all clean start

GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m

DOCKER_COMPOSE := docker-compose
GO := go
GOLANGCI_LINT := golangci-lint
APP_NAME := Service-for-assigning-reviewers-for-Pull-Requests

# впоследствии код будет запускаться только через docker-compose up
docker-up:
	@echo "$(YELLOW)Starting PostgreSQL container...$(NC)"
	@$(DOCKER_COMPOSE) up -d
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
	@echo "$(YELLOW)Running linter...$(NC)"
	@$(GOLANGCI_LINT) run
	@echo "$(GREEN)✓ Successful linter execution$(NC)"

test:
	@echo "$(YELLOW)Running tests in all directories...$(NC)"
	@$(GO) test ./internal/service/... ./internal/handlers/... ./util/... -coverprofile=coverage.out && $(GO) tool cover -func=coverage.out
	@echo "$(GREEN)✓ Successful tests execution - all tests passed$(NC)"

run:
	@echo "$(YELLOW)Starting application...$(NC)"
	@echo "$(GREEN)✓ Successful application startup$(NC)"
	@$(GO) run ./cmd/app/main.go

start: docker-up test run