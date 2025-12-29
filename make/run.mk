# make/run.mk
# ----------------------------------------
# Run microservices locally
# ----------------------------------------

ENV_FILE ?= .env.local

.PHONY: run-local run-all-local

run-local: ## Run one service locally (SERVICE=...)
	@echo "$(SERVICE) service running in dev mode"
	@set -a; [ -f "$(ENV_FILE)" ] && . "$(ENV_FILE)"; set +a; \
	go run ./services/$(SERVICE)/cmd/main.go

run-all-local: ## Run all services locally (background). Stop with Ctrl+C
	@echo "Running all services in dev mode (background)"
	@set -e; \
	set -a; [ -f "$(ENV_FILE)" ] && . "$(ENV_FILE)"; set +a; \
	pids=""; \
	for svc in $(SERVICES); do \
		echo "==> $$svc"; \
		go run ./services/$$svc/cmd/main.go & \
		pids="$$pids $$!"; \
	done; \
	trap 'echo "Stopping..."; kill $$pids 2>/dev/null || true' INT TERM; \
	wait $$pids
