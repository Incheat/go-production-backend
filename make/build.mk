# make/docker.mk
# ----------------------------------------
# Docker buildx (ARM64) + kind load
# ----------------------------------------

ENV ?= dev
KIND_NAME ?= dev

.PHONY: build-local build-all build-all-kind build-probe

build-probe: ## Build probe image
	@echo "Building probe image"
	docker buildx build -f infra/obs/probe/Dockerfile \
		-t probe:$(ENV) \
		--platform linux/arm64 \
		--load .

build-local: ## Build one service image (SERVICE=... ENV=...) (ARM64 only)
	@echo "Building image for $(SERVICE) (linux/arm64 ONLY)"
	docker buildx build -f services/$(SERVICE)/Dockerfile \
		--build-arg SERVICE=$(SERVICE) \
		--platform linux/arm64 \
		-t $(SERVICE):$(ENV) \
		--load .

build-all: ## Build all service images (ARM64 only)
	@echo "Building image for all services (ARM64 ONLY)"
	$(call for_each_service,build-local)
	$(MAKE) build-probe

build-all-kind-load: ## Build all images and load into kind (ARM64 only)
	@echo "Building image for all services (ARM64 ONLY) and loading into kind cluster: $(KIND_NAME)"
	@set -e; \
	for svc in $(SERVICES); do \
		$(MAKE) build-local SERVICE=$$svc ENV=$(ENV); \
		kind load docker-image --name $(KIND_NAME) $$svc:$(ENV); \
	done
