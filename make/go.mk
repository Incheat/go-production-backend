# make/go.mk
GO ?= go
GOFILES := ./...
GOLANGCI_LINT ?= golangci-lint

.PHONY: lint test unit-test pact-test fmt tidy ci tools

lint: ## Run golangci-lint
	@echo "==> Running golangci-lint..."
	$(GOLANGCI_LINT) run ./...

test: ## Run all go tests
	@echo "==> Running all go tests..."
	$(GO) test $(GOFILES) -v

unit-test: ## Run unit tests (-run Unit)
	@echo "==> Running unit tests..."
	$(GO) test $(GOFILES) -run Unit -short

pact-test: ## Run pact tests (-run Pact)
	@echo "==> Running pact tests..."
	$(GO) test $(GOFILES) -run Pact -short

fmt: ## Format code
	@echo "==> Running fmt..."
	$(GO) fmt $(GOFILES)

tidy: ## go mod tidy
	@echo "==> Running mod tidy..."
	$(GO) mod tidy

ci: lint test ## CI pipeline
