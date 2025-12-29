# make/common.mk

SHELL := /bin/bash

SERVICE ?= helloworld
SERVICES_DIR ?= services
API_DIR ?= api

# More stable than ls: only take directories
SERVICES := $(shell find $(API_DIR) -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)

define for_each_service
	@set -e; \
	for svc in $(SERVICES); do \
		$(MAKE) $(1) SERVICE=$$svc; \
	done
endef

help: ## Show commands
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "  %-22s %s\n", $$1, $$2}'
