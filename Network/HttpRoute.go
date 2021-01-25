package Network

import "fmt"

type HttpRoute struct {
	Cmd    uint32                 `json:"cmd"`
	Params map[string]interface{} `json:"params"`
	WebSocketDDM
}

func (e *HttpRoute) GetCmd() uint32 {
	return e.Cmd
}

func (e *HttpRoute) Header() string {
	return fmt.Sprintf("Cmd: %d, Params: %+v", e.Cmd, e.Params)
}
