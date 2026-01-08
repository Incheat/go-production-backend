.DEFAULT_GOAL := help

include make/common.mk
include make/go.mk
include make/oapi.mk
include make/grpc.mk
include make/sqlc.mk
include make/migrate.mk
include make/build.mk
include make/helm.mk
include make/docker-compose.mk
include make/run.mk
include make/security.mk
include make/step-ca.mk
include make/openssl.mk