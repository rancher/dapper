ARG SERVERCORE_VERSION

FROM mcr.microsoft.com/windows/servercore:${SERVERCORE_VERSION}
SHELL ["powershell", "-NoLogo", "-Command", "$ErrorActionPreference = 'Stop'; $ProgressPreference = 'SilentlyContinue';"]
RUN [Environment]::SetEnvironmentVariable('PATH', ('C:\git\cmd;C:\git\mingw64\bin;C:\git\usr\bin;c:\innoextract;c:\app;c:\rancher;{0}' -f $env:PATH), [EnvironmentVariableTarget]::Machine)
RUN $URL = 'https://github.com/git-for-windows/git/releases/download/v2.21.0.windows.1/MinGit-2.21.0-64-bit.zip'; \
    \
    Write-Host ('Downloading git from {0} ...' -f $URL); \
    \
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12; \
    Invoke-WebRequest -UseBasicParsing -OutFile c:\git.zip -Uri $URL; \
    \
    Write-Host 'Expanding ...'; \
    \
    Expand-Archive c:\git.zip -DestinationPath c:\git\.; \
    \
    Write-Host 'Cleaning ...'; \
    \
    Remove-Item -Force -Path c:\git.zip; \
    \
    Write-Host 'Complete.';
RUN $URL = 'http://constexpr.org/innoextract/files/innoextract-1.7-windows.zip'; \
    \
    Write-Host ('Downloading innoextract from {0} ...' -f $URL); \
    \
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12; \
    Invoke-WebRequest -UseBasicParsing -OutFile c:\innoextract.zip -Uri $URL; \
    \
    Write-Host 'Expanding ...'; \
    \
    Expand-Archive c:\innoextract.zip -DestinationPath c:\innoextract\.; \
    \
    Write-Host 'Cleaning ...'; \
    \
    Remove-Item -Force -Recurse -Path c:\innoextract.zip; \
    \
    Write-Host 'Complete.';
RUN $URL = 'https://github.com/docker/toolbox/releases/download/v18.09.3/DockerToolbox-18.09.3.exe'; \
    \
    Write-Host ('Downloading docker from {0} ...' -f $URL); \
    \
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12; \
    Invoke-WebRequest -UseBasicParsing -OutFile c:\dockertoolbox.exe -Uri $URL; \
    \
    Write-Host 'Expanding ...'; \
    \
    pushd c:\; \
    \
    innoextract c:\dockertoolbox.exe; \
    \
    Write-Host 'Cleaning ...'; \
    \
    Remove-Item -Force -Recurse -Path @('c:\dockertoolbox.exe', 'c:\app\*') -Exclude @('docker.exe'); \
    \
    popd; \
    \
    Write-Host 'Complete.';
ADD bin/dapper.exe c:/rancher/dapper.exe