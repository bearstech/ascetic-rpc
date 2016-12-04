package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
)

func Write(wire io.Writer, msg proto.Message) error {
	if msg == nil {
		return binary.Write(wire, binary.LittleEndian, uint16(0))
	}
	txt, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	size := len(txt)
	if size >= 65536 {
		return fmt.Errorf("Message is too big : %d >= 65536", size)
	}
	err = binary.Write(wire, binary.LittleEndian, uint16(size))
	if err != nil {
		return err
	}
	s, err := wire.Write([]byte(txt))
	if err != nil {
		return err
	}
	if s < size {
		return errors.New("Partial write")
	}
	return nil
}

func ReadBytes(wire io.Reader) ([]byte, error) {
	var size uint16
	err := binary.Read(wire, binary.LittleEndian, &size)
	if err != nil {
		return nil, err
	}
	if size == 0 {
		return []byte{}, nil
	}
	buf := make([]byte, size)
	s, err := io.ReadFull(wire, buf)
	if err != nil {
		return nil, err
	}
	if uint16(s) < size {
		return nil, errors.New("Partial read")
	}
	return buf, nil
}

func Read(wire io.Reader, msg proto.Message) error {
	buf, err := ReadBytes(wire)
	if err != nil {
		return err
	}
	if msg == nil {
		if len(buf) != 0 {
			return errors.New("Nil message should be empty")
		}
		return nil
	}
	return proto.Unmarshal(buf, msg)
}
