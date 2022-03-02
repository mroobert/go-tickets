SHELL := /bin/bash

run-default: 
	go run cmd/server/main.go

# ==============================================================================
# Modules support
tidy:
	go mod tidy
	go mod vendor

# ==============================================================================
# Building containers
VERSION := develop

auth-api:
	docker build \
		-f docker/Dockerfile.auth-api \
		-t auth-api:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		.