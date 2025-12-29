# make/migrate.mk
# ----------------------------------------
# Migrate (per service)
# ----------------------------------------

MIGRATE ?= migrate
SERVICES_DIR ?= services

# Per-service paths (depends on SERVICE)
SERVICE_DB_DIR := $(SERVICES_DIR)/$(SERVICE)/db
MIGRATIONS_DIR := $(SERVICE_DB_DIR)/migrations

.PHONY: migrate-help migrate-up migrate-down1 migrate-version migrate-force migrate-create migrate-up-all

migrate-help: ## Show migrate usage
	@echo "Usage:"
	@echo "  make migrate-up SERVICE=foo MYSQL_DSN='mysql://user:pass@tcp(127.0.0.1:3306)/db?parseTime=true&multiStatements=true'"
	@echo "  make migrate-down1 SERVICE=foo MYSQL_DSN='...'"
	@echo "  make migrate-version SERVICE=foo MYSQL_DSN='...'"
	@echo "  make migrate-force SERVICE=foo MYSQL_DSN='...' VERSION=12"
	@echo "  make migrate-create SERVICE=foo NAME=add_users"
	@echo "  make migrate-up-all MYSQL_DSN='...'"
	@echo ""
	@echo "Notes:"
	@echo "  - Looks for migrations in: services/<SERVICE>/db/migrations"
	@echo "  - Skips a service if the migrations directory doesn't exist"

define require_service_dsn_and_dir
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

migrate-up: ## Migrate up for one service (SERVICE=... MYSQL_DSN=...)
	@$(call require_service_dsn_and_dir)
	@echo "=== Migrating UP for service: $(SERVICE) ==="
	@$(MIGRATE) -database "$(MYSQL_DSN)" -path "$(MIGRATIONS_DIR)" up

migrate-down1: ## Migrate down 1 for one service
	@$(call require_service_dsn_and_dir)
	@echo "=== Migrating DOWN 1 for service: $(SERVICE) ==="
	@$(MIGRATE) -database "$(MYSQL_DSN)" -path "$(MIGRATIONS_DIR)" down 1

migrate-version: ## Show migration version for one service
	@$(call require_service_dsn_and_dir)
	@echo "=== Migration VERSION for service: $(SERVICE) ==="
	@$(MIGRATE) -database "$(MYSQL_DSN)" -path "$(MIGRATIONS_DIR)" version

migrate-force: ## Force migration version (SERVICE=... MYSQL_DSN=... VERSION=...)
	@$(call require_service_dsn_and_dir)
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required (e.g. VERSION=12)"; exit 1; fi
	@echo "=== Forcing version $(VERSION) for service: $(SERVICE) ==="
	@$(MIGRATE) -database "$(MYSQL_DSN)" -path "$(MIGRATIONS_DIR)" force $(VERSION)

migrate-create: ## Create migration files (SERVICE=... NAME=...)
	@if [ -z "$(SERVICE)" ]; then echo "SERVICE is required"; exit 1; fi
	@if [ -z "$(NAME)" ]; then echo "NAME is required (e.g. NAME=add_users)"; exit 1; fi
	@mkdir -p "$(MIGRATIONS_DIR)"
	@echo "=== Creating migration for service: $(SERVICE) name: $(NAME) ==="
	@$(MIGRATE) create -ext sql -dir "$(MIGRATIONS_DIR)" -seq "$(NAME)"

migrate-up-all: ## Migrate up for all services (MYSQL_DSN=...)
	@if [ -z "$(MYSQL_DSN)" ]; then echo "MYSQL_DSN is required"; exit 1; fi
	@set -e; \
	for d in $(SERVICES_DIR)/*/db/migrations; do \
		if [ -d "$$d" ]; then \
			svc=$$(echo "$$d" | sed -E 's#^$(SERVICES_DIR)/([^/]+)/db/migrations#\1#'); \
			echo "=== Migrating UP for service: $$svc ==="; \
			$(MIGRATE) -database "$(MYSQL_DSN)" -path "$$d" up; \
		fi; \
	done
