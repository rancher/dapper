# Dapper - Docker Build Wrapper

Dapper is a tool to wrap any existing build tool in an consistent environment.  This allows people to build your software from source or modify it without worrying about setting up a build environment.  The approach is very simple and taken from a common pattern that has adopted by many open source projects.  Create a file called `Dockerfile.dapper` in the root of your repository.  Dapper will build that Dockerfile and then execute a container based off of the resulting image.  Dapper will also copy in source files and copy out resulting artifacts or will use bind mounting if you choose.

## Installation

```sh
curl -sL https://releases.rancher.com/dapper/latest/dapper-`uname -s`-`uname -m` > /usr/local/bin/dapper
chmod +x /usr/local/bin/dapper
```

From source

```sh
go get github.com/rancher/dapper
```

## Example

Dapper is built using dapper so the following is a decent example

```sh
go get github.com/rancher/dapper
git clone https://github.com/rancher/dapper.git
cd dapper
dapper
```

This is the `Dockerfile.dapper` used

```Dockerfile
FROM golang:1.4
RUN go get github.com/tools/godep
ENV DAPPER_SOURCE /go/src/github.com/rancher/dapper
ENV DAPPER_OUTPUT bin
WORKDIR ${DAPPER_SOURCE}
ENTRYPOINT ["./script/build"]
```

## Using

### Dockerfile.dapper

The `Dockerfile.dapper` is intended to create a build environment but not really build your code.  For example if you need build tools such as `make` or `bundler` or language environments for Ruby, Python, Java, etc.

The `ENTRYPOINT`, `CMD`, and `WORKDIR` defined in `Dockerfile.dapper` are what are used to initiate your build.  When running `dapper foo bar`, `foo bar` will be passed as the docker CMD.  For example, running `dapper make install` will do the basic equivalent of `docker run -it --rm build-image make install`.  If you want you can set the `ENTRYPOINT` to `make` and then `dapper install` will be the same as `make install`.  Either approach is fine.

You can also customize your build container image with build arguments (via `ARG` Dockerfile instructions), which are populated from environment variables on dapper image build. That is useful if you want to parameterize your build for different platforms and you're using essentially the same build environment, only on different platforms. For example, if you have `ARG ARCH` in Dockerfile.dapper, you can have `ARCH=arm` in your environment variables, and when you run `dapper -s` your dapper image is built with `--build-arg ARCH=arm` and `$ARCH` is effectively replaced with `arm` in the resulting dapper image.

### Dapper Modes: Bind mount or CP

Dapper runs in two modes `bind` or `cp`, meaning bind mount in the source or cp in the source.  Depending on your environment one or the other could be preferred.  If your host is Linux bind mounting is typically preferred because it is very fast.  If you are running on Mac, Windows, or with a remote Docker daemon, CP is usually your only option.  You can force a specific mode with

    dapper --mode|m MODE

For example `dapper -m cp` or `dapper -m bind`.

### Interactive Shell

If you just want a shell in the build environment run `dapper -s`.

## Configuring

Configuring the behavior of Dapper is done through ENV variables in the `Dockerfile.dapper`.

### DAPPER_SOURCE

`DAPPER_SOURCE` is the location in the container of where your source should be.  For go applications this might look like `ENV DAPPER_SOURCE /go/src/github.com/rancher/dapper`

In bind mode `DAPPER_SOURCE` is used in the Docker `run` command as follows

    docker run -v .:${DAPPER_SOURCE} build-image

In CP mode `DAPPER_SOURCE` is used in a Docker `cp` command as follows

    docker cp . build-container:${DAPPER_SOURCE}

The default value of `DAPPER_SOURCE` is `/source`.

### DAPPER_CP

`DAPPER_CP` is the location in the host that should be copied to `DAPPER_SOURCE` in the container.

In bind mode `DAPPER_CP` is used in the Docker `run` command as follows

    docker run -v ${DAPPER_CP}:/source/ build-image

In CP mode `DAPPER_CP` is used in a Docker `cp` command as follows

    docker cp ${DAPPER_CP} build-container:/source/

The default value of `DAPPER_CP` is `.`.

### DAPPER_OUTPUT

`DAPPER_OUTPUT` is used after the build is done to copy the build artifacts back to the host.  The setting is only used in CP mode.  After the build is done equivalent Docker `cp` command is ran

    docker cp ${DAPPER_SOURCE}/${DAPPER_OUTPUT} .

If you don't want the `DAPPER_OUTPUT` to be relative to the `DAPPER_SOURCE` then set `DAPPER_OUTPUT` to a strings that starts with `/`. 


### DAPPER_DOCKER_SOCKET

Setting `DAPPER_DOCKER_SOCKET` will cause the Docker socket to be bind mounted into your build.  This is so that your build can use Docker without requiring Docker-in-Docker.  The equivalent parameter will be added to Docker.

   docker run -v /var/run/docker.sock:/var/run/docker.sock build-image

### DAPPER_RUN_ARGS

`DAPPER_RUN_ARGS` is used to add any parameters to the Docker `run` command for the build container.  For example you may want to set `--privileged` if you need to do advanced operations as root.

### DAPPER_ENV

`DAPPER_ENV` is a list of ENV variables that should be copied for the host context.  Setting `DAPPER_ENV=A B C` is the equivalent of adding to the Docker `run` command the following

    docker run -e A -e B -e C build-image

## License

Copyright (c) 2015-2018 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
