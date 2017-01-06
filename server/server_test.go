package server

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bearstech/ascetic-rpc/client"
	"github.com/bearstech/ascetic-rpc/message"
)

func hello(req *message.Request) (*message.Response, error) {
	var hello message.Hello
	err := req.GetBody(&hello)
	if err != nil {
		panic(err)
	}
	world := message.World{
		Message: fmt.Sprintf("Hello %s♥️", hello.Name),
	}
	res, err := message.NewOKResponse(&world)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func TestHelloServer(t *testing.T) {
	socketPath := "/tmp/test.sock"
	s, err := NewServerUnix(socketPath)
	if err != nil {
		t.Error(err)
	}
	s.Register("hello", hello)
	serverStopped := false
	go func() {
		s.Serve()
		serverStopped = true
	}()

	c, err := client.NewClientUnix(socketPath)
	if err != nil {
		t.Error(err)
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
	if world.Message != "Hello Alice♥️" {
		t.Error(errors.New("Bad message: " + world.Message))
	}
	s.Stop()
	if !serverStopped {
		t.Error(errors.New("Bad stop"))
	}
	if s.IsRunning() {
		t.Error(errors.New("Ghost running"))
	}
}

func dontpanic(req *message.Request) (*message.Response, error) {
	panic(errors.New("oups"))
}

func TestPanic(t *testing.T) {

	socketPath := "/tmp/test.sock"
	s, err := NewServerUnix(socketPath)
	if err != nil {
		t.Error(err)
	}
	defer s.Stop()
	s.Register("panic", dontpanic)
	go s.Serve()

	c, err := client.NewClientUnix(socketPath)
	if err != nil {
		t.Error(err)
	}
	err = c.Do("panic", nil, nil)
	if err == nil {
		t.Error(errors.New("Should not be nil"))
	}
	if err.Error() != "oups" {
		t.Error(errors.New("Should be oups"))
	}
}
