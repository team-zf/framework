package messages

import (
	"fmt"
)

type HttpMessage struct {
	Cmd    uint32                 `json:"cmd"`
	Params map[string]interface{} `json:"params"`
}

func (e *HttpMessage) GetCmd() uint32 {
	return e.Cmd
}

func (e *HttpMessage) Header() string {
	return fmt.Sprintf("Cmd: %d, Params: %+v", e.Cmd, e.Params)
}

type HttpResponse struct {
	Code uint32                 `json:"code"`
	Data map[string]interface{} `json:"data"`
}

type Options func(resp *HttpResponse)

func ResponseSetCode(v uint32) Options {
	return func(resp *HttpResponse) {
		resp.Code = v
	}
}

func ResponseSetData(v map[string]interface{}) Options {
	return func(resp *HttpResponse) {
		resp.Data = v
	}
}

func NewHttpResponse(opts ...Options) *HttpResponse {
	resp := &HttpResponse{
		Code: 0,
		Data: make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(resp)
	}
	return resp
}

func NewCodeHttpResponse(code uint32) *HttpResponse {
	return &HttpResponse{
		Code: code,
		Data: make(map[string]interface{}),
	}
}
