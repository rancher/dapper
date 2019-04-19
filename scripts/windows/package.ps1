#Requires -Version 5.0
$ErrorActionPreference = "Stop"

trap {
    Write-Host -ForegroundColor DarkRed "$_"

    exit 1
}

$DIR_PATH = Split-Path -Parent $MyInvocation.MyCommand.Definition
$SRC_PATH = (Resolve-Path "$DIR_PATH\..\..").Path
cd $SRC_PATH

$null = New-Item -Type Directory -Path dsit\artifacts -ErrorAction Ignore
$null = Copy-Item -Force -Path bin\dapper* -Destination dsit\artifacts
