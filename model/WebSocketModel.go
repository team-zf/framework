package model

import "golang.org/x/net/websocket"

type IWebSocketModel interface {
}

type WebSocketModel struct {
	Conn     *websocket.Conn
	CloseFun func(wsmd *WebSocketModel) // 关闭连接时的方法
	ConInfo  interface{}                // 自定义的连接信息，给上层逻辑使用
	KeyID    int                        // 用来标记的ID
}

//发的是字符串
func (e *WebSocketModel) SendStr(data string) error {
	return websocket.Message.Send(e.Conn, data)
}

//发的是二进制数据
func (e *WebSocketModel) SendByte(data []byte) error {
	return websocket.Message.Send(e.Conn, data)
}
