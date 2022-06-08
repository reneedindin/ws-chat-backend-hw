FROM golang:1.11.2-alpine

RUN apk add --no-cache git

WORKDIR /server

COPY . /server

WORKDIR /go/src
RUN go get -u github.com/satori/go.uuid golang.org/x/net/websocket

WORKDIR /server
RUN cd /server && go build

EXPOSE 12345
ENTRYPOINT ./server