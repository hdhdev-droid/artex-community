$ErrorActionPreference = 'Stop'
$projectRoot = Split-Path $PSScriptRoot -Parent
Push-Location $projectRoot
try {
    # Only this disposable, non-published test container is recreated.
    docker compose -f docker-compose.dev.yml --profile test up -d --force-recreate --wait postgres-test
    if ($LASTEXITCODE -ne 0) { throw 'Could not prepare isolated test PostgreSQL.' }
    foreach ($package in @('db', 'agent', 'server')) {
        docker compose -f docker-compose.dev.yml exec -T postgres-test createdb -U artex "artex_test_$package"
        if ($LASTEXITCODE -ne 0) { throw "Could not create test database for $package." }
    }
    docker compose -f docker-compose.dev.yml exec -T backend sh /src/scripts/check-dev-backend.sh
    if ($LASTEXITCODE -ne 0) { throw 'Backend checks failed; see /tmp/artex-backend-check.log inside backend.' }
    Get-Content -Raw (Join-Path $PSScriptRoot 'check-dev-restore.sh') |
        docker compose -f docker-compose.dev.yml exec -T postgres-test sh -c "tr -d '\015' | sh"
    if ($LASTEXITCODE -ne 0) { throw 'Isolated database dump/restore verification failed.' }
} finally {
    Pop-Location
}
