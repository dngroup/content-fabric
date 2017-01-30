FROM alpine
#TODO: NEED TO EDIT THIS
#FROM golang:1.6-alpine
MAINTAINER david bourasseau david.bourasseau@viotech.net
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ENV CHAINCODE=fb35eb0071ba47fdcf7ce76359e9f9e4c847a74d8bf0187c1a79ce3468322d378530ca9fc6aefb9c7b074c3d423dfb664142e0c0a3d905a173f4fce253f67372
ENV CP-ID=""
ENV EVENT_PORT=7053
ENV PEERADDR=172.17.0.1
ENV PEERADDR=localhost
ENV PERCENT=100
ENV REST_PORT=7050
ENV TLS=false
ENV USER=admin

COPY content-contract-cp content-contract-cp

#RUN ./b
CMD ./content-contract-cp -user=${USER} -tls=${TLS} -events-address=${PEERADDR}:${EVENT_PORT} -events-from-chaincode=${CHAINCODE} -CP-ID=${CP-ID} -percent=${PERCENT} -rest-address=${PEERADDR}:${REST_PORT}