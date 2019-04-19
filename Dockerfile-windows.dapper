FROM golang:1.12-windowsservercore
SHELL ["powershell", "-NoLogo", "-Command", "$ErrorActionPreference = 'Stop'; $ProgressPreference = 'SilentlyContinue';"]

RUN [Environment]::SetEnvironmentVariable('PATH', ('c:\innoextract;c:\app;{0}' -f $env:PATH), [EnvironmentVariableTarget]::Machine)
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

ENV DAPPER_ENV REPO TAG DRONE_TAG
ENV DAPPER_SOURCE /gopath/src/github.com/rancher/dapper
ENV DAPPER_OUTPUT ./bin
ENV DAPPER_DOCKER_SOCKET true
ENV TRASH_CACHE ${DAPPER_SOURCE}/.trash-cache
ENV HOME ${DAPPER_SOURCE}

WORKDIR ${DAPPER_SOURCE}
ENTRYPOINT ["powershell", "-NoLogo", "-NonInteractive", "-File", "./scripts/windows/entry.ps1"]
CMD ["ci"]