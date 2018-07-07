FROM golang:1.10

WORKDIR /go/src/app

RUN go get -u gopkg.in/alecthomas/gometalinter.v2 && \
	go get -u github.com/golang/dep/cmd/dep && \
	go get -u golang.org/x/tools/cmd/goimports && \
    go get -u github.com/mattn/goveralls

RUN apt-get update && \
    apt-get install -y libpcre++-dev
