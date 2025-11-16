#!/bin/sh
set -e

echo "Waiting for database ${DB_HOST:-postgres}:${DB_PORT:-5432} with user ${POSTGRES_USER:-postgres}..."

until pg_isready -h ${DB_HOST:-postgres} -p ${DB_PORT:-5432} -U ${POSTGRES_USER:-postgres}; do
    echo "Database is unavailable - sleeping"
    sleep 1
done

echo "Database is ready!"

echo "Running database migrations..."
./main migrate

echo "Starting application..."
exec ./main