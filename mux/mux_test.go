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

func ping(*model.Request, []byte) (model.Response, proto.Message) {
	return model.Response{Code: 1}, nil
}

func TestPing(t *testing.T) {
	w := wire.New()
	s := NewServer(nil)
	s.Route("ping", ping)

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

func hello(req_h *model.Request, req_b []byte) (model.Response, proto.Message) {
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
	s.Route("hello", hello)

	var err error
	req := model.Request{
		Name: "plop",
	}
	err = protocol.WriteHeaderAndBody(w.ClientToServer(), &req, nil)
	err = s.Read(w.ServerToClient())
	if err != nil {
		t.Error(err)
	}
	var resp model.Response
	err = protocol.Read(w.ClientToServer(), &resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Code != -1 {
		t.Error(errors.New("It should be unknown"))
	}
	fmt.Println(w.Len())
	err = protocol.Read(w.ClientToServer(), nil)
	if err != nil {
		t.Error(err)
	}
	lin, lout := w.Len()
	fmt.Println("Wire len: ", lin, lout)
	if lin != 0 {
		t.Error(errors.New("bad size"))
	}
	if lout != 0 {
		t.Error(errors.New("bad size"))
	}

	req2 := model.Request{
		Name: "hello",
	}
	hello := model.Hello{Name: "Bob"}
	fmt.Println("deuz: ", req2)
	err = protocol.WriteHeaderAndBody(w.ClientToServer(), &req2, &hello)
	if err != nil {
		t.Error(err)
	}
	err = s.Read(w.ServerToClient())
	if err != nil {
		t.Error(err)
	}

	resp.Reset()
	err = protocol.Read(w.ClientToServer(), &resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Code < 0 {
		t.Error(errors.New("It's an error: " + resp.Message))
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
