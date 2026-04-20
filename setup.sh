#!/bin/bash
set -e

echo "setting up gymapp architecture"
docker compose up --build -d

echo "waiting for postgres be ready"
until docker exec -it postgres pg_isready -U admin > /dev/null 2>&1; do
  sleep 3
done

echo "Running database migrations..."
# We use the official migrate image to run the migrations against the postgres container.
NETWORK_NAME=$(docker network ls --filter "name=${PWD##*/}" -q | head -n 1)
docker run --rm \
    -v $(pwd)/migrations:/migrations \
    --network "$NETWORK_NAME" \
    migrate/migrate \
    -path=/migrations/ \
    -database "postgres://postgres:postgres@postgres:5432/appweb?sslmode=disable" \
    up

echo "the architecture is ready"