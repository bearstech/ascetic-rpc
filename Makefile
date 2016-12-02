export GOPATH:=$(shell pwd)/gopath

lib: gopath/src/github.com/bearstech/ascetic-rpc gopath/src/github.com/golang/protobuf/proto

clean:
	rm -rf gopath

protoc:
	protoc --go_out=. model/*.proto

test:
	go test -v github.com/bearstech/ascetic-rpc/mux


# Kitchen sinks

gopath/src/github.com/bearstech/ascetic-rpc:
	mkdir -p gopath/src/github.com/bearstech/
	ln -s `pwd` gopath/src/github.com/bearstech/ascetic-rpc

gopath/src/github.com/golang/protobuf/proto:
	go get github.com/golang/protobuf/proto
