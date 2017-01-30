#!/usr/bin/env bash
echo ${1}/chaincode
curl --request POST \
 --url ${1}/chaincode \
 --header 'content-type: application/json' \
 --data '{ "jsonrpc": "2.0", "method": "deploy", "params": { "type": 1, "chaincodeID": { "path": "https://github.com/dngroup/content-fabric/content-contract-cc" }, "ctorMsg": { "function": "init", "args": [ ] }, "secureContext": "user_type1_0" }, "id": 1}'
echo ""
