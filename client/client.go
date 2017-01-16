package client

import (
	"errors"
	"io"
	"net"
	"strings"
	"time"

	"github.com/bearstech/ascetic-rpc/message"
	"github.com/bearstech/ascetic-rpc/protocol"
	"github.com/golang/protobuf/proto"
)

type Client struct {
	wire io.ReadWriteCloser
}

func New(wire io.ReadWriteCloser) *Client {
	return &Client{
		wire: wire,
	}
}

func NewClientUnix(socketPath string) (*Client, error) {
	for i := int64(0); i < 4; i++ {
		conn, err := net.DialUnix("unix", nil, &net.UnixAddr{
			Name: socketPath,
			Net:  "unix"})
		if err != nil {
			if strings.HasSuffix(err.Error(), "connect: no such file or directory") {
				time.Sleep(time.Duration(i*100) * time.Millisecond)
				continue
			}
			return nil, err
		}
		return New(conn), nil
	}
	return nil, errors.New("Too many connections attempt on socket : " + socketPath)
}

func (c *Client) Do(fun string, arg proto.Message, r proto.Message) error {
	req := message.Request{
		Name: fun,
	}
	err := req.SetBody(arg)
	if err != nil {
		return err
	}
	err = protocol.Write(c.wire, &req)
	if err != nil {
		return err
	}

	var resp message.Response
	err = protocol.Read(c.wire, &resp)
	if err != nil {
		return err
	}
	e := resp.GetError()
	if e != nil { // it's an error
		return e
	}
	return resp.ReadOK(r)
}

func (c *Client) Stream(fun string, arg proto.Message) (Streamer, error) {
	req := message.Request{
		Name: fun,
	}
	err := req.SetBody(arg)
	if err != nil {
		return nil, err
	}
	err = protocol.Write(c.wire, &req)
	if err != nil {
		return nil, err
	}

	var resp message.Response
	err = protocol.Read(c.wire, &resp)
	if err != nil {
		return nil, err
	}
	e := resp.GetError()
	if e != nil { // it's an error
		return nil, e
	}

	if !resp.GetStream() {
		return nil, errors.New("It's not a stream")
	}

	return &simpleStramer{wire: c.wire}, nil
}

func (c *Client) Close() error {
	return c.wire.Close()
}

type Streamer interface {
	Recv(proto.Message) error
}

type simpleStramer struct {
	wire io.ReadWriteCloser
}

func (s *simpleStramer) Recv(r proto.Message) error {
	var chunk message.Chunk
	err := protocol.Read(s.wire, &chunk)
	if err != nil {
		return err
	}
	e := chunk.GetError()
	if e != nil {
		return e
	}

	if chunk.GetEOF() {
		return io.EOF
	}

	return proto.Unmarshal(chunk.GetRawOK(), r)
}
