$ErrorActionPreference = "Stop"

Set-Location (Split-Path -Parent $PSScriptRoot)

docker compose up -d --build
if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}

docker compose ps

Write-Host "TravelRaft is running at http://localhost:8081"
