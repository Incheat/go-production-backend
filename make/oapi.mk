# make/oapi.mk
.PHONY: oapi-generate oapi-generate-all

GLOBAL_OAPI_DIR := api/$(SERVICE)/oapi

PUBLIC_SERVER_CONFIG := $(GLOBAL_OAPI_DIR)/public.server.yaml
INTERNAL_SERVER_CONFIG := $(GLOBAL_OAPI_DIR)/internal.server.yaml
PUBLIC_CLIENT_CONFIG := $(GLOBAL_OAPI_DIR)/public.client.yaml
INTERNAL_CLIENT_CONFIG := $(GLOBAL_OAPI_DIR)/internal.client.yaml
PUBLIC_TYPES_CONFIG := $(GLOBAL_OAPI_DIR)/public.type.yaml
INTERNAL_TYPES_CONFIG := $(GLOBAL_OAPI_DIR)/internal.type.yaml
PUBLIC_OPENAPI := $(GLOBAL_OAPI_DIR)/public.yaml
INTERNAL_OPENAPI := $(GLOBAL_OAPI_DIR)/internal.yaml

oapi-generate: ## Generate OpenAPI code for one service (SERVICE=...)
	@echo "==> Generating OpenAPI code for service: $(SERVICE) ==="
	@if [ -f "$(PUBLIC_OPENAPI)" ]; then \
		oapi-codegen --config $(PUBLIC_SERVER_CONFIG) $(PUBLIC_OPENAPI); \
		oapi-codegen --config $(PUBLIC_CLIENT_CONFIG) $(PUBLIC_OPENAPI); \
		oapi-codegen --config $(PUBLIC_TYPES_CONFIG) $(PUBLIC_OPENAPI); \
		echo "✓ PUBLIC OpenAPI code generation completed"; \
	else \
		echo "⚠ Skip PUBLIC: $(PUBLIC_OPENAPI) not found"; \
	fi
	@if [ -f "$(INTERNAL_OPENAPI)" ]; then \
		oapi-codegen --config $(INTERNAL_SERVER_CONFIG) $(INTERNAL_OPENAPI); \
		oapi-codegen --config $(INTERNAL_CLIENT_CONFIG) $(INTERNAL_OPENAPI); \
		oapi-codegen --config $(INTERNAL_TYPES_CONFIG) $(INTERNAL_OPENAPI); \
		echo "✓ INTERNAL OpenAPI code generation completed"; \
	else \
		echo "⚠ Skip INTERNAL: $(INTERNAL_OPENAPI) not found"; \
	fi

oapi-generate-all: ## Generate OpenAPI code for all services
	$(call for_each_service,oapi-generate)
