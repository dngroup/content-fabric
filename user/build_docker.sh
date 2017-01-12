#!/usr/bin/env bash

CGO_ENABLED=0 go build -o user user.go
chmod +x user
docker build -t dngroup/user .
