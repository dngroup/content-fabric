FROM alpine
#TODO: NEED TO EDIT THIS
#FROM golang:1.6-alpine
MAINTAINER david bourasseau david.bourasseau@viotech.net
ENV REST_PORT=7050
ENV EVENT_PORT=7053
ENV PEERADDR=172.17.0.1
ENV CHAINCODE=fb35eb0071ba47fdcf7ce76359e9f9e4c847a74d8bf0187c1a79ce3468322d378530ca9fc6aefb9c7b074c3d423dfb664142e0c0a3d905a173f4fce253f67372
ENV USERID=LALA
ENV CONTENTID=BBB.mp4
ENV TIMEMAX=10
ENV TLS=false
COPY content-contract-user content-contract-user

#RUN ./b
CMD ./content-contract-user -tls=${TLS} -userId=${USERID} -contentId=${CONTENTID} -chaincodeid=${CHAINCODE} -time-max=${TIMEMAX} -rest-address=${PEERADDR}:${REST_PORT} -events-address=${PEERADDR}:${EVENT_PORT}
