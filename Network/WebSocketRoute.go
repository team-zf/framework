package Network

import "fmt"

type WebSocketRoute struct {
	Cmd    uint32                 `json:"cmd"`
	Params map[string]interface{} `json:"params"`
	WebSocketDDM
}

func (e *WebSocketRoute) GetCmd() uint32 {
	return e.Cmd
}

func (e *WebSocketRoute) Header() string {
	return fmt.Sprintf("Cmd: %d, Params: %+v", e.Cmd, e.Params)
}
