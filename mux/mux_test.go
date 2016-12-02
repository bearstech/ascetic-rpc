package mux

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bearstech/ascetic-rpc/model"
	"github.com/bearstech/ascetic-rpc/protocol"
	"github.com/bearstech/ascetic-rpc/wire"
	"github.com/golang/protobuf/proto"
)

type ping struct{}

func (p ping) Handle(*model.Request, []byte) (model.Response, proto.Message) {
	return model.Response{Code: 1}, nil
}

func TestPing(t *testing.T) {
	w := wire.New()
	s := NewServer(nil)
	s.Route("ping", ping{})

	req := model.Request{
		Name: "ping",
	}
	err := protocol.WriteHeaderAndBody(w.ClientToServer(), &req, nil)
	if err != nil {
		t.Error(err)
	}
	err = s.Read(w.ServerToClient())
	if err != nil {
		t.Error(err)
	}

	var resp model.Response
	err = protocol.Read(w.ClientToServer(), &resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Code != 1 {
		t.Fatal()
	}
}

type hello struct{}

func (h hello) Handle(req_h *model.Request, req_b []byte) (model.Response, proto.Message) {
	var hello model.Hello
	err := proto.Unmarshal(req_b, &hello)
	if err != nil {
		panic(err)
	}
	world := model.World{
		Message: fmt.Sprintf("Hello %s♥️", hello.Name),
	}
	return model.Response{Code: 1}, &world
}

func TestHello(t *testing.T) {
	w := wire.New()
	s := NewServer(nil)
	s.Route("hello", hello{})
	req := model.Request{
		Name: "hello",
	}
	hello := model.Hello{Name: "Bob"}
	err := protocol.WriteHeaderAndBody(w.ClientToServer(), &req, &hello)
	if err != nil {
		t.Error(err)
	}
	err = s.Read(w.ServerToClient())
	if err != nil {
		t.Error(err)
	}

	var resp model.Response
	err = protocol.Read(w.ClientToServer(), &resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Code < 0 {
		t.Error(errors.New("It's an error"))
	}

	var world model.World
	err = protocol.Read(w.ClientToServer(), &world)
	if err != nil {
		t.Error(err)
	}
	if world.Message != "Hello Bob♥️" {
		t.Error(errors.New("Bad message: " + world.Message))
	}
}
