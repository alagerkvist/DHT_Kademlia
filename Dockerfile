# iron/go:dev is the alpine image with the go tools added
FROM nanoservice/protobuf:3.0-alpha

RUN apk add --update go
RUN apk add --update git

RUN apk add --update go
RUN apk add --update git

RUN mkdir -p /go/src

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

RUN go get -u github.com/golang/protobuf/proto github.com/golang/protobuf/protoc-gen-go

RUN apk del git

ENV SRC_DIR=/go/src/github.com/ramirobg94/DHT_Kademlia
ADD . $SRC_DIR

# Build it:
RUN cd $SRC_DIR; go build; go run main.go