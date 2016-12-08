package wire

import (
	"bytes"
	"errors"
	"io"
)

type mockWire struct {
	in     *bytes.Buffer
	out    *bytes.Buffer
	closed bool
}

func New() *mockWire {
	return &mockWire{
		in:     new(bytes.Buffer),
		out:    new(bytes.Buffer),
		closed: false,
	}
}

type wire struct {
	a    *bytes.Buffer
	b    *bytes.Buffer
	mock *mockWire
}

func (w *wire) Read(p []byte) (n int, err error) {
	if w.mock.closed {
		return 0, errors.New("The wire is closed")
	}
	return w.a.Read(p)
}

func (w *wire) Write(p []byte) (n int, err error) {
	if w.mock.closed {
		return 0, errors.New("The wire is closed")
	}
	return w.b.Write(p)
}

func (w *wire) Close() error {
	w.mock.closed = true
	return nil
}

func (m *mockWire) ClientToServer() io.ReadWriteCloser {
	return &wire{m.in, m.out, m}
}

func (m *mockWire) ServerToClient() io.ReadWriteCloser {
	return &wire{m.out, m.in, m}
}

func (m *mockWire) Len() (int, int) {
	return m.in.Len(), m.out.Len()
}
