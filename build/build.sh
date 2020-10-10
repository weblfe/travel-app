#!/usr/bin/env bash

GOOS=windows GOARCH=amd64 go build  -ldflags '-d -w -s'  ./...

