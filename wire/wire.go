package wire

import (
	"bytes"
	"io"
)

type mockWire struct {
	in  *bytes.Buffer
	out *bytes.Buffer
}

func New() *mockWire {
	return &mockWire{
		in:  new(bytes.Buffer),
		out: new(bytes.Buffer),
	}
}

type wire struct {
	a *bytes.Buffer
	b *bytes.Buffer
}

func (w *wire) Read(p []byte) (n int, err error) {
	return w.a.Read(p)
}

func (w *wire) Write(p []byte) (n int, err error) {
	return w.b.Write(p)
}

func (m *mockWire) ClientToServer() io.ReadWriter {
	return &wire{m.in, m.out}
}

func (m *mockWire) ServerToClient() io.ReadWriter {
	return &wire{m.out, m.in}
}
