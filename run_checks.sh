#!/bin/bash

# echo "Running go vet"
# go vet

# echo "Running gofmt check"
# gofmt -s -l .

# echo "Running go test"
# go test ./...

failed=0

echo "Running go vet"
go vet
if [ $? -ne 0 ]; then
    echo "Go Vet Failed."
    failed=1
fi

echo "Running gofmt check"
gofmt -s -l .

if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
    echo "Gofmt Check Failed."
    failed=1
fi

echo "Running go test"

go test ./...

if [ $? -ne 0 ]; then
    echo "Go Test Failed."
    failed=1
fi

if [ "$failed" -eq 1 ]; then
    exit 1
else
    exit 0
fi
