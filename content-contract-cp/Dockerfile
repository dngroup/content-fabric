FROM alpine
#TODO: NEED TO EDIT THIS
#FROM golang:1.6-alpine
MAINTAINER david bourasseau david.bourasseau@viotech.net
ENV PEERADDR=localhost
ENV A=a
ENV B=b
ENV REST_PORT=7050
ENV EVENT_PORT=7053
ENV CHAINCODE=07471aa6f663b437b318792909f93b804498e60329b89d5668631827b6d99d11b2344e606805ba80ef7024bc271fb319d69b2526db40e944829f4a4a579bded7
#COPY user.go user.go
COPY user user

#RUN ./b
CMD ./user -events-address=${PEERADDR}:${EVENT_PORT} -events-from-chaincode=${CHAINCODE} -value-to-analyse=${A} -value-to-change=${B} -rest-address=${PEERADDR}:${REST_PORT}