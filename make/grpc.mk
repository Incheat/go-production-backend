# make/grpc.mk
# ----------------------------------------
# Generate GRPC code
# ----------------------------------------

.PHONY: protoc protoc-all

# depends on SERVICE
GRPC_GEN_DIR := api/$(SERVICE)/grpc/gen
INTERNAL_PROTO := api/$(SERVICE)/grpc/internal.proto

protoc: ## Generate grpc code for one service (SERVICE=...)
	@echo "=== Generating GRPC code for service: $(SERVICE) ==="
	@if [ -f "$(INTERNAL_PROTO)" ]; then \
		mkdir -p "$(GRPC_GEN_DIR)"; \
		protoc --go_out="$(GRPC_GEN_DIR)" --go-grpc_out="$(GRPC_GEN_DIR)" "$(INTERNAL_PROTO)"; \
		echo "✓ GRPC code generation completed"; \
	else \
		echo "⚠ Skipping GRPC: $(INTERNAL_PROTO) not found."; \
	fi

protoc-all: ## Generate grpc code for all services
	$(call for_each_service,protoc)
