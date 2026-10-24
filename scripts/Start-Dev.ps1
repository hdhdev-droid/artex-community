$ErrorActionPreference = 'Stop'
$projectRoot = Split-Path $PSScriptRoot -Parent
$secretDir = Join-Path $projectRoot '.secrets/dev'
New-Item -ItemType Directory -Path $secretDir -Force | Out-Null

# Preserve credentials on subsequent starts. Never print secret contents.
foreach ($name in @('postgres-password', 'offensive-token', 'setup-token')) {
    $path = Join-Path $secretDir $name
    if (-not (Test-Path -LiteralPath $path)) {
        $bytes = New-Object byte[] 32
        $rng = [System.Security.Cryptography.RandomNumberGenerator]::Create()
        try { $rng.GetBytes($bytes) } finally { $rng.Dispose() }
        $value = [BitConverter]::ToString($bytes).Replace('-', '').ToLowerInvariant()
        [System.IO.File]::WriteAllText($path, $value, [System.Text.UTF8Encoding]::new($false))
    }
}

Push-Location $projectRoot
try {
    docker compose -f docker-compose.dev.yml up -d --build
    if ($LASTEXITCODE -ne 0) { throw 'Development containers could not be started.' }
} finally {
    Pop-Location
}
