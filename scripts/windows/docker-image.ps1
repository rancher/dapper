#Requires -Version 5.0
$ErrorActionPreference = "Stop"

trap {
    Write-Host -ForegroundColor DarkRed "$_"

    exit 1
}

$DIR_PATH = Split-Path -Parent $MyInvocation.MyCommand.Definition
$SRC_PATH = (Resolve-Path "$DIR_PATH\..\..").Path
cd $SRC_PATH

if (-not (Test-Path "bin\dapper")) {
    Invoke-Expression -Command "$DIR_PATH\build.ps1"
}

# Get release id as image tag suffix
$HOST_RELEASE_ID = (Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\' -ErrorAction Ignore).ReleaseId
$RELEASE_ID = $env:RELEASE_ID
if (-not $RELEASE_ID) {
    $RELEASE_ID = $HOST_RELEASE_ID
}

$IMAGE = ('rancher/dapper:windows-{0}' -f $RELEASE_ID)
if ($RELEASE_ID -eq $HOST_RELEASE_ID) {
    docker build `
        --build-arg SERVERCORE_VERSION=$RELEASE_ID `
        -t $IMAGE `
        -f Dockerfile-windows .
} else {
    docker build `
        --isolation hyperv `
        --build-arg SERVERCORE_VERSION=$RELEASE_ID `
        -t $IMAGE `
        -f Dockerfile-windows .
}
$null = New-Item -Type Directory -Path dsit -ErrorAction Ignore
$IMAGE | Out-File -Encoding ascii -Force -FilePath dsit\images