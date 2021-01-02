package http_server

import (
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/messages"
	"github.com/team-zf/framework/utils"
	"net/http"
	"time"
)

type ListEvent struct {
	messages.HttpMessage

	Token string // Token值
}

func (e *ListEvent) Parse() {
	e.Token = utils.NewStringAny(e.Params["token"]).ToString()
}

func (e *ListEvent) HttpDirectCall(req *http.Request, resp *messages.HttpResponse) {
	logger.Debug("Token: %s", e.Token)

	time.Sleep(time.Second * 4)
	return
}

func M_List() *ListEvent {
	return &ListEvent{}
}
