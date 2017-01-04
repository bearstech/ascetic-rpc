package message

import "github.com/golang/protobuf/proto"

// Error

type TypedError interface {
	Type() ErrorType
}

type RpcError struct {
	message string
	_type   ErrorType
}

func (r *RpcError) Error() string {
	return r.message
}

func (r *RpcError) Type() ErrorType {
	return r._type
}

func (r *RpcError) String() string {
	return "<RpcError " + ErrorType_name[int32(r._type)] + " " + r.message + " >"
}

func NewRpcError(err ErrorType, message string) *RpcError {
	return &RpcError{
		message: message,
		_type:   err,
	}
}

func NewApplicationError(err error) *RpcError {
	return &RpcError{
		message: err.Error(),
		_type:   Error_APPLICATION,
	}
}

func (e *Error) Error() string {
	return e.Message
}

// Request

func (r *Request) SetBody(body proto.Message) error {
	if body == nil {
		r.RawBody = []byte{}
		return nil
	}
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
	var blob []byte
	var err error
	if body != nil {
		blob, err = proto.Marshal(body)
		if err != nil {
			return err
		}
	} else {
		blob = []byte{}
	}
	r.Body = &Response_RawOK{RawOK: blob}
	return nil
}

func (r *Response) ReadOK(body proto.Message) error {
	if body == nil {
		return nil
	}
	return proto.Unmarshal(r.GetRawOK(), body)
}

func (r *Response) SetError(e *Error) {
	r.Body = &Response_Error{Error: e}
}

func NewErrorResponse(_type ErrorType, message string) *Response {
	return &Response{
		Body: &Response_Error{
			Error: &Error{
				Type:    _type,
				Message: message,
			},
		},
	}
}

func NewOKResponse(body proto.Message) (*Response, error) {
	r := &Response{}
	err := r.SetOK(body)
	if err != nil {
		return nil, err
	}
	return r, nil
}
