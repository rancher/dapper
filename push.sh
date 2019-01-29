#!/bin/bash
docker push ibuildthecloud/dapper:test-amd64
docker push ibuildthecloud/dapper:test-arm
docker push ibuildthecloud/dapper:test-arm64
manifest-tool push from-args --platforms linux/arm,linux/arm64,linux/amd64 --template ibuildthecloud/dapper:test-ARCH --target ibuildthecloud/dapper:test
