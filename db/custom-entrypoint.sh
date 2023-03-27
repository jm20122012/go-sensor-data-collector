#!/bin/bash
set -e

# Run the original entrypoint script in the background
docker-entrypoint.sh postgres &

# Set the database connection URL
export DB_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable"

# Wait for PostgreSQL to be ready
until pg_isready -d "$DB_URL" 2>/dev/null; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 1
done

# Apply database migrations
echo "Running database migrations..."
atlas migrate apply --dir file://migrations --url "${DB_URL}"

# Check to seed if the device_type_ids table is empty. If it is, seed the database
TABLE_ROW_COUNT=$(psql -d "$DB_URL" -t -A -c "SELECT COUNT(id) FROM device_type_ids;")

if [ "$TABLE_ROW_COUNT" -eq 0 ]; then
  echo "Seeding the database..."
  psql -d "$DB_URL" -f seed.sql
fi

# Unset the password environment variable
unset POSTGRES_PASSWORD

# Create a flag file indicating initialization is complete
touch /tmp/db_initialized

# Wait indefinitely to keep the container running
wait
