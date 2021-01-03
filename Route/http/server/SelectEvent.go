package http_server

import (
	"github.com/team-zf/framework/Data"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/messages"
	"github.com/team-zf/framework/utils"
	"net/http"
)

type SelectEvent struct {
	messages.HttpMessage

	Token    string // Token值
	ServerId int64  // 服务器ID
}

func (e *SelectEvent) Parse() {
	e.Token = utils.NewStringAny(e.Params["token"]).ToString()
	e.ServerId = utils.NewStringAny(e.Params["server_id"]).ToInt64V()
}

func (e *SelectEvent) HttpDirectCall(req *http.Request, resp *messages.HttpResponse) {
	account := Data.GetAccountByToken(e.Token)

	// Token错误
	if account == nil {
		logger.Debug("Token错误")
		return
	}

	// 该账户已被禁用, 也许是被封了
	if account.Status == 1 {
		logger.Debug("该账户已被禁用, 也许是被封了")
		return
	}

	server := Data.GetServerById(e.ServerId)

	// 该服不存在
	if server == nil {
		logger.Debug("该服不存在")
		return
	}

	logger.Debug("请及时连接逻辑服.")

	resp.Data["host"] = server.Host
	resp.Data["port"] = server.Port
	resp.Code = messages.RC_Success
	return
}

func M_Select() *SelectEvent {
	return &SelectEvent{}
}
