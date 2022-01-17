#!/bin/bash
set -e

go vet $(go list ./... | grep -v example)
go test $(go list ./... | grep -v example)