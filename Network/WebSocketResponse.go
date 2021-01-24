package Network

type WebSocketResponse struct {
	Cmd  uint32 `json:"cmd"`
	Code uint32 `json:"code"`
}
