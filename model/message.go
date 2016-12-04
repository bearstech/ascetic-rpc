package model

import (
	"errors"

	"github.com/golang/protobuf/proto"
)

// Request

func (r *Request) SetBody(body proto.Message) error {
	blob, err := proto.Marshal(body)
	if err != nil {
		return err
	}
	r.RawBody = blob
	return nil
}

func (r *Request) GetBody(body proto.Message) error {
	return proto.Unmarshal(r.RawBody, body)
}

func NewRequest(name string, body proto.Message) (*Request, error) {
	r := &Request{
		Name: name,
	}
	err := r.SetBody(body)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Response

func (r *Response) SetOK(body proto.Message) error {
	blob, err := proto.Marshal(body)
	if err != nil {
		return err
	}
	r.Body = &Response_RawOK{RawOK: blob}
	return nil
}

func (r *Response) ReadOK(body proto.Message) error {
	return proto.Unmarshal(r.GetRawOK(), body)
}

func (r *Response) SetError(e *Error) {
	r.Body = &Response_Error{Error: e}
}

func NewError(code int32, message string) *Response {
	return &Response{
		Code: code,
		Body: &Response_Error{
			Error: &Error{
				Message: message,
			},
		},
	}
}

func NewOK(code int32, body proto.Message) (*Response, error) {
	r := &Response{Code: code}
	err := r.SetOK(body)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Response) GetErrorError() error {
	if r.Code >= 0 {
		return nil
	}
	return errors.New(r.GetError().Message)
}
