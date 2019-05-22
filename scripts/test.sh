#!/bin/bash
set -e

go get -v -t ./... && go test -v ./...

