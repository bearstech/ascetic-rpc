// Code generated by protoc-gen-go.
// source: model/message.proto
// DO NOT EDIT!

/*
Package model is a generated protocol buffer package.

It is generated from these files:
	model/message.proto
	model/test.proto

It has these top-level messages:
	Request
	Response
	Error
	Hello
	World
*/
package model

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Request struct {
	Name      string `protobuf:"bytes,1,opt,name=Name" json:"Name,omitempty"`
	Signature string `protobuf:"bytes,2,opt,name=Signature" json:"Signature,omitempty"`
	RawBody   []byte `protobuf:"bytes,3,opt,name=RawBody,proto3" json:"RawBody,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Response struct {
	Code      int32  `protobuf:"varint,1,opt,name=Code" json:"Code,omitempty"`
	Signature string `protobuf:"bytes,2,opt,name=Signature" json:"Signature,omitempty"`
	// Types that are valid to be assigned to Body:
	//	*Response_Error
	//	*Response_RawOK
	Body isResponse_Body `protobuf_oneof:"Body"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type isResponse_Body interface {
	isResponse_Body()
}

type Response_Error struct {
	Error *Error `protobuf:"bytes,3,opt,name=Error,oneof"`
}
type Response_RawOK struct {
	RawOK []byte `protobuf:"bytes,4,opt,name=RawOK,proto3,oneof"`
}

func (*Response_Error) isResponse_Body() {}
func (*Response_RawOK) isResponse_Body() {}

func (m *Response) GetBody() isResponse_Body {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *Response) GetError() *Error {
	if x, ok := m.GetBody().(*Response_Error); ok {
		return x.Error
	}
	return nil
}

func (m *Response) GetRawOK() []byte {
	if x, ok := m.GetBody().(*Response_RawOK); ok {
		return x.RawOK
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Response) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Response_OneofMarshaler, _Response_OneofUnmarshaler, _Response_OneofSizer, []interface{}{
		(*Response_Error)(nil),
		(*Response_RawOK)(nil),
	}
}

func _Response_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Response)
	// Body
	switch x := m.Body.(type) {
	case *Response_Error:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Error); err != nil {
			return err
		}
	case *Response_RawOK:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		b.EncodeRawBytes(x.RawOK)
	case nil:
	default:
		return fmt.Errorf("Response.Body has unexpected type %T", x)
	}
	return nil
}

func _Response_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Response)
	switch tag {
	case 3: // Body.Error
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Error)
		err := b.DecodeMessage(msg)
		m.Body = &Response_Error{msg}
		return true, err
	case 4: // Body.RawOK
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeRawBytes(true)
		m.Body = &Response_RawOK{x}
		return true, err
	default:
		return false, nil
	}
}

func _Response_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Response)
	// Body
	switch x := m.Body.(type) {
	case *Response_Error:
		s := proto.Size(x.Error)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Response_RawOK:
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.RawOK)))
		n += len(x.RawOK)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Error struct {
	Message string `protobuf:"bytes,1,opt,name=Message" json:"Message,omitempty"`
}

func (m *Error) Reset()                    { *m = Error{} }
func (m *Error) String() string            { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()               {}
func (*Error) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func init() {
	proto.RegisterType((*Request)(nil), "model.Request")
	proto.RegisterType((*Response)(nil), "model.Response")
	proto.RegisterType((*Error)(nil), "model.Error")
}

func init() { proto.RegisterFile("model/message.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 204 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0xce, 0xcd, 0x4f, 0x49,
	0xcd, 0xd1, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0x62, 0x05, 0x0b, 0x2a, 0x85, 0x72, 0xb1, 0x07, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x09,
	0x71, 0xb1, 0xf8, 0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x42,
	0x32, 0x5c, 0x9c, 0xc1, 0x99, 0xe9, 0x79, 0x89, 0x25, 0xa5, 0x45, 0xa9, 0x12, 0x4c, 0x60, 0x09,
	0x84, 0x80, 0x90, 0x04, 0x50, 0x73, 0x62, 0xb9, 0x53, 0x7e, 0x4a, 0xa5, 0x04, 0x33, 0x50, 0x8e,
	0x27, 0x08, 0xc6, 0x55, 0x6a, 0x62, 0xe4, 0xe2, 0x08, 0x4a, 0x2d, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e,
	0x05, 0x19, 0xec, 0x0c, 0xb4, 0x0c, 0x6c, 0x30, 0x6b, 0x10, 0x98, 0x4d, 0xc0, 0x60, 0x15, 0x2e,
	0x56, 0xd7, 0xa2, 0xa2, 0xfc, 0x22, 0xb0, 0xb1, 0xdc, 0x46, 0x3c, 0x7a, 0x60, 0xc7, 0xea, 0x81,
	0xc5, 0x3c, 0x18, 0x82, 0x20, 0x92, 0x42, 0x62, 0x5c, 0xac, 0x40, 0xfb, 0xfc, 0xbd, 0x25, 0x58,
	0x40, 0x96, 0x83, 0xc4, 0xc1, 0x5c, 0x27, 0x36, 0x2e, 0x16, 0xb0, 0x23, 0x14, 0xa1, 0xa6, 0x80,
	0xdc, 0xe9, 0x0b, 0xf1, 0x3c, 0xd4, 0x73, 0x30, 0x6e, 0x12, 0x1b, 0x38, 0x30, 0x8c, 0x01, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x43, 0x4f, 0x12, 0x11, 0x23, 0x01, 0x00, 0x00,
}
