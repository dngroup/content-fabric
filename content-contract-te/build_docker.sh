#!/usr/bin/env bash
CGO_ENABLED=0 go build -o content-contract-te content-contract-te.go
chmod +x content-contract-te
docker build -t dngroup/content-contract-te .
