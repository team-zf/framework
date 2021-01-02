package messages

import (
	"github.com/team-zf/framework/model"
	"net/http"
)

// 消息收发接口
type IMessageHandle interface {
	// 编码
	Marshal(data interface{}) ([]byte, error)
	// 解码
	Unmarshal(buff []byte) (data interface{}, err error)
	// 设置消息路由
	SetRoute(cmd uint32, msg interface{})
	// 按消息拿出消息处理实例
	GetRoute(cmd uint32) (result interface{}, err error)
	// 一个消息是否收完了
	// 返回这个消息应该的长度，和是否收完的信息
	CheckMaxLenVaild(buff []byte) (msglen uint32, ok bool)
}

type options func(msghandle IMessageHandle)

// 路由消息接口
type IGateMessage interface {
	// 编码，传出编码的数据和数据的长度
	GateMarshal() ([]byte, uint32)
	// 解码,传入数据，传出使用后剩下的数据，和使用了多少字节
	GateUnmarshal(buff []byte) ([]byte, uint32)
}

type IMessage interface {
	GetCmd() uint32
	Header() string
}

type IHttpMessageHandle interface {
	IMessage
	// 解析参数
	Parse()
	// HTTP的回调
	HttpDirectCall(req *http.Request, resp *HttpResponse)
}

type IWebSocketMessageHandle interface {
	IMessage
	// WebSocket的回调
	WebSocketDirectCall(*model.WebSocketModel)
}
