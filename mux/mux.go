package mux

import (
	"fmt"
	"io"
	"net"

	"github.com/bearstech/ascetic-rpc/model"
	"github.com/bearstech/ascetic-rpc/protocol"
)

type server struct {
	socket   *net.UnixListener
	handlers map[string]func(req *model.Request) *model.Response
}

func NewServer(socket *net.UnixListener) *server {
	return &server{
		socket:   socket,
		handlers: make(map[string]func(req *model.Request) *model.Response),
	}
}

func (s *server) Route(name string, fun func(req *model.Request) *model.Response) {
	s.handlers[name] = fun
}

func (s *server) Listen() {
	for {
		conn, err := s.socket.AcceptUnix()
		if err != nil {
			panic(err)
		}
		go func() {
			for {
				err := s.Read(conn)
				if err != nil {
					// Do something
					s.socket.Close()
					fmt.Println(err.Error())
					return
				}
			}
		}()
	}
}

func (s *server) Read(wire io.ReadWriter) error {
	var req model.Request
	err := protocol.Read(wire, &req)
	if err != nil {
		return err
	}
	fmt.Println("header:", req)
	if req.Name == "" {
		return protocol.Write(wire, model.NewError(-1, "Empty method"))
	}
	h, ok := s.handlers[req.Name]
	if !ok {
		return protocol.Write(wire, model.NewError(-1, "Unknown method: "+req.Name))
	}
	return protocol.Write(wire, h(&req))
}
