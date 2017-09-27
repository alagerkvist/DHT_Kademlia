FROM golang:1.7-alpine

ADD . /home
        
WORKDIR /home

RUN apk update
RUN apk add screen

RUN \
       apk add --no-cache bash git openssh && \
       go get -u github.com/minio/minio-go 
       
RUN ls