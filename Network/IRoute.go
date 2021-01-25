package Network

import "net/http"

type IRoute interface {
	GetCmd() uint32
	Header() string
}

type IHttpRoute interface {
	IRoute
	// 解析参数
	Parse()
	// 处理业务逻辑, 并返回结果代码
	Handle(req *http.Request) uint32
	// 输出JsonMap
	ToJsonMap() map[string]interface{}
}

type IWebSocketRoute interface {
	IRoute
	// 解析参数
	Parse()
	// 处理业务逻辑, 并返回结果代码
	Handle(agent *WebSocketAgent) uint32
	// 输出JsonMap
	ToJsonMap() map[string]interface{}
}
