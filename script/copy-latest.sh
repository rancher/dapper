#!/bin/bash
gsutil -m rsync -r bin/   gs://releases.rancher.com/dapper/latest
