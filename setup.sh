#!/bin/bash
set -e

echo "setting up gymapp architecture"
docker compose up --build -d

echo "waiting for postgres be ready"
until docker exec -it postgres pg_isready -U admin > /dev/null 2>&1; do
  sleep 3
done

echo "creating exercises table"
docker exec -i postgres psql -U admin -d gymapp -c "
CREATE TYPE exercise_type AS ENUM ('legs', 'arms', 'core');"
docker exec -i postgres psql -U admin -d gymapp -c "
CREATE TABLE IF NOT EXISTS exercises (
    id UUID PRIMARY KEY,
    type exercise_type NOT NULL,
    alias TEXT NOT NULL,
    description TEXT NOT NULL
);"

echo "the architecture is ready"