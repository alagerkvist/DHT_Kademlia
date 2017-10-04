
FROM ubuntu:14.04
RUN apt-get update
RUN apt-get -y upgrade
RUN apt-get -y install curl
RUN apt-get -y install git

RUN curl -O https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz
RUN mv go1.9.linux-amd64.tar.gz /usr/local/
RUN tar -xf /usr/local/go1.9.linux-amd64.tar.gz -C /usr/local/

ENV PATH $PATH:/usr/local/go/bin
ENV SRC_DIR=/go/local/github.com/ramirobg94/DHT_Kademlia        
ADD . $SRC_DIR    

RUN mv /go/local/github.com/ramirobg94/DHT_Kademlia/* .

RUN cd $SRC_DIR;

RUN go env
RUN apt-get install -y protobuf-compiler
RUN apt-get install -y golang-goprotobuf-dev
RUN go get -u -v github.com/golang/protobuf/proto
RUN go get -u -v github.com/golang/protobuf/protoc-gen-go

RUN ls
