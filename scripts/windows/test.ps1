#Requires -Version 5.0
$ErrorActionPreference = "Stop"

trap {
    Write-Host -ForegroundColor DarkRed "$_"

    exit 1
}

go test ./...