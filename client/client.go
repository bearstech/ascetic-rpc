package client

import (
	"io"

	"github.com/bearstech/ascetic-rpc/model"
	"github.com/bearstech/ascetic-rpc/protocol"
	"github.com/golang/protobuf/proto"
)

type client struct {
	wire io.ReadWriter
}

func New(wire io.ReadWriter) *client {
	return &client{
		wire: wire,
	}
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
