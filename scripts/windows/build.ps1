#Requires -Version 5.0
$ErrorActionPreference = "Stop"

trap {
    Write-Host -ForegroundColor DarkRed "$_"

    exit 1
}

Invoke-Expression -Command "$PSScriptRoot\version.ps1"

$DIR_PATH = Split-Path -Parent $MyInvocation.MyCommand.Definition
$SRC_PATH = (Resolve-Path "$DIR_PATH\..\..").Path
cd $SRC_PATH

$null = New-Item -Type Directory -Path bin -ErrorAction Ignore
$env:GOARCH=$env:ARCH
$env:GOOS='windows'
$env:CGO_ENABLED=0
$LINKFLAGS = ('-X main.VERSION={0} -s -w -extldflags "-static"' -f $env:VERSION)
go build -ldflags $LINKFLAGS -o .\bin\dapper.exe main.go