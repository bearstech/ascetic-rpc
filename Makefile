export GOPATH:=$(shell pwd)/gopath

.PHONY: plugin

lib: gopath/src/github.com/bearstech/ascetic-rpc gopath/src/github.com/golang/protobuf/proto

clean:
	rm -rf gopath

protoc:
	protoc --go_out=. message/message.proto

test:
	go test -v -cover github.com/bearstech/ascetic-rpc/server
	go test -v -cover github.com/bearstech/ascetic-rpc/protocol
	go test -v -cover github.com/bearstech/ascetic-rpc/client
	#go test -v github.com/bearstech/ascetic-rpc/register

plugin:
	go build -o bin/protoc-gen-ascetic github.com/bearstech/ascetic-rpc/protoc-gen-ascetic
	#protoc --plugin=bin/protoc-gen-ascetic --go_out=plugins=ascetic:. model/test.proto
	#protoc --plugin=bin/protoc-gen-ascetic --go_out=plugins=grpc:. model/test.proto
	PATH=$(PATH):bin protoc --ascetic_out=plugins=ascetic:. model/test.proto

coverage:
	go test -v -coverprofile=coverage.out github.com/bearstech/ascetic-rpc/server
	go tool cover -html=coverage.out
	go test -v -coverprofile=coverage.out github.com/bearstech/ascetic-rpc/protocol
	go tool cover -html=coverage.out
	go test -v -coverprofile=coverage.out github.com/bearstech/ascetic-rpc/client
	go tool cover -html=coverage.out
	#go test -v github.com/bearstech/ascetic-rpc/register
	#go tool cover -html=coverage.out

# Kitchen sinks

bin:
	mkdir bin

gopath/src/github.com/bearstech/ascetic-rpc:
	mkdir -p gopath/src/github.com/bearstech/
	ln -s `pwd` gopath/src/github.com/bearstech/ascetic-rpc

gopath/src/github.com/golang/protobuf/proto:
	go get github.com/golang/protobuf/proto
