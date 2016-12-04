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
					panic(err)
				}
			}
		}()
	}
}

func (s *server) Read(wire io.ReadWriter) error {
	var req model.Request
	err := protocol.Read(wire, &req)
	if err != nil {
		s.socket.Close()
		return err
	}
	fmt.Println("header:", req)
	if req.Name == "" {
		err2 := protocol.Write(wire, model.NewError(-1, "Empty method"))
		if err2 != nil {
			// oups
			fmt.Println(err2)
		}
		return nil
	}
	h, ok := s.handlers[req.Name]
	if !ok {
		err2 := protocol.Write(wire, model.NewError(-1, "Unknown method: "+req.Name))
		if err2 != nil {
			// oups
			fmt.Println(err2)
		}
		return nil
	}
	res := h(&req)
	err = protocol.Write(wire, res)
	if err != nil {
		return err
	}
	return nil
}
