#!/usr/bin/env bash
CGO_ENABLED=0 go build -o content-contract-user content-contract-user.go
chmod +x content-contract-user
docker build -t dngroup/content-contract-user .
