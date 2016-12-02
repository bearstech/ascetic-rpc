package protocol

import (
	"bytes"
	"errors"
	"testing"

	"github.com/bearstech/ascetic-rpc/model"
)

func TestProtocol(t *testing.T) {
	wire := new(bytes.Buffer)
	req := model.Request{
		Name: "plop",
	}
	hello := model.Hello{
		Name: "Charles",
	}
	err := WriteHeaderAndBody(wire, &req, &hello)
	if err != nil {
		t.Error(err)
	}

	var r model.Request
	var h model.Hello
	err = Read(wire, &r)
	if err != nil {
		t.Error(err)
	}
	if r.Name != "plop" {
		t.Error(errors.New("Bad name"))
	}
	err = Read(wire, &h)
	if err != nil {
		t.Error(err)
	}
	if h.Name != "Charles" {
		t.Error(errors.New("Bad name"))
	}

}
