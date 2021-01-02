package http_account

import (
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/messages"
	"github.com/team-zf/framework/utils"
	"net/http"
	"time"
)

type RegisterEvent struct {
	messages.HttpMessage

	UserName string // 账号
	PassWord string // 密码
}

func (e *RegisterEvent) Parse() {
	e.UserName = utils.NewStringAny(e.Params["username"]).ToString()
	e.PassWord = utils.NewStringAny(e.Params["password"]).ToString()
}

func (e *RegisterEvent) HttpDirectCall(req *http.Request, resp *messages.HttpResponse) {
	logger.Debug("UserName: %s", e.UserName)
	logger.Debug("PassWord: %s", e.PassWord)

	time.Sleep(time.Second * 4)
	return
}

func M_Register() *RegisterEvent {
	return &RegisterEvent{}
}
