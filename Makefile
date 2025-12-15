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
GLOBAL_OAPI_DIR := api/$(SERVICE)/oapi
LOCAL_API_DIR := services/$(SERVICE)/internal/api/oapi

PUBLIC_SERVER_CONFIG := $(GLOBAL_OAPI_DIR)/public.server.yaml
INTERNAL_SERVER_CONFIG := $(GLOBAL_OAPI_DIR)/internal.server.yaml

PUBLIC_CLIENT_CONFIG := $(GLOBAL_OAPI_DIR)/public.client.yaml
INTERNAL_CLIENT_CONFIG := $(GLOBAL_OAPI_DIR)/internal.client.yaml

PUBLIC_TYPES_CONFIG := $(GLOBAL_OAPI_DIR)/public.type.yaml
INTERNAL_TYPES_CONFIG := $(GLOBAL_OAPI_DIR)/internal.type.yaml

PUBLIC_OPENAPI := $(GLOBAL_OAPI_DIR)/public.yaml
INTERNAL_OPENAPI := $(GLOBAL_OAPI_DIR)/internal.yaml

# ----------------------------------------
# Generate code for one service
# ----------------------------------------
.PHONY: generate

generate:
	@echo "=== Generating API code for service: $(SERVICE) ==="

	@if [ -f "$(PUBLIC_OPENAPI)" ]; then \
		oapi-codegen --config $(PUBLIC_SERVER_CONFIG) $(PUBLIC_OPENAPI); \
		oapi-codegen --config $(PUBLIC_CLIENT_CONFIG) $(PUBLIC_OPENAPI); \
		oapi-codegen --config $(PUBLIC_TYPES_CONFIG) $(PUBLIC_OPENAPI); \
		echo "✓ Public API generation completed for $(SERVICE)"; \
	else \
		echo "⚠ Skipping PUBLIC API: $(PUBLIC_OPENAPI) not found."; \
	fi

	@if [ -f "$(INTERNAL_OPENAPI)" ]; then \
		oapi-codegen --config $(INTERNAL_SERVER_CONFIG) $(INTERNAL_OPENAPI); \
		oapi-codegen --config $(INTERNAL_CLIENT_CONFIG) $(INTERNAL_OPENAPI); \
		oapi-codegen --config $(INTERNAL_TYPES_CONFIG) $(INTERNAL_OPENAPI); \
		echo "✓ Internal API generation completed for $(SERVICE)"; \
	else \
		echo "⚠ Skipping INTERNAL API: $(INTERNAL_OPENAPI) not found."; \
	fi
	@echo "\n"

# ----------------------------------------
# Generate code for ALL services
# ----------------------------------------
SERVICES := $(shell ls api)

.PHONY: generate-all
generate-all:
	@echo "\nGenerating API code for ALL services: $(SERVICES)\n"
	@for svc in $(SERVICES); do \
		$(MAKE) generate SERVICE=$$svc; \
	done
	@echo "All services updated.\n"

# ----------------------------------------
# Generate SQLC code for one service
# ----------------------------------------
DB ?= mysql
.PHONY: sqlc

sqlc:
	@if [ -d "services/$(SERVICE)/db/$(DB)" ]; then \
		echo "=== Generating SQLC code for service: $(SERVICE) ==="; \
		cd services/$(SERVICE)/db/$(DB) && sqlc generate; \
		echo "=== Done for service: $(SERVICE) ==="; \
	else \
		echo "=== Skip SQLC: services/$(SERVICE)/db/$(DB) not found ==="; \
	fi

