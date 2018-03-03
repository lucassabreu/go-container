#!/bin/bash

which goverage &> /dev/null || go get -u github.com/haya14busa/goverage
goverage -v -race -coverprofile=coverage.out ./...