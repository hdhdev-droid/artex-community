#!/bin/sh
set -eu

# This overlay upgrades an initialized local installation, never a new one.
if [ ! -s /app/auth/jwt.key ]; then
  echo "Existing local authentication volume is missing; refusing to create a new JWT key." >&2
  exit 1
fi

# Preserve the interpreter path stored by the previous Debian-based dev image.
if [ ! -e /usr/bin/python3 ] && [ -x /usr/local/bin/python3 ]; then
  ln -s /usr/local/bin/python3 /usr/bin/python3
fi

# Keep credentials in mounted files; never echo their values or compile sources.
export ARTEX_PG_DSN="postgres://artex:$(cat /run/secrets/postgres-password)@postgres:5432/artex?sslmode=disable"
export ARTEX_SETUP_TOKEN="$(cat /run/secrets/setup-token)"
exec /app/start.sh -addr :8787 -proxy :8788 -data /app/data
