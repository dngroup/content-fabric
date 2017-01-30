#!/usr/bin/env bash
url=https://d0ffb689045e4dfeb25fd8df4bafca84-vp$2.us.blockchain.ibm.com:5002
l=$(sed -n "$(($1+3)) p" user.csv)
echo $l
login=$(echo ${l} | cut -b -12 )
password=$(echo ${l} | cut -b 14- )
curl -v --request DELETE \
  --url $url/registrar/${login}
echo login: \"${login}\"
