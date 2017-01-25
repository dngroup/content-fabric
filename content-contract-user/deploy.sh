#!/usr/bin/env bash
curl --request POST \
 --url http://${1}:7050/chaincode \
 --header 'content-type: application/json' \
 --data '{ "jsonrpc": "2.0", "method": "deploy", "params": { "type": 1, "chaincodeID": { "path": "https://github.com/dngroup/content-fabric/content-contract-cc" }, "ctorMsg": { "function": "init", "args": [ ] } }, "id": 1}'
echo ""
