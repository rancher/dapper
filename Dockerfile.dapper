FROM golang:1.11.4

env ARCH amd64

RUN apt-get update && \
    apt-get install -y apt-transport-https ca-certificates

RUN wget -O - https://download.docker.com/linux/ubuntu/gpg | apt-key add - && \
    echo "deb [arch=${ARCH}] https://download.docker.com/linux/ubuntu xenial stable" >> /etc/apt/sources.list && \
    cat /etc/apt/sources.list

RUN apt-get update && \
    apt-get install -y 'docker-ce=5:18.09.1~3-0~ubuntu-xenial' bash git jq

ENV DOCKER_CLI_EXPERMENTAL enabled
ENV DAPPER_SOURCE /go/src/github.com/rancher/dapper
ENV DAPPER_OUTPUT ./bin ./dist
ENV DAPPER_DOCKER_SOCKET true
ENV TRASH_CACHE ${DAPPER_SOURCE}/.trash-cache
ENV DAPPER_ENV CROSS
ENV HOME ${DAPPER_SOURCE}
WORKDIR ${DAPPER_SOURCE}

ENTRYPOINT ["./scripts/entry"]
CMD ["ci"]
