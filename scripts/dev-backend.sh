#!/bin/sh
set -eu

# Read local development credentials without putting them in command arguments.
export ARTEX_PG_DSN="postgres://artex:$(cat /run/secrets/postgres-password)@postgres:5432/artex?sslmode=disable"
export ARTEX_SETUP_TOKEN="$(cat /run/secrets/setup-token)"
go build -mod=readonly -o /app/artex ./cmd/artex
cd /app
exec /app/artex -addr :8787 -proxy :8788 -data /app/data
