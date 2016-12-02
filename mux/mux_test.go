package mux

import (
	"bytes"
	"testing"

	"github.com/bearstech/ascetic-rpc/model"
	"github.com/bearstech/ascetic-rpc/protocol"
	"github.com/golang/protobuf/proto"
)

type ping struct{}

func (p ping) Handle(*model.Request, []byte) (model.Response, proto.Message) {
	return model.Response{Code: 1}, nil
}

type mockClient struct {
	in  *bytes.Buffer
	out *bytes.Buffer
}

func (m *mockClient) Read(p []byte) (n int, err error) {
	return m.in.Read(p)
}

func (m *mockClient) Write(p []byte) (n int, err error) {
	return m.out.Write(p)
}

func newMockClient() *mockClient {
	return &mockClient{
		in:  new(bytes.Buffer),
		out: new(bytes.Buffer),
	}
}

func TestPing(t *testing.T) {
	wire := newMockClient()
	s := NewServer(wire)
	s.Route("ping", ping{})

	req := model.Request{
		Name: "ping",
	}
	err := protocol.Write(wire.in, &req)
	if err != nil {
		t.Error(err)
	}
	err = protocol.Write(wire.in, nil)
	if err != nil {
		t.Error(err)
	}
	err = s.Read()
	if err != nil {
		t.Error(err)
	}

	var resp model.Response
	err = protocol.Read(wire.out, &resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Code != 1 {
		t.Fatal()
	}
}
