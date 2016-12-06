package server

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/bearstech/ascetic-rpc/model"
	"github.com/bearstech/ascetic-rpc/protocol"
)

type server struct {
	socket   *net.UnixListener
	handlers map[string]func(req *model.Request) (*model.Response, error)
	lock     sync.Mutex
}

func NewServer(socket *net.UnixListener) *server {
	return &server{
		socket:   socket,
		handlers: make(map[string]func(req *model.Request) (*model.Response, error)),
	}
}

func (s *server) Register(name string, fun func(req *model.Request) (*model.Response, error)) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.handlers[name] = fun
}

func (s *server) Deregister(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.handlers, name)
}

func (s *server) Listen() {
	// FIXME Handling quiet stop
	for {
		conn, err := s.socket.AcceptUnix()
		if err != nil {
			panic(err)
		}
		go func() {
			for {
				err := s.Read(conn)
				if err != nil {
					s.socket.Close()
					// FIXME it's error logging
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
		return protocol.Write(wire, model.NewErrorResponse(-1, "Empty method"))
	}
	h, ok := s.handlers[req.Name]
	if !ok {
		return protocol.Write(wire, model.NewErrorResponse(-1, "Unknown method: "+req.Name))
	}
	resp, err := h(&req)
	if err == nil {
		return protocol.Write(wire, resp)
	}
	return protocol.Write(wire, model.NewErrorResponse(-2, err.Error()))
}
