export GOPATH:=$(shell pwd)/gopath

lib: gopath/src/github.com/bearstech/ascetic-rpc gopath/src/github.com/golang/protobuf/proto

clean:
	rm -rf gopath

protoc:
	protoc --go_out=. model/*.proto

test:
	go test -v -cover github.com/bearstech/ascetic-rpc/server
	go test -v -cover github.com/bearstech/ascetic-rpc/protocol
	go test -v -cover github.com/bearstech/ascetic-rpc/client
	#go test -v github.com/bearstech/ascetic-rpc/register

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

gopath/src/github.com/bearstech/ascetic-rpc:
	mkdir -p gopath/src/github.com/bearstech/
	ln -s `pwd` gopath/src/github.com/bearstech/ascetic-rpc

gopath/src/github.com/golang/protobuf/proto:
	go get github.com/golang/protobuf/proto
