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
	// 解析参数
	Parse()
	// WebSocket的回调
	WebSocketDirectCall(wsmd *model.WebSocketModel, resp *WebSocketResponse)
}

type IDataBaseMessage interface {
	// 所在DB协程
	DBThreadID() int
	// 数据表,如果你的表放入时，不是马上保存的，那么后续可以用这个KEY来进行覆盖，
	// 这样就可以实现多次修改一次保存的功能
	// 所以这个字段建议是：用户ID+数据表名+数据主键
	GetDataKey() string
	// 调用方法
	SaveDB() error
}
