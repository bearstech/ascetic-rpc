package protocol

import (
	"bytes"
	"errors"
	"github.com/bearstech/ascetic-rpc/model"
	"testing"
)

func TestProtocol(t *testing.T) {
	wire := new(bytes.Buffer)
	req := model.Request{
		Name: "plop",
	}
	req.SetBody(&model.Hello{
		Name: "Charles",
	})
	err := Write(wire, &req)
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
	err = r.GetBody(&h)
	if err != nil {
		t.Error(err)
	}
	if h.Name != "Charles" {
		t.Error(errors.New("Bad name"))
	}

}
