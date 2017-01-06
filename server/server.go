package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/bearstech/ascetic-rpc/message"
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
	socket         *net.UnixListener
	handlers       map[string]*method
	lock           sync.Mutex
	ch             chan bool
	waitGroup      *sync.WaitGroup
	running        bool
	timeout        time.Duration
	queryTimeout   time.Duration
	sessionTimeout time.Duration
}

type method struct {
	name     string
	function func(req *message.Request) (*message.Response, error)
	timeout  time.Duration
}

func (m *method) WithTimeout(t time.Duration) *method {
	m.timeout = t
	return m
}

func NewServer(socket *net.UnixListener) *server {
	return &server{
		socket:         socket,
		handlers:       make(map[string]*method),
		ch:             make(chan bool),
		waitGroup:      &sync.WaitGroup{},
		running:        false,
		timeout:        time.Second, // function timeout
		queryTimeout:   100 * time.Millisecond,
		sessionTimeout: 10 * time.Second,
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

func (s *server) Register(name string, fun func(req *message.Request) (*message.Response, error)) *method {
	s.lock.Lock()
	defer s.lock.Unlock()
	m := &method{
		name:     name,
		function: fun,
		timeout:  s.timeout,
	}
	s.handlers[name] = m
	return m
}

func (s *server) Deregister(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.handlers, name)
}

func (s *server) Serve() {
	s.waitGroup.Add(1)
	s.running = true
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
		s.socket.SetDeadline(time.Now().Add(s.queryTimeout))
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
	wire.SetDeadline(time.Now().Add(s.queryTimeout))
	defer wire.Close()
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch:
			return nil
		default:
		}
		err := s.Handle(wire)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			if err == io.EOF { // client deconnects
				return nil
			}
			// FIXME it's error logging
			fmt.Println("Handle error", err)
			return err
		}
		wire.SetDeadline(time.Now().Add(s.sessionTimeout))
	}
}

func (s *server) Handle(wire ReadWriteCloseDeadliner) error {
	wire.SetDeadline(time.Now().Add(s.queryTimeout))

	var req message.Request
	err := protocol.Read(wire, &req)
	if err != nil {
		wire.Close()
		return err
	}
	if req.Name == "" {
		return protocol.Write(wire, message.NewErrorResponse(message.Error_BAD_METHOD, "Empty method"))
	}
	h, ok := s.handlers[req.Name]
	if !ok {
		return protocol.Write(wire, message.NewErrorResponse(message.Error_BAD_METHOD, "Unknown method: "+req.Name))
	}

	wire.SetDeadline(time.Now().Add(h.timeout))

	var resp *message.Response

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
		resp, err = h.function(&req)
	}()
	wire.SetDeadline(time.Now().Add(s.queryTimeout))
	if err != nil {
		return protocol.Write(wire, message.NewErrorResponse(message.Error_APPLICATION, err.Error()))
	}
	if resp == nil { // It's a lazy answer, but I can handle it.
		resp = &message.Response{}
	}
	return protocol.Write(wire, resp)
}

func (s *server) Stop() error {
	if !s.running {
		return errors.New("Server is not running")
	}
	close(s.ch)
	s.waitGroup.Wait()
	s.running = false
	return nil
}

func (s *server) IsRunning() bool {
	return s.running
}

func Ping(req *message.Request) (*message.Response, error) {
	// assert body is nil
	return nil, nil
}
