package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/bearstech/ascetic-rpc/model"
	"github.com/bearstech/ascetic-rpc/protocol"
)

// Graceful stop pattern came from :
// https://rcrowley.org/articles/golang-graceful-stop.html

type Deadliner interface {
	SetDeadline(time.Time) error
}

type ReadWriteCloseDeadliner interface {
	io.ReadWriteCloser
	Deadliner
}

type server struct {
	socket    *net.UnixListener
	handlers  map[string]func(req *model.Request) (*model.Response, error)
	lock      sync.Mutex
	ch        chan bool
	waitGroup *sync.WaitGroup
}

func NewServer(socket *net.UnixListener) *server {
	return &server{
		socket:    socket,
		handlers:  make(map[string]func(req *model.Request) (*model.Response, error)),
		ch:        make(chan bool),
		waitGroup: &sync.WaitGroup{},
	}
}

func NewServerUnix(socketPath string) (*server, error) {
	err := os.Remove(socketPath)
	// FIXME Handle error other than "file not exist"
	l, err := net.ListenUnix("unix", &net.UnixAddr{
		Name: socketPath,
		Net:  "unix",
	})
	if err != nil {
		return nil, err
	}
	return NewServer(l), nil
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
	s.waitGroup.Add(1)
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch:
			s.socket.Close()
			return
		default:
		}
		if s.socket == nil {
			return
		}
		s.socket.SetDeadline(time.Now().Add(1e9))
		conn, err := s.socket.AcceptUnix()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			panic(err)
		}
		s.waitGroup.Add(1)
		go s.HandleSession(conn)
	}
}

func (s *server) HandleSession(wire ReadWriteCloseDeadliner) error {
	defer wire.Close()
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch:
			return nil
		default:
		}
		wire.SetDeadline(time.Now().Add(1e9))
		err := s.Handle(wire)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			// FIXME it's error logging
			fmt.Println("Handle error", err)
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
	if req.Name == "" {
		return protocol.Write(wire, model.NewErrorResponse(model.Error_BAD_METHOD, "Empty method"))
	}
	h, ok := s.handlers[req.Name]
	if !ok {
		return protocol.Write(wire, model.NewErrorResponse(model.Error_BAD_METHOD, "Unknown method: "+req.Name))
	}

	var resp *model.Response

	func() {
		defer func() {
			if r := recover(); r != nil {
				if er, ok := r.(error); ok {
					err = er
				} else {
					fmt.Println("Panic :", r)
					err = errors.New("Uncatchable error")
				}
				resp = nil
			}
		}()
		resp, err = h(&req)
	}()
	if err != nil {
		return protocol.Write(wire, model.NewErrorResponse(model.Error_APPLICATION, err.Error()))
	}
	if resp == nil { // It's a lazy answer, but I can handle it.
		resp = &model.Response{}
	}
	return protocol.Write(wire, resp)
}

func (s *server) Stop() {
	close(s.ch)
	s.waitGroup.Wait()
}

func Ping(req *model.Request) (*model.Response, error) {
	// assert body is nil
	return nil, nil
}
