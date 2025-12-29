# make/sqlc.mk
DB ?= mysql
.PHONY: sqlc sqlc-all

sqlc: ## Generate sqlc for one service (SERVICE=... DB=mysql)
	@if [ -d "services/$(SERVICE)/db/$(DB)" ]; then \
		cd services/$(SERVICE)/db/$(DB) && sqlc generate; \
	else \
		echo "Skip SQLC: services/$(SERVICE)/db/$(DB) not found"; \
	fi

sqlc-all: ## Generate sqlc for all services (tries each service's db/*)
	@set -e; \
	for svc in $(SERVICES); do \
		for dir in services/$$svc/db/*; do \
			if [ -d "$$dir" ]; then \
				db=$$(basename "$$dir"); \
				$(MAKE) sqlc SERVICE=$$svc DB=$$db; \
			fi; \
		done; \
	done
