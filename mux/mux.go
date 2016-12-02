package mux

import (
	"errors"
	"fmt"
	"io"

	"github.com/bearstech/ascetic-rpc/model"
	"github.com/bearstech/ascetic-rpc/protocol"
	"github.com/golang/protobuf/proto"
)

type Handler interface {
	Handle(req_h *model.Request, req_b []byte) (model.Response, proto.Message)
}

type server struct {
	wire     io.ReadWriter
	handlers map[string]Handler
}

func NewServer(wire io.ReadWriter) *server {
	return &server{
		wire:     wire,
		handlers: make(map[string]Handler),
	}
}

func (s *server) Route(name string, handler Handler) {
	s.handlers[name] = handler
}

func (s *server) Read() error {
	var req_h model.Request
	err := protocol.Read(s.wire, &req_h)
	if err != nil {
		return err
	}
	fmt.Println(req_h)
	h, ok := s.handlers[req_h.Name]
	if !ok {
		return errors.New("Not found")
	}
	req_b, err := protocol.ReadBytes(s.wire)
	if err != nil {
		return err
	}
	res_h, res_b := h.Handle(&req_h, req_b)
	err = protocol.Write(s.wire, &res_h)
	if err != nil {
		return err
	}
	err = protocol.Write(s.wire, res_b)
	if err != nil {
		return err
	}
	return nil
}
