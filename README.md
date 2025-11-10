# Intro:
The following project was created as a show case for golang code, for my job application as a golang developer.

I wanted to include some of the basic components that I thought would be useful for any Go project. Such as:
- Golang.
- Docker.
- Make.
- Testing:
  - Unit tests.
  - Integration tests.
  - System tests.
  - Manual test scripts.
- Grpc and Proto messages.
- Postgresql Databases.
- ...


As I had limited time to build this project from scratch, the code is not perfect in the way I would usually have it in, and is definitely not production ready.
Rather to serves as an example of one of the ways a go project might be build.

As I proceed with the implementation to make the code more robust. I plan to include more components, that I have identified in [future improvements](#future-improvements) paragraph. That can be found at the end of the document.

## Contents:

1. [Prerequisites](#prerequisites).
1. [Quick Start](#quick-start).
1. [Services](#services).
1. [Databases](#databases).
1. [Tests](#tests).
1. [Architecture](#architecture).
1. [Make](#make-commands).
1. [Demo](#demo).
1. [Future improvements](#future-improvements).

## Prerequisites
- Go 1.19+:
https://go.dev/doc/install
- PostgreSQL:
https://www.postgresql.org/download/
- Docker:
https://docs.docker.com/get-started/get-docker/
- protoc compiler:
https://protobuf.dev/installation/
- Make:
https://www.gnu.org/software/make/#download

## Quick Start:
To quickly explore the project, you only need to build and run the docker compose image.

This can be achieved by running the following in `./services/`:

```bash
# Build all docker compose images
make build

# Run all docker compose
make up

# Check logs
make logs
```

When done, you can stop all
```bash
# Stop all containers
make down
```

# Services
The example provides 2 basic services each with its own Database:
- User service:
A basic user management example with CRUD operations. [User README.md](./services/user/README.md)

- Operation service: 
A basic service for queuing LRO jobs. [Operation README.md](./operation/README.md)

# Databases
Each of the previosly mentioned servcies comes with a postgress database. More details about it in each services Readme.

# Tests
The project contains multiple tests in 4 categories:
1. Manual test scripts:

Golang code to manually test each rpc. The script can be found under `./services/<service_name>/scripts/`
and can be quickly started using.
```bash
make script-*           # Each service has different test cases
```

The test scripts expect the service to be up and running at the default port.

2. Unit tests: Automated unit tests to test the code.
3. Integrations tests: Automated test to test a service and db.
4. System tests: Automated test to test all services working together.

## Architecture
The project follows the following Architecture. It was mainly constructed with the aim of easy extendability in mind.
```
services/          # All Golang Services
├── user/               # User management service       (port 50051)
├── operation/          # Background operations service (port 50052)
├── <more_services>     # More services can be added to the example as needed
├── test/               # Integration & system tests for features between services.
├── .env.example        # Environment variables example
├── go.mod              # Go module configuration
├── go.work             # Go workspace configuration
├── Makefile            # Docker Build and run commands
└── Makefile.service    # Common Build/test commands for each service
```

## Make Commands
The project is equipped with basic Make commands, that helps building, running and maintaining the code.

Each Directory will have different make commands.
You can run `Make help` to get more insights of the provided functionalities.

Basic dirs to run make commands are:
- `./Services/` Here its possible to run make commands to:
  - Startup the whole project in docker compose container.
```bash
make build                     # build all services and dbs
make up                        # start all services and dbs
make down                      # stop all services and dbs
make logs                      # shows logs from all services and dbs
```
  - Connect to specific services:
```bash
make user-shell                # Connect to user service shell
make operation-shell           # Connect to operation service shell
```
  - Running any other make command, will proxy the command to be executed in all child services directories:
```bash
make test                      # Runs `make test` in user then operation service.
```

- `./Services/<service_name>/` Here its possible to run make commands to
  - Build and handle main go code for this specific service..
```bash
make all                   # Install deps, generate proto, and build
make build                 # Build the application
make deps                  # Install dependencies
make proto                 # Generate protobuf code
make clean                 # Clean build artifacts
```
  - Assure code quality:
```bash
make fmt                   # Format code
make vet                   # Vet code
make lint                  # Lint code (requires golangci-lint)
```
  - Run service unit tests:
```bash
make test                  # Run unit tests
make test-coverage         # Run tests with coverage report
```


# Demo:
Recommended flow to show case all features would be:
1. Start the project:
```bash
cd ./services/
make build
make up
make logs
```

2. Build a service.
```bash
cd ./services/user/
make all
```

3. Run service unit tests + check generated report.
```bash
make test-coverage
# Then check the generated coverage.html
```

4. Run manual test scripts
```bash
cd ./services/user/
make script-create               # create user
make script-get ARGS="1"         # get user with id
make script-list                 # list users
make script-update ARGS="1"      # update user with id
make script-delete ARGS="1"      # delete user with id
```

5. Run integration tests:
```bash
cd ./services/test
make test-integration
```

6. Run system tests:
```bash
cd ./services/test
make test-system
```

# Future improvements:
- Handle configurations properly.
- Create CI/CD Pipeline script.
  - Enforce code quality checks before being able to commit the code. (fmt and lint)
  - Add stages for build, test and deploy.
- Logging:
  - Have a logger interface.
  - Add proper monitoring solution.
- Auth.
- Refine Code and clean up:
  - Move packages to common place.
- Operation:
  - Store error code + message in DB.
  - Add data to be stored and read.
  - Add tests for restarting the service.
- User.
  - Move DB client to common place.
  - Add seeding scripts.
