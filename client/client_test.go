package client

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/bearstech/ascetic-rpc/model"
	"github.com/bearstech/ascetic-rpc/mux"
)

func hello(req *model.Request) (*model.Response, error) {
	var hello model.Hello
	err := req.GetBody(&hello)
	if err != nil {
		panic(err)
	}
	world := model.World{
		Message: fmt.Sprintf("Hello %s🐈", hello.Name),
	}
	resp, err := model.NewOKResponse(1, &world)
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

	s := mux.NewServer(l)
	s.Register("hello", hello)
	go s.Listen()

	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{
		Name: socketPath,
		Net:  "unix"})
	if err != nil {
		t.Error(err)
	}
	c := New(conn)

	// Unknown function
	err = c.Do("oups", nil, nil)
	if !strings.HasPrefix(err.Error(), "Unknown method") {
		t.Error(errors.New("Wrong error: " + err.Error()))
	}

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
