#!/usr/bin/env bash
url=https://9dbf187e9e4247e2a59e122e812101f6-vp$3.us.blockchain.ibm.com:5003
LINE=`cat users.csv|grep $1_$2`
LOGIN=${LINE%,*}
PASSWORD=${LINE#*,}

curl -v --request POST \
  --url $url/registrar \
  --header 'content-type: application/json' \
  --data "{ \"enrollId\": \"${LOGIN}\",  \"enrollSecret\": \"${PASSWORD}\"}"
