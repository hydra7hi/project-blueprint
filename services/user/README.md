# User Service
A Go gRPC service that performs basic CRUD operations on a PostgreSQL User database.

## Features
- Create, Read, Update, Delete users
- List users with pagination
- PostgreSQL integration
- Environment variable configuration

# API:
The service provide CRUD operations on the user DB,
with validation for request

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
make script-create               # create test user
make script-create ARGS="test"   # create user with specific name
make script-get                  # get user
make script-get ARGS="1"         # get user with specific id
make script-list                 # list users
make script-list ARGS="1 2"      # list users with page number and limit
make script-update               # update user
make script-update ARGS="1"      # update user with specific id
make script-delete               # delete user
make script-delete ARGS="1"      # delete user with specific id
```
