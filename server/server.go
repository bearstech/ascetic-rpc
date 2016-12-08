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
	ch       chan bool
}

func NewServer(socket *net.UnixListener) *server {
	return &server{
		socket:   socket,
		handlers: make(map[string]func(req *model.Request) (*model.Response, error)),
		ch:       make(chan bool),
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

func (s *server) Serve() {
	for {
		select {
		case <-s.ch:
			s.socket.Close()
			return
		default:
		}
		if s.socket == nil {
			fmt.Println("No more socket")
			return
		}
		conn, err := s.socket.AcceptUnix()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		go s.HandleSession(conn)
	}
}

func (s *server) HandleSession(wire io.ReadWriteCloser) error {
	for {
		err := s.Handle(wire)
		if err != nil {
			// FIXME it's error logging
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}

func (s *server) Handle(wire io.ReadWriteCloser) error {
	var req model.Request
	err := protocol.Read(wire, &req)
	if err != nil {
		wire.Close()
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

func (s *server) Stop() {
	close(s.ch)
	fmt.Println("Stopped")
}
