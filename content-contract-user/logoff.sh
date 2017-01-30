#!/usr/bin/env bash
url=https://9dbf187e9e4247e2a59e122e812101f6-vp$2.us.blockchain.ibm.com:5003
l=$(sed -n "$(($1+3)) p" user.csv)
echo $l
login=$(echo ${l} | cut -b -12 )
password=$(echo ${l} | cut -b 14- )
curl -v --request DELETE \
  --url $url/registrar/${login}
echo login: \"${login}\"
