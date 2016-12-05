package mux

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bearstech/ascetic-rpc/model"
	"github.com/bearstech/ascetic-rpc/protocol"
	"github.com/bearstech/ascetic-rpc/wire"
)

func ping(*model.Request) *model.Response {
	return &model.Response{Code: 1}
}

func TestPing(t *testing.T) {
	w := wire.New()
	s := NewServer(nil)
	s.Route("ping", ping)

	req := model.Request{
		Name: "ping",
	}
	err := protocol.Write(w.ClientToServer(), &req)
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

func hello(req *model.Request) *model.Response {
	var hello model.Hello
	err := req.GetBody(&hello)
	if err != nil {
		panic(err)
	}
	world := model.World{
		Message: fmt.Sprintf("Hello %s♥️", hello.Name),
	}
	res, err := model.NewOK(1, &world)
	if err != nil {
		return model.NewError(-2, err.Error())
	}
	return res
}

func TestHello(t *testing.T) {
	w := wire.New()
	s := NewServer(nil)
	s.Route("hello", hello)

	var err error
	req := model.Request{
		Name: "plop",
	}
	err = protocol.Write(w.ClientToServer(), &req)
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

	req2, err := model.NewRequest("hello", &model.Hello{Name: "Bob"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("deuz: ", req2)
	err = protocol.Write(w.ClientToServer(), req2)
	if err != nil {
		t.Error(err)
	}

	err = s.Read(w.ServerToClient())
	if err != nil {
		t.Error(err)
	}
	err = protocol.Read(w.ClientToServer(), &resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Code < 0 {
		t.Error(errors.New("It's an error: " + resp.GetError().Message))
	}

	var world model.World
	err = resp.ReadOK(&world)
	if err != nil {
		t.Error(err)
	}
	if world.Message != "Hello Bob♥️" {
		t.Error(errors.New("Bad message: " + world.Message))
	}
}
