package http_server

import (
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/messages"
	"github.com/team-zf/framework/utils"
	"net/http"
	"time"
)

type SelectEvent struct {
	messages.HttpMessage

	Token    string // Token值
	ServerId int    // 服务器ID
}

func (e *SelectEvent) Parse() {
	e.Token = utils.NewStringAny(e.Params["token"]).ToString()
	e.ServerId = utils.NewStringAny(e.Params["server_id"]).ToIntV()
}

func (e *SelectEvent) HttpDirectCall(req *http.Request, resp *messages.HttpResponse) {
	logger.Debug("Token: %s", e.Token)
	logger.Debug("ServerId: %s", e.ServerId)

	time.Sleep(time.Second * 4)
	return
}

func M_Select() *SelectEvent {
	return &SelectEvent{}
}
