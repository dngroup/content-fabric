FROM alpine
#TODO: NEED TO EDIT THIS
#FROM golang:1.6-alpine
MAINTAINER david bourasseau david.bourasseau@viotech.net
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ENV CHAINCODE=fb35eb0071ba47fdcf7ce76359e9f9e4c847a74d8bf0187c1a79ce3468322d378530ca9fc6aefb9c7b074c3d423dfb664142e0c0a3d905a173f4fce253f67372
ENV CONTENTID=BBB.mp4
ENV EVENT_PORT=7053
ENV PEERADDR=172.17.0.1
ENV REST_PORT=7050
ENV TIMEMAX=10
ENV TLS=false
ENV USER=admin
ENV USERID=LALA

COPY content-contract-user content-contract-user

#RUN ./b
CMD ./content-contract-user -user=${USER} -tls=${TLS} -userId=${USERID} -contentId=${CONTENTID} -chaincodeid=${CHAINCODE} -time-max=${TIMEMAX} -rest-address=${PEERADDR}:${REST_PORT} -events-address=${PEERADDR}:${EVENT_PORT}
