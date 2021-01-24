package Network

import "golang.org/x/net/websocket"

type WebSocketAgent struct {
	Conn        *websocket.Conn
	routeHandle *WebSocketRouteHandle
}

func (e *WebSocketAgent) SendData(data interface{}) error {
	buff, err := e.routeHandle.Marshal(data)
	if err != nil {
		return err
	}
	return e.SendByte(buff)
}

func (e *WebSocketAgent) SendByte(buff []byte) error {
	return websocket.Message.Send(e.Conn, buff)
}
