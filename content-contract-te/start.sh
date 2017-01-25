#!/usr/bin/env bash
for ((number=1; number<=$1;number++))
do
    docker rm -f te${number}
done
if [ "$#" -ne 4 ]; then
   echo "./start NUMBEROFCP PERCENTOFCHANCETOHAVECONTENT PERCENTOFCHANCEFORPRICELOWER CONTRACTID PeerIP"

   exit -1
fi

for ((number=1; number<=$1;number++))
do
#    docker rm -f te${number}
    docker run --name te${number} -d -e PEERADDR=$5 -e CP-ID=cp${number} -e CHAINCODE=$4 -e PERCENT=$2 -e PERCENTPRICE=$3 dngroup/content-contract-te
done
