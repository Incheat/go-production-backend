# make/security.mk
.PHONY: gosec

gosec: ## Run gosec
	@echo "==> Running gosec..."
	gosec -exclude-generated ./...