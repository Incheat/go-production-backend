.PHONY: init-step-ca step-ca-up step-ca-down init-step-ca-service init-step-ca-services

init-step-ca: ## Initialize step-ca
	mkdir -p infra/security/mtls/pki/stepca infra/security/mtls/pki/secrets
	test -f infra/security/mtls/pki/secrets/ca.password || openssl rand -base64 32 > infra/security/mtls/pki/secrets/ca.password
	docker run --rm -it \
		-v "$$PWD/infra/security/mtls/pki/stepca:/home/step" \
		-v "$$PWD/infra/security/mtls/pki/secrets:/run/secrets:ro" \
		smallstep/step-ca:latest \
		step ca init \
			--name "local-ca" \
			--dns "step-ca" \
			--dns "localhost" \
			--address ":9000" \
			--provisioner "admin" \
			--password-file "/run/secrets/ca.password"

step-ca-up:
	@if [ -z "$$(docker compose ps -q step-ca)" ]; then \
		echo "[step-ca] starting..."; \
		docker compose up -d step-ca; \
	else \
		echo "[step-ca] already running"; \
	fi

step-ca-down:
	@echo "[step-ca] stopping..."
	@docker compose stop step-ca >/dev/null 2>&1 || true

## Warning: step-ca service must be running before initializing the service
init-step-ca-service: ## Initialize step-ca service
	mkdir -p infra/security/mtls/pki/$(SERVICE)-sds
	mkdir -p infra/security/mtls/pki/secrets
	test -f infra/security/mtls/pki/secrets/$(SERVICE).sds.password || openssl rand -base64 32 > infra/security/mtls/pki/secrets/$(SERVICE).sds.password
	docker exec -i $$(docker compose ps -q step-ca) step ca root > infra/security/mtls/pki/$(SERVICE)-sds/root.crt

# 1. Check the default network
#
# docker network ls | grep _default
#
# 2. Run step-sds(e.g. auth) init command in the default network(e.g. go-production-backend_default)
#
# docker run --rm -it --network go-production-backend_default \
# 	-v "$PWD/infra/security/mtls/pki/auth-sds:/home/step" \
# 	smallstep/step-sds:latest \
# 	step-sds init --ca-url "https://step-ca:9000" --root "/home/step/root.crt"
#
# ✔ What would you like to name your new PKI? (e.g. SDS): auth-sds
# ✔ What do you want your PKI password to be? [leave empty and we'll generate one]: $(SERVICE).sds.password -> use the password generated in init-step-ca-service
# ✔ What address will your new SDS server listen at? (e.g. :443): :8443
# ✔ What DNS names or IP addresses would you like to add to your SDS server? (e.g. sds.smallstep.com[,1.1.1.1,etc.]): auth-sds
# ✔ What would you like to name your SDS client certificate? (e.g. envoy.smallstep.com): auth-envoy
# ✔ What do you want your certificates password to be? [leave empty and we'll generate one]: 
#
# 3. Close the step-sds service
#
# make step-ca-down
#
# 4. Revise sds.json file