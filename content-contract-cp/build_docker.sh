#!/usr/bin/env bash
CGO_ENABLED=0 go build -o content-contract-cp content-contract-cp.go
chmod +x content-contract-cp
docker build -t dngroup/content-contract-cp .
