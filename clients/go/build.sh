#!/bin/bash

GOOS=linux GOARCH=arm64 go build -o releases/Rex-linux-arm64
GOOS=linux GOARCH=amd64 go build -o releases/Rex-linux-amd64
GOOS=windows GOARCH=amd64 go build -o releases/Rex-windows-amd64
GOOS=darwin GOARCH=arm64 go build -o releases/Rex-darwin-arm64