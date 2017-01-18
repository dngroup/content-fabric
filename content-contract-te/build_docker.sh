#!/usr/bin/env bash
#TODO: NEED TO EDIT THIS
CGO_ENABLED=0 go build -o user user.go
chmod +x user
docker build -t dngroup/user .
