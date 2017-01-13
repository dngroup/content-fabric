#!/usr/bin/env bash
#load the chaine code
# Remove old docker
docker stop ab bc cd de da
docker rm ab bc cd de da

ADDR=172.17.0.1
echo "Load  the chaincode and get the id of this"
# you have to edit the path if you fork this repository
CHAINCODEID=$(curl --request POST \
  --url http://localhost:7050/chaincode \
  --header 'content-type: application/json' \
  --data "{\"jsonrpc\":\"2.0\",\"method\":\"deploy\",\"params\": {\"type\": 1,\"chaincodeID\": {\"path\":\"https://github.com/dngroup/content-fabric/chainecode\"  },\"ctorMsg\": {\"function\":\"init\",\"args\": [ ]  },\"secureContext\":\"test_user0\"  },\"id\": 1 }" \
  | json_pp | grep -o -P '(?<="message" : ").*(?=")' )
  # sed -e 's/"message" : "\(.*\)"/\1/')



#.user/user -events-address=localhost:7053 -events-from-chaincode=${id} -value-to-analyse=a -value-to-change=b -rest-address=localhost:7050
echo "CHAINCODEID=${CHAINCODEID}"
sleep 3

echo "Start docker listener"
echo "docker run -d --name=ab -e PEERADDR=${ADDR} -e A=a -e B=b -e CHAINCODE=${CHAINCODEID} dngroup/user"
docker run -d --name=ab -e PEERADDR=${ADDR} -e A=a -e B=b -e CHAINCODE=${CHAINCODEID} dngroup/user
docker run -d --name=bc -e PEERADDR=${ADDR} -e A=b -e B=c -e CHAINCODE=${CHAINCODEID} dngroup/user
docker run -d --name=cd -e PEERADDR=${ADDR} -e A=c -e B=d -e CHAINCODE=${CHAINCODEID} dngroup/user
docker run -d --name=da -e PEERADDR=${ADDR} -e A=d -e B=a -e CHAINCODE=${CHAINCODEID} dngroup/user





echo "wait 3 second before set a to 1"
sleep 3

curl --request POST \
 --url http://${ADDR}:7050/chaincode \
 --header 'content-type: application/json' \
 --data "{\"jsonrpc\":\"2.0\",\"method\":\"invoke\",\"params\": {\"type\": 1,\"chaincodeID\": {\"name\":\"${CHAINCODEID}\" },\"ctorMsg\": {\"function\":\"write\",\"args\": [\"a\",\"1\" ] },\"secureContext\":\"test_user0\" },\"id\": 3}"

echo "copy this to show the value of a"
echo "while true; do  curl -sb -X POST -H \"Content-Type: application/json\" -d '{ \"jsonrpc\": \"2.0\", \"method\": \"query\", \"params\": { \"type\": 1, \"chaincodeID\":{ \"name\":\"${CHAINCODEID}\" }, \"ctorMsg\": { \"function\":\"read\", \"args\":[\"a\"] }}, \"id\": 2}' \"http://localhost:7050/chaincode\" | json_pp | grep message ;  done"
