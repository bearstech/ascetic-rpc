package client

import (
	"errors"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/bearstech/ascetic-rpc/model"
	"github.com/bearstech/ascetic-rpc/mux"
	"github.com/golang/protobuf/proto"
)

func hello(req_h *model.Request, req_b []byte) (model.Response, proto.Message) {
	var hello model.Hello
	err := proto.Unmarshal(req_b, &hello)
	if err != nil {
		panic(err)
	}
	world := model.World{
		Message: fmt.Sprintf("Hello %s🐈", hello.Name),
	}
	return model.Response{Code: 1}, &world
}

func TestClientHello(t *testing.T) {
	socketPath := "/tmp/test_client.sock"
	os.Remove(socketPath)

	l, err := net.ListenUnix("unix", &net.UnixAddr{
		Name: socketPath,
		Net:  "unix",
	})
	if err != nil {
		t.Error(err)
	}

	s := mux.NewServer(l)
	s.Route("hello", hello)
	go s.Listen()

	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{
		Name: socketPath,
		Net:  "unix"})
	if err != nil {
		t.Error(err)
	}
	c := New(conn)

	hello := model.Hello{Name: "Alice"}
	var world model.World

	err = c.Do("hello", &hello, &world)
	if err != nil {
		t.Error(err)
	}

	if world.Message != "Hello Alice🐈" {
		t.Error(errors.New("Bad message: " + world.Message))
	}
}
