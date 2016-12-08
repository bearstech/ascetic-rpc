package client

import (
	"errors"
	"io"
	"net"
	"strings"
	"time"

	"github.com/bearstech/ascetic-rpc/model"
	"github.com/bearstech/ascetic-rpc/protocol"
	"github.com/golang/protobuf/proto"
)

type client struct {
	wire io.ReadWriteCloser
}

func New(wire io.ReadWriteCloser) *client {
	return &client{
		wire: wire,
	}
}

func NewClientUnix(socketPath string) (*client, error) {
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
	return nil, errors.New("Too many connections attempt")
}

func (c *client) Do(fun string, arg proto.Message, r proto.Message) error {
	req := model.Request{
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

	var resp model.Response
	err = protocol.Read(c.wire, &resp)
	if err != nil {
		return err
	}
	if resp.Code < 0 { // it's an error
		return resp.GetErrorError()
	}
	return resp.ReadOK(r)
}

func (c *client) Close() error {
	return c.wire.Close()
}