# ----------------------------------------
# Generate SQLC code for ALL services
# ----------------------------------------
DBS := $(shell ls services/*/db)
.PHONY: sqlc-all

sqlc-all:
	@echo "=== Generating SQLC code for ALL services: $(SERVICES) ==="
	@for svc in $(SERVICES); do \
		echo "Generating SQLC code for service: $$svc"; \
		for db in $(DBS); do \
			$(MAKE) sqlc SERVICE=$$svc DB=$$db; \
		done; \
	done
	@echo "All services updated."


# ----------------------------------------
# Migrate (per service)
# ----------------------------------------

MIGRATE ?= migrate
SERVICE ?=
SERVICES_DIR ?= services

# Per-service paths
SERVICE_DB_DIR := $(SERVICES_DIR)/$(SERVICE)/db
MIGRATIONS_DIR := $(SERVICE_DB_DIR)/migrations

# Usage helper
.PHONY: migrate-help
migrate-help:
	@echo "Usage:"
	@echo "  make migrate-up SERVICE=foo MYSQL_DSN='mysql://user:pass@tcp(127.0.0.1:3306)/ordersdb?parseTime=true&multiStatements=true'"
	@echo "  make migrate-down1 SERVICE=foo MYSQL_DSN='...'"
	@echo "  make migrate-version SERVICE=foo MYSQL_DSN='...'"
	@echo "  make migrate-create SERVICE=foo NAME=add_users"
	@echo "  make migrate-up-all MYSQL_DSN='...'"
	@echo ""
	@echo "Notes:"
	@echo "  - Looks for migrations in: services/<SERVICE>/db/migrations"
	@echo "  - Skips if the directory doesn't exist"

# Internal check: requires SERVICE + MYSQL_DSN, and migrations dir exists
define require_service_and_dsn_and_dir
	@if [ -z "$(SERVICE)" ]; then \
		echo "SERVICE is required (e.g. SERVICE=user)"; exit 1; \
	fi
	@if [ -z "$(MYSQL_DSN)" ]; then \
		echo "MYSQL_DSN is required"; exit 1; \
	fi
	@if [ ! -d "$(MIGRATIONS_DIR)" ]; then \
		echo "=== Skip migrate: $(MIGRATIONS_DIR) not found ==="; exit 0; \
	fi
endef

.PHONY: migrate-up migrate-down1 migrate-version migrate-force migrate-create

migrate-up:
	@$(call require_service_and_dsn_and_dir)
	@echo "=== Migrating UP for service: $(SERVICE) ==="
	@$(MIGRATE) -database "$(MYSQL_DSN)" -path "$(MIGRATIONS_DIR)" up

migrate-down1:
	@$(call require_service_and_dsn_and_dir)
	@echo "=== Migrating DOWN 1 for service: $(SERVICE) ==="
	@$(MIGRATE) -database "$(MYSQL_DSN)" -path "$(MIGRATIONS_DIR)" down 1

migrate-version:
	@$(call require_service_and_dsn_and_dir)
	@echo "=== Migration VERSION for service: $(SERVICE) ==="
	@$(MIGRATE) -database "$(MYSQL_DSN)" -path "$(MIGRATIONS_DIR)" version

# Force schema_migrations to a version (useful when fixing a broken state)
# Usage: make migrate-force SERVICE=foo MYSQL_DSN='...' VERSION=12
migrate-force:
	@$(call require_service_and_dsn_and_dir)
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required (e.g. VERSION=12)"; exit 1; fi
	@echo "=== Forcing version $(VERSION) for service: $(SERVICE) ==="
	@$(MIGRATE) -database "$(MYSQL_DSN)" -path "$(MIGRATIONS_DIR)" force $(VERSION)

# Create new migration pair files in the service migrations dir
# Usage: make migrate-create SERVICE=foo NAME=add_users
migrate-create:
	@if [ -z "$(SERVICE)" ]; then echo "SERVICE is required"; exit 1; fi
	@if [ -z "$(NAME)" ]; then echo "NAME is required (e.g. NAME=add_users)"; exit 1; fi
	@mkdir -p "$(MIGRATIONS_DIR)"
	@echo "=== Creating migration for service: $(SERVICE) name: $(NAME) ==="
	@$(MIGRATE) create -ext sql -dir "$(MIGRATIONS_DIR)" -seq "$(NAME)"

# ---- All services -----------------------------------------------------------

.PHONY: migrate-up-all
migrate-up-all:
	@if [ -z "$(MYSQL_DSN)" ]; then echo "MYSQL_DSN is required"; exit 1; fi
	@set -e; \
	for d in $(SERVICES_DIR)/*/db/migrations; do \
		if [ -d "$$d" ]; then \
			svc=$$(echo "$$d" | sed -E 's#^$(SERVICES_DIR)/([^/]+)/db/migrations#\1#'); \
			echo "=== Migrating UP for service: $$svc ==="; \
			$(MIGRATE) -database "$(MYSQL_DSN)" -path "$$d" up; \
		fi; \
	done

# ----------------------------------------
# Build image and load into kind cluster
# ----------------------------------------
ENV ?= dev

build-all:
	@echo "Building image for all services"
	@for svc in $(SERVICES); do \
		$(MAKE) build-local SERVICE=$$svc ENV=$(ENV); \
	done

build-local:
	@echo "Building image for $(SERVICE)"
	docker build -t $(SERVICE):$(ENV) ./services/$(SERVICE)
	kind load docker-image $(SERVICE):$(ENV)

# ----------------------------------------
# Run microservices locally
# ----------------------------------------

# make run-local SERVICE=auth
run-all-local:
	@echo "Running all services in dev mode"
	@set -a; . ./.env.local; set +a; \
	for svc in $(SERVICES); do \
		echo "==> $$svc"; \
		go run ./services/$$svc/cmd/main.go; \
	done

run-local:
	@echo "$(SERVICE) service running in dev mode"
	set -a; . .env.local; set +a; \
	go run ./services/$(SERVICE)/cmd/main.go