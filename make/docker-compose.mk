# make/compose.mk
# ----------------------------------------
# docker compose
# ----------------------------------------

.PHONY: compose-up compose-down compose-reset compose-logs

compose-up: ## docker compose up -d
	docker compose up -d

compose-down: ## docker compose down
	docker compose down

compose-reset: ## docker compose down -v
	docker compose down -v

compose-logs: ## docker compose logs -f
	docker compose logs -f
