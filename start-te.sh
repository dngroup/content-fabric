#!/usr/bin/env bash
for ((number=1; number<=$1;number++))
do
    docker rm -f te1-${number} te2-${number} te3-${number} te4-${number}
done
if [ "$#" -ne 4 ]; then
   echo "./start NUMBEROFCP PERCENTOFCHANCETOHAVECONTENT PERCENTOFCHANCEFORPRICELOWER CONTRACTID"

   exit -1
fi

for ((number=1; number<=$1;number++))
do
#    docker rm -f te${number}
    docker run --name te1-${number} -d -e PEERADDR=34.249.48.106 -e CP-ID=cp${number} -e CHAINCODE=$4 -e PERCENT=$2 -e PERCENTPRICE=$3 dngroup/content-contract-te
done
for ((number=1; number<=$1;number++))
do
#    docker rm -f te${number}
    docker run --name te2-${number} -d -e PEERADDR=34.249.177.11 -e CP-ID=cp${number} -e CHAINCODE=$4 -e PERCENT=$2 -e PERCENTPRICE=$3 dngroup/content-contract-te
done
for ((number=1; number<=$1;number++))
do
#    docker rm -f te${number}
    docker run --name te3-${number} -d -e PEERADDR=54.211.4.170 -e CP-ID=cp${number} -e CHAINCODE=$4 -e PERCENT=$2 -e PERCENTPRICE=$3 dngroup/content-contract-te
done
for ((number=1; number<=$1;number++))
do
#    docker rm -f te${number}
    docker run --name te4-${number} -d -e PEERADDR=13.54.121.139 -e CP-ID=cp${number} -e CHAINCODE=$4 -e PERCENT=$2 -e PERCENTPRICE=$3 dngroup/content-contract-te
done