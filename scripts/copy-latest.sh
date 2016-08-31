#!/bin/bash
gsutil -m cp -r dist/artifacts/* gs://releases.rancher.com/dapper/latest
