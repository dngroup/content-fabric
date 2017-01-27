#!/usr/bin/env bash
for ((number=1; number<=$1;number++))
do
    docker rm -f cp${number}
done
if [ "$#" -ne 3 ]; then
    echo "./start NUMBEROFCP PERCENTOFCHANCETOHAVECONTENT CONTRACTID"

    exit -1
fi

for ((number=1; number<=$1;number++))
do
#    docker rm -f cp${number}
    docker run --name cp${number} -d -e PEERADDR=d0ffb689045e4dfeb25fd8df4bafca84-vp0.us.blockchain.ibm.com -e EVENT_PORT=31002 -e CP-ID=cp${number} -e CHAINCODE=$3 -e PERCENT=$2 dngroup/content-contract-cp
done