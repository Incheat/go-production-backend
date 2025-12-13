# ----------------------------------------
# Makefile for the go-playground project
# ----------------------------------------
# --- Go settings ---
GO              ?= go
GOFILES         := ./...
GOMOD           := go.mod

# --- Tools ---
GOLANGCI_LINT   ?= golangci-lint

# --- Commands ---
.PHONY: all lint test fmt tidy ci tools

all: lint test

## Run linters (golangci-lint)
lint:
	@echo "==> Running golangci-lint..."
	$(GOLANGCI_LINT) run ./...

## Run Go tests
test:
	@echo "==> Running go test..."
	$(GO) test $(GOFILES) -v

## Run Unit tests
unit-test:
	@echo "==> Running go test..."
	$(GO) test $(GOFILES) -run Unit -short

## Run Pact tests
pact-test:
	@echo "==> Running pact tests..."
	$(GO) test $(GOFILES) -run Pact -short

## Format code (goimports + gofmt)
fmt:
	@echo "==> Running gofmt..."
	$(GO) fmt $(GOFILES)
	@echo "==> Optionally run goimports (uncomment below if installed)..."
	# goimports -w .

## Keep go.mod / go.sum tidy
tidy:
	@echo "==> Running go mod tidy..."
	$(GO) mod tidy

## CI: run everything you'd want in CI
ci: lint test

## Install tools (optional helper)
tools:
	@echo "==> Installing golangci-lint if missing (Mac with Homebrew)..."
	@if ! command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		echo "golangci-lint not found, installing with brew..."; \
		brew install golangci-lint || echo "Install golangci-lint manually"; \
	else \
		echo "golangci-lint already installed"; \
	fi

# ----------------------------------------
# Generate code for ALL services
# ----------------------------------------

# make generate SERVICE=helloworld

# The service to generate code for:
SERVICE ?= helloworld

# Paths based on SERVICE
PUBLIC_SERVER_CONFIG := api/$(SERVICE)/oapi/public.server.yaml
PUBLIC_CLIENT_CONFIG := api/$(SERVICE)/oapi/public.client.yaml
INTERNAL_SERVER_CONFIG := api/$(SERVICE)/oapi/internal.server.yaml
INTERNAL_CLIENT_CONFIG := api/$(SERVICE)/oapi/internal.client.yaml
PUBLIC_OPENAPI := api/$(SERVICE)/oapi/public.yaml
INTERNAL_OPENAPI := api/$(SERVICE)/oapi/internal.yaml
API_DIR := services/$(SERVICE)/internal/api
GEN_DIR := $(API_DIR)/gen
GEN_PUBLIC_DIR := $(GEN_DIR)/oapi/public
GEN_INTERNAL_DIR := $(GEN_DIR)/oapi/internal
GEN_PUBLIC_SERVER_FILE := $(GEN_PUBLIC_DIR)/server.gen.go
GEN_PUBLIC_CLIENT_FILE := $(GEN_PUBLIC_DIR)/client.gen.go
GEN_INTERNAL_SERVER_FILE := $(GEN_INTERNAL_DIR)/server.gen.go
GEN_INTERNAL_CLIENT_FILE := $(GEN_INTERNAL_DIR)/client.gen.go

# ----------------------------------------
# Generate code for one service
# ----------------------------------------
.PHONY: generate

generate:
	@echo "=== Generating API code for service: $(SERVICE) ==="

	# -------------------------------
	# PUBLIC API (skip if missing)
	# -------------------------------
	@if [ -f "$(PUBLIC_OPENAPI)" ]; then \
		echo "-- Public OpenAPI found: $(PUBLIC_OPENAPI)"; \
		mkdir -p $(GEN_PUBLIC_DIR); \
		oapi-codegen --config $(PUBLIC_SERVER_CONFIG) $(PUBLIC_OPENAPI); \
		oapi-codegen --config $(PUBLIC_CLIENT_CONFIG) $(PUBLIC_OPENAPI); \
		echo "✓ Public API generation completed for $(SERVICE)"; \
	else \
		echo "⚠ Skipping PUBLIC API: $(PUBLIC_OPENAPI) not found."; \
	fi

	# -------------------------------
	# INTERNAL API (skip if missing)
	# -------------------------------
	@if [ -f "$(INTERNAL_OPENAPI)" ]; then \
		echo "-- Internal OpenAPI found: $(INTERNAL_OPENAPI)"; \
		mkdir -p $(GEN_INTERNAL_DIR); \
		oapi-codegen --config $(INTERNAL_SERVER_CONFIG) $(INTERNAL_OPENAPI); \
		oapi-codegen --config $(INTERNAL_CLIENT_CONFIG) $(INTERNAL_OPENAPI); \
		echo "✓ Internal API generation completed for $(SERVICE)"; \
	else \
		echo "⚠ Skipping INTERNAL API: $(INTERNAL_OPENAPI) not found."; \
	fi

	@echo "=== Done for service: $(SERVICE) ==="


# ----------------------------------------
# Generate code for ALL services
# ----------------------------------------
SERVICES := $(shell ls api)

.PHONY: generate-all
generate-all:
	@echo "Generating API code for ALL services: $(SERVICES)"
	@for svc in $(SERVICES); do \
		echo "Generating API code for service: $$svc"; \
		$(MAKE) generate SERVICE=$$svc; \
	done
	@echo "All services updated."

# ----------------------------------------
# Generate SQLC code for one service
# ----------------------------------------
.PHONY: sqlc

sqlc:
	@if [ -d "services/$(SERVICE)/db" ]; then \
		echo "=== Generating SQLC code for service: $(SERVICE) ==="; \
		cd services/$(SERVICE)/db && sqlc generate; \
		echo "=== Done for service: $(SERVICE) ==="; \
	else \
		echo "=== Skip SQLC: services/$(SERVICE)/db not found ==="; \
	fi

# ----------------------------------------
# Generate SQLC code for ALL services
# ----------------------------------------
.PHONY: sqlc-all

sqlc-all:
	@echo "=== Generating SQLC code for ALL services: $(SERVICES) ==="
	@for svc in $(SERVICES); do \
		echo "Generating SQLC code for service: $$svc"; \
		$(MAKE) sqlc SERVICE=$$svc; \
	done
	@echo "All services updated."

# ----------------------------------------
# Run the auth service with different environments
# ----------------------------------------
# make run-dev SERVICE=auth
run-dev:
	@echo "$(SERVICE) service running in dev mode"
	cd services/${SERVICE} && APP_ENV=dev go run ./cmd/main.go

run-stage:
	@echo "$(SERVICE) service running in stage mode"
	cd services/${SERVICE} && APP_ENV=stage go run ./cmd/main.go

run-prod:
	@echo "$(SERVICE) service running in prod mode"
	cd services/${SERVICE} && APP_ENV=prod go run ./cmd/main.go