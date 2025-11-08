# User Service
A Go gRPC service that performs basic CRUD operations on a PostgreSQL User database.

## Features
- Create, Read, Update, Delete users
- List users with pagination
- PostgreSQL integration
- Environment variable configuration

# API:
The service provide CRUD operations on the user DB,
with vaidation for request

For more details check:
- [Proto](./proto/user.proto) For usage.
- [Handler](./server/handler.go) For returned error codes.

# Database
Stores user information with automatic timestamp management.

**Table Structure:**
- `id`: Auto-incrementing primary key
- `name`: User's full name (required)
- `email`: Unique email address (required)
- `age`: User's age (required)s
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
make script-create               # create user
make script-get ARGS="1"         # get user with id
make script-list                 # list users
make script-update ARGS="1"      # update user with id
make script-delete ARGS="1"      # delete user with id
```
