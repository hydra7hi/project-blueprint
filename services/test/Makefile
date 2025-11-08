# Makefile for gRPC User Service with PostgreSQL

# Service Name
SERVICE=user

# Base service Makefile
include ../Makefile.service

.PHONY: test-system

# Test cleanup
test-cleanup:
	@echo "Cleaning test databases..."
	@go run cleanup/main.go

# Run integration tests (with cleanup)
test-integration: test-cleanup
	@echo "Running integration tests..."
	make test-tags ARGS="integration"

# Run system tests with cleanup
test-system: test-cleanup
	@echo "Running system tests..."
	make test-tags ARGS="system"
