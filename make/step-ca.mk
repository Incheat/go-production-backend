.PHONY: init-step-ca

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
