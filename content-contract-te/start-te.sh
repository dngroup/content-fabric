#!/usr/bin/env bash
URL0=172.22.0.
URL1=
for ((number=1; number<=$1;number++))
do
    docker rm -f te1-${number} te2-${number} te3-${number} te0-${number}
done
if [ "$#" -ne 4 ]; then
   echo "./start NUMBEROFCP PERCENTOFCHANCETOHAVECONTENT PERCENTOFCHANCEFORPRICELOWER CONTRACTID"

   exit -1
fi

for ((number=1; number<=$1;number++))
do
#    docker rm -f te${number}
    docker run --name te0-${number} -d \
        -e USER=user_type1_0 \
        -e PEERADDR=${URL0}4${URL1} \
        -e TE-ID=cp${number} \
        -e CHAINCODE=$4 \
        -e PERCENT=$2 \
        -e REST_PORT=7050 \
        -e EVENT_PORT=7053 \
        -e PERCENTPRICE=$3 \
        --net fabric_default \
        dngroup/content-contract-te
done
for ((number=1; number<=$1;number++))
do
#    docker rm -f te${number}
    docker run --name te1-${number} -d \
        -e USER=user_type1_1 \
        -e PEERADDR=${URL0}3${URL1} \
        -e TE-ID=cp${number} \
        -e CHAINCODE=$4 \
        -e PERCENT=$2 \
        -e REST_PORT=7050 \
        -e EVENT_PORT=7053 \
        -e PERCENTPRICE=$3 \
        --net fabric_default \
        dngroup/content-contract-te
done
for ((number=1; number<=$1;number++))
do
#    docker rm -f te${number}
    docker run --name te2-${number} -d \
        -e USER=user_type1_2 \
        -e PEERADDR=${URL0}5${URL1} \
        -e TE-ID=cp${number} \
        -e CHAINCODE=$4 \
        -e PERCENT=$2 \
        -e REST_PORT=7050 \
        -e EVENT_PORT=7053 \
        -e PERCENTPRICE=$3 \
        --net fabric_default \
        dngroup/content-contract-te
done
for ((number=1; number<=$1;number++))
do
#    docker rm -f te${number}
    docker run --name te3-${number} -d \
        -e USER=user_type1_3 \
        -e PEERADDR=${URL0}2${URL1} \
        -e TE-ID=cp${number} \
        -e CHAINCODE=$4 \
        -e PERCENT=$2 \
        -e REST_PORT=7050 \
        -e EVENT_PORT=7053 \
        -e PERCENTPRICE=$3 \
        --net fabric_default \
        dngroup/content-contract-te
done