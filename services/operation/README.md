# Operation Service
A Go gRPC service that allows starting a long running operation in the background, and check on the status of it later.

Currently it include an example usecase to queue a job that Creates Users in user service.

The service is still a work in progress

## Features
- Start, check, and cancel long-running operations
- Get operation results when they're done
- Background job management
- Postgress integration for storing operations
- Environment variable configuration

# API:
The service allows starting a long running operation in the background.
Then sending a request to check the status.

For more details check:
- [Proto](./proto/operation.proto) For usage.
- [Handler](./server/handler.go) For returned error codes.

# Database
Keeps track of operations and their progress, like a recipe book for your background jobs.

**Table Structure:**
- `id`: id number for each operation
- `marshalled_request`: Metadata for the operation
- `step_id`: Which step is the operation currently at
- `state` 
- `created_at`: Record creation timestamp (auto-set)
- `updated_at`: Last update timestamp (auto-updated)

# Testing:
- Unit tests:
```bash
make test
# or
make test-coverage
```

- Manual test
```bash
make script-start               # starts the LRO operation
make script-check               # checks the LRO operation state, until it finishs.
```
