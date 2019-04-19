#Requires -Version 5.0
$ErrorActionPreference = "Stop"

trap {
    Write-Host -ForegroundColor DarkRed "$_"

    exit 1
}

Invoke-Expression -Command "$PSScriptRoot\ci.ps1"
Invoke-Expression -Command "$PSScriptRoot\docker-image.ps1"
