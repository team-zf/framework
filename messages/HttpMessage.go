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
