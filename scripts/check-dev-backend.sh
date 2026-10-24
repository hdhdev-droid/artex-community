#!/bin/sh
set -eu
# This hostname is a separate, disposable Compose test service, never the app DB.
export ARTEX_CONFIG=/app/nonexistent-config.json
export ARTEX_LANG=ko
result=0
go test -count=1 ./defenseeval ./investigation ./securitysource ./closedloop ./offensivereview ./report ./notify ./internal/l10n ./internal/chatgpt ./llmpool ./selfupdate > /tmp/artex-backend-check.log 2>&1 || result=$?
# Existing suites leave global settings and task fixtures behind. Give each
# package its own database so a server startup cannot resume another suite's tasks.
for package in db agent server; do
    export ARTEX_PG_DSN="postgres://artex:$(cat /run/secrets/postgres-password)@postgres-test:5432/artex_test_${package}?sslmode=disable"
    go test -count=1 "./$package" >> /tmp/artex-backend-check.log 2>&1 || result=$?
done
# Auth tests may log temporary setup tokens. Only expose package/test status.
grep -E '^(ok[[:space:]]|FAIL|--- FAIL:|\?[[:space:]])' /tmp/artex-backend-check.log || true
exit "$result"
