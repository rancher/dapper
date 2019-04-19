#Requires -Version 5.0
$ErrorActionPreference = "Stop"

trap {
    Write-Host -ForegroundColor DarkRed "$_"

    exit 1
}

$SCRIPT_PATH = ("{0}\{1}.ps1" -f $PSScriptRoot, $Args[0])
if (Test-Path $SCRIPT_PATH -ErrorAction Ignore) {
    Invoke-Expression -Command $SCRIPT_PATH
    exit
}

Start-Process -Wait -FilePath $Args[0] -ArgumentList $Args[1..$Args.Length]