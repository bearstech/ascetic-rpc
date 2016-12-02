package client

import (
	"errors"
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

func (c *client) Do(fun string, arg proto.Message, resp proto.Message) error {
	req := model.Request{
		Name: fun,
	}
	err := protocol.WriteHeaderAndBody(c.wire, &req, arg)
	if err != nil {
		fmt.Println("Error while writing request", err)
		return err
	}

	var resp_h model.Response
	err = protocol.Read(c.wire, &resp_h)
	if err != nil {
		fmt.Println("Error while reading response", err)
		return err
	}
	if resp_h.Code < 0 { // it's an error
		return errors.New(resp_h.Message)
	}
	return protocol.Read(c.wire, resp)
}
