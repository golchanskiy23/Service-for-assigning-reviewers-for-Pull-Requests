.PHONY: help docker-up docker-down lint test run all clean start

GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m

DOCKER_COMPOSE := docker compose
GO := go
GOLANGCI_LINT := golangci-lint
APP_NAME := Service-for-assigning-reviewers-for-Pull-Requests


lint:
	@echo "$(YELLOW)Running linters...$(NC)"
	@$(GOLANGCI_LINT) run --timeout 5m
	@echo "$(GREEN)✓ Successful linter execution$(NC)"

test:
	@echo "$(YELLOW)Running tests in all directories...$(NC)"
	@$(GO) test ./internal/service/... ./internal/handlers/... ./util/... -coverprofile=coverage.out && $(GO) tool cover -func=coverage.out
	@echo "$(GREEN)✓ Successful tests execution - all tests passed$(NC)"

run:
	@echo "$(YELLOW)Starting application...$(NC)"
	@echo "$(GREEN)✓ Successful application startup$(NC)"
	@$(GO) run ./cmd/app/main.go

clean:
	@echo "$(YELLOW)Stopping and removing containers...$(NC)"
	@$(DOCKER_COMPOSE) stop -v

#start: lint test run clean
start: run