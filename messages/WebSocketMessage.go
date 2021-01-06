package messages

import "fmt"

type WebSocketMessage struct {
	Cmd    uint32                 `json:"cmd"`
	Params map[string]interface{} `json:"params"`
}

func (e *WebSocketMessage) GetCmd() uint32 {
	return e.Cmd
}

func (e *WebSocketMessage) Header() string {
	return fmt.Sprintf("Cmd: %d, Params: %+v", e.Cmd, e.Params)
}

type WebSocketResponse struct {
	Cmd  uint32                 `json:"cmd"`
	Code uint32                 `json:"code"`
	Data map[string]interface{} `json:"data"`
}
