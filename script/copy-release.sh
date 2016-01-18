#!/bin/bash
VERSION=$(./bin/dapper-Linux-x86_64 -v | awk '{print $3}')
gsutil -m cp -r bin/  gs://releases.rancher.com/dapper/${VERSION}
