# make/helm.mk
# ----------------------------------------
# Helm
# ----------------------------------------

HELM_RELEASE ?= monorepo-dev
HELM_NS ?= dev
HELM_CHART ?= ./deploy/helm/monorepo
HELM_VALUES ?= -f ./deploy/helm/monorepo/values.yaml -f ./deploy/helm/monorepo/values.dev.secrets.yaml

.PHONY: helm-install helm-uninstall helm-template

helm-install: ## Install/upgrade helm chart
	@echo "Installing Helm chart: $(HELM_RELEASE) into ns: $(HELM_NS)"
	helm upgrade --install $(HELM_RELEASE) $(HELM_CHART) \
		-n $(HELM_NS) --create-namespace \
		$(HELM_VALUES)

helm-uninstall: ## Uninstall helm release
	@echo "Uninstalling Helm release: $(HELM_RELEASE) from ns: $(HELM_NS)"
	helm uninstall $(HELM_RELEASE) -n $(HELM_NS)

helm-template: ## Render helm templates (dry-run)
	helm template $(HELM_RELEASE) $(HELM_CHART) -n $(HELM_NS) $(HELM_VALUES)
