#!/bin/bash

set -a
source ./../env
set +a

required_vars=("DB_HOST" "DB_PORT" "DB_USER" "DB_PASSWORD" "DB_NAME" "GRPC_PORT")

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "Error: Required environment variable $var is not set"
        echo "Please set all required variables in .env file:"
        echo "DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, GRPC_PORT"
        exit 1
    fi
done

echo "All required environment variables are set"