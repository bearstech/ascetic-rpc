package client

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/bearstech/ascetic-rpc/message"
	"github.com/bearstech/ascetic-rpc/server"
)

func hello(req *message.Request) (*message.Response, error) {
	var hello message.Hello
	err := req.GetBody(&hello)
	if err != nil {
		panic(err)
	}
	world := message.World{
		Message: fmt.Sprintf("Hello %s🐈", hello.Name),
	}
	resp, err := message.NewOKResponse(&world)
	if err != nil {
		return nil, err
	}
	fmt.Println("Response: ", resp)
	return resp, nil
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

	s := server.NewServer(l)
	s.Register("hello", hello)
	go s.Serve()

	c, err := NewClientUnix(socketPath)
	if err != nil {
		t.Error(err)
	}

	// Unknown function
	err = c.Do("oups", nil, nil)
	if !strings.HasPrefix(err.Error(), "Unknown method") {
		t.Error(errors.New("Wrong error: " + err.Error()))
	}

	hello := message.Hello{Name: "Alice"}
	var world message.World

	err = c.Do("hello", &hello, &world)
	if err != nil {
		t.Error(err)
	}
	err = c.Close()
	if err != nil {
		t.Error(err)
	}
	if world.Message != "Hello Alice🐈" {
		t.Error(errors.New("Bad message: " + world.Message))
	}

	s.Stop()

}
