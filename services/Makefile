DOCKER_COMPOSE = docker-compose
COMPOSE_FILE = docker-compose.yml
SERVICES := user operation  # add more as needed

# Proxy any command to ALL services
%:
	@echo "🔁 Running '$@' for all services: $(SERVICES)"
	@for s in $(SERVICES); do \
		echo "→ 🔵 Running '$@' for '$$s'"; \
		$(MAKE) -C $$s $@ || echo "⚠️  Command '$@' failed for $$s (might not exist)"; \
	done
	@echo "✅ Command '$@' executed for all services"


# Quick commands:
build:
	$(MAKE) user-build
	$(MAKE) operation-build
up:
	$(MAKE) user-up
	$(MAKE) operation-up
down:
	$(MAKE) compose-down
logs:
	$(MAKE) compose-logs-all


# Service-specific shortcuts
## User
user-up:
	$(MAKE) compose-up SERVICE=user

user-shell:
	$(MAKE) compose-shell SERVICE=user

user-logs:
	$(MAKE) compose-logs SERVICE=user

user-build:
	$(MAKE) compose-build SERVICE=user

## Operation
operation-up:
	$(MAKE) compose-up SERVICE=operation

operation-shell:
	$(MAKE) compose-shell SERVICE=operation

operation-logs:
	$(MAKE) compose-logs SERVICE=operation

operation-build:
	$(MAKE) compose-build SERVICE=operation

##################
## Deprecated:  ##
##################

# Commands on all targets:
compose-up-all:
	@echo "Starting all services (build happens once)..."
	# Single command to build and start everything, leveraging the shared build stage cache
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up --build -d
	@echo "✅ All services are up!"

compose-down-all:
	@echo "Stopping all services..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) down
	@echo "✅ All services are down!"

compose-build-all:
	@echo "Building all services..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) build
	@echo "✅ All services are docker built!"

compose-logs-all:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) logs -f

# Shared Compose commands:
compose-up:
	@echo "Starting $(SERVICE) with its DB..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up -d $(SERVICE)-db $(SERVICE)-service

compose-down:
	@echo "Stopping all services..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) down

compose-logs:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) logs -f $(SERVICE)-service

compose-build:
	@echo "Building $(SERVICE)-service..."
	# This correctly builds the target stage from the Dockerfile
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) build $(SERVICE)-service

compose-shell:
	@echo "Opening shell into $(SERVICE)-service..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) exec $(SERVICE)-service /bin/sh

help:
	@echo "# Quick Commands"
	@echo "make build                     # build all services and dbs"
	@echo "make up                        # start all services and dbs"
	@echo "make down                      # stop all services and dbs"
	@echo "make logs                      # shows logs from services and dbs"

# Default target
.DEFAULT_GOAL := help