#!/bin/sh
set -eu

project_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$project_root"

# Check Docker before creating any local credentials.
docker compose version >/dev/null
docker info >/dev/null

secret_dir="$project_root/.secrets/dev"
umask 077
mkdir -p "$secret_dir"
for name in postgres-password offensive-token setup-token; do
    secret_path="$secret_dir/$name"
    if [ -e "$secret_path" ] || [ -L "$secret_path" ]; then
        if [ ! -f "$secret_path" ] || [ ! -s "$secret_path" ]; then
            printf 'Invalid credential file: %s\n' "$secret_path" >&2
            exit 1
        fi
        continue
    fi
    # /dev/urandom, od and tr are available on both Linux and macOS.
    secret=$(od -An -N32 -tx1 /dev/urandom | tr -d ' \n')
    if [ "${#secret}" -ne 64 ]; then
        printf 'Could not generate development credentials.\n' >&2
        exit 1
    fi
    # Refuse to overwrite a file created by another concurrent start.
    (set -C; printf '%s' "$secret" > "$secret_path")
    unset secret
done

exec docker compose -f docker-compose.dev.yml up -d --build
