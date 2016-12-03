package mux

import (
	"fmt"
	"io"
	"net"

	"github.com/bearstech/ascetic-rpc/model"
	"github.com/bearstech/ascetic-rpc/protocol"
	"github.com/golang/protobuf/proto"
)

type server struct {
	socket   *net.UnixListener
	handlers map[string]func(req_h *model.Request, req_b []byte) (model.Response, proto.Message)
}

func NewServer(socket *net.UnixListener) *server {
	return &server{
		socket:   socket,
		handlers: make(map[string]func(req_h *model.Request, req_b []byte) (model.Response, proto.Message)),
	}
}

func (s *server) Route(name string, fun func(req_h *model.Request, req_b []byte) (model.Response, proto.Message)) {
	s.handlers[name] = fun
}

func (s *server) Listen() {
	for {
		conn, err := s.socket.AcceptUnix()
		if err != nil {
			panic(err)
		}
		err = s.Read(conn)
		if err != nil {
			// Do something
		}
	}
}

func (s *server) Read(wire io.ReadWriter) error {
	var req_h model.Request
	err := protocol.Read(wire, &req_h)
	if err != nil {
		s.socket.Close()
		return err
	}
	fmt.Println("header:", req_h)
	if req_h.Name == "" {
		err2 := protocol.Read(wire, nil)
		if err2 != nil {
			// oups
			fmt.Println(err2)
		}
		res_h := model.Response{
			Code:    -1,
			Message: "Empty method"}
		err2 = protocol.WriteHeaderAndBody(wire, &res_h, nil)
		if err2 != nil {
			// oups
			fmt.Println(err2)
		}
		return nil
	}
	h, ok := s.handlers[req_h.Name]
	if !ok {
		err2 := protocol.Read(wire, nil) // Drain body
		res_h := model.Response{
			Code:    -1,
			Message: "Unknown method: " + req_h.Name}
		err2 = protocol.WriteHeaderAndBody(wire, &res_h, nil)
		if err2 != nil {
			// oups
			fmt.Println(err2)
		}
		return nil
	}
	req_b, err := protocol.ReadBytes(wire)
	if err != nil {
		return err
	}
	res_h, res_b := h(&req_h, req_b)
	err = protocol.WriteHeaderAndBody(wire, &res_h, res_b)
	if err != nil {
		return err
	}
	return nil
}
