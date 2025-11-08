# Makefile for gRPC User Service with PostgreSQL

# Service Name
SERVICE=user

# Base service Makefile
include ../Makefile.service


script-create:
	@go run ./scripts/create/main.go $(ARGS)

script-delete:
	@go run ./scripts/delete/main.go $(ARGS)

script-get:
	@go run ./scripts/get/main.go $(ARGS)

script-list:
	@go run ./scripts/list/main.go $(ARGS)

script-update:
	@go run ./scripts/update/main.go $(ARGS)