SHELL := /bin/bash

run-default: 
	go run cmd/server/main.go

# ==============================================================================
# Modules support
tidy:
	go mod tidy
	go mod vendor