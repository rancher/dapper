#!/bin/bash
VERSION=$(./bin/dapper-Linux-x86_64 -v | awk '{print $3}')
gsutil -m cp -r -p winged-math-749 bin/*  gs://releases.rancher.com/dapper/${VERSION}
