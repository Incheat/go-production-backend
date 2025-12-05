# make generate SERVICE=helloworld

# The service to generate code for:
SERVICE ?= helloworld

# Paths based on SERVICE
SERVER_CONFIG := api/$(SERVICE)/oapi-codegen.server.yaml
CLIENT_CONFIG := api/$(SERVICE)/oapi-codegen.client.yaml
OPENAPI := api/$(SERVICE)/openapi.yaml
API_DIR := services/$(SERVICE)/internal/api
GEN_DIR := $(API_DIR)/gen

# ----------------------------------------
# Generate code for one service
# ----------------------------------------
.PHONY: generate

generate:
	@echo "Generating API code for service: $(SERVICE)"
	@if [ ! -f "$(OPENAPI)" ]; then \
		echo "Error: $(OPENAPI) not found!"; \
		exit 1; \
	fi
	mkdir -p $(GEN_DIR)
	oapi-codegen \
		--config $(SERVER_CONFIG) \
		$(OPENAPI)
	oapi-codegen \
		--config $(CLIENT_CONFIG) \
		$(OPENAPI)
	@echo "Done: $(GEN_DIR)/server.gen.go $(GEN_DIR)/client.gen.go"

# ----------------------------------------
# Generate code for ALL services
# ----------------------------------------
SERVICES := $(shell ls api)

.PHONY: generate-all
generate-all:
	@echo "Generating API code for ALL services: $(SERVICES)"
	@for svc in $(SERVICES); do \
		$(MAKE) generate SERVICE=$$svc; \
	done
	@echo "All services updated."

# ----------------------------------------
# Run the auth service with different environments
# ----------------------------------------
auth-dev:
	cd services/auth && APP_ENV=dev go run ./cmd/main.go

auth-stage:
	cd services/auth && APP_ENV=stage go run ./cmd/main.go

auth-prod:
	cd services/auth && APP_ENV=prod go run ./cmd/main.go