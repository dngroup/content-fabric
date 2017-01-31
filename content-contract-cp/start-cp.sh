#!/usr/bin/env bash
URL0=172.22.0.
URL1=
for ((number=1; number<=$1;number++))
do
    docker rm -f cp1-${number} cp2-${number} cp3-${number} cp0-${number}
done
if [ "$#" -ne 3 ]; then
   echo "./start NUMBEROFCP PERCENTOFCHANCETOHAVECONTENT CONTRACTID"

   exit -1
fi

for ((number=1; number<=$1;number++))
do
#    docker rm -f cp${number}
    docker run --name cp0-${number} -d \
        -e USER=user_type1_0 \
        -e PEERADDR=${URL0}4${URL1} \
        -e CP-ID=cp${number} \
        -e CHAINCODE=$3 \
        -e PERCENT=$2 \
        -e REST_PORT=7050 \
        -e EVENT_PORT=7053 \
        --net fabric_default \
        dngroup/content-contract-cp
done

for ((number=1; number<=$1;number++))
do
#    docker rm -f cp${number}
    docker run --name cp1-${number} -d \
        -e USER=user_type1_1 \
        -e PEERADDR=${URL0}3${URL1} \
        -e CP-ID=cp${number} \
        -e CHAINCODE=$3 \
        -e PERCENT=$2 \
        -e REST_PORT=7050 \
        -e EVENT_PORT=7053 \
        --net fabric_default \
        dngroup/content-contract-cp
done
for ((number=1; number<=$1;number++))
do
#    docker rm -f cp${number}
    docker run --name cp2-${number} -d \
        -e USER=user_type1_2 \
        -e PEERADDR=${URL0}5${URL1} \
        -e CP-ID=cp${number} \
        -e CHAINCODE=$3 \
        -e PERCENT=$2 \
        -e REST_PORT=7050 \
        -e EVENT_PORT=7053 \
        --net fabric_default \
        dngroup/content-contract-cp
done
for ((number=1; number<=$1;number++))
do
#    docker rm -f cp${number}
    docker run --name cp3-${number} -d \
        -e USER=user_type1_3 \
        -e PEERADDR=${URL0}2${URL1} \
        -e CP-ID=cp${number} \
        -e CHAINCODE=$3 \
        -e PERCENT=$2 \
        -e REST_PORT=7050 \
        -e EVENT_PORT=7053 \
        --net fabric_default \
        dngroup/content-contract-cp
done