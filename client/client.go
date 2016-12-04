package client

import (
	"fmt"
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
	req.SetBody(arg)
	err := protocol.Write(c.wire, &req)
	if err != nil {
		fmt.Println("Error while writing request", err)
		return err
	}

	var resp model.Response
	err = protocol.Read(c.wire, &resp)
	if err != nil {
		fmt.Println("Error while reading response", err)
		return err
	}
	if resp.Code < 0 { // it's an error
		return resp.GetErrorError()
	}
	return resp.ReadOK(r)
}
