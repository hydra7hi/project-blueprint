#!/bin/sh

set -e  # Exit on any error

ENV_FILE="./.env"

# Check if .env file exists
if [ ! -f "$ENV_FILE" ]; then
    echo "Error: .env file not found at $ENV_FILE"
    exit 1
fi

# Check if .env file is readable
if [ ! -r "$ENV_FILE" ]; then
    echo "Error: .env file is not readable"
    exit 1
fi

# Source the .env file
echo "Loading environment variables from $ENV_FILE"
set -a
. "$ENV_FILE"
set +a

# Validate that required variables are set
required_vars="DB_HOST DB_PORT DB_USER DB_PASSWORD DB_NAME"

for var in $required_vars; do
    eval "value=\"\$$var\""
    if [ -z "$value" ]; then
        echo "Warning: Required variable $var is not set in .env"
    else
        echo "$var is set"
    fi
done

echo "Environment variables successfully loaded"