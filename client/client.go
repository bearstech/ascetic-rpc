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

func (c *Client) Close() error {
	return c.wire.Close()
}
