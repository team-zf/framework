package http_account

import (
	"github.com/team-zf/framework/Data"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/messages"
	"github.com/team-zf/framework/utils"
	"net/http"
	"strings"
)

type TokenLoginEvent struct {
	messages.HttpMessage

	Token string // Token值
}

func (e *TokenLoginEvent) Parse() {
	e.Token = utils.NewStringAny(e.Params["token"]).ToString()
}

func (e *TokenLoginEvent) HttpDirectCall(req *http.Request, resp *messages.HttpResponse) {
	account := Data.GetAccountByToken(e.Token)

	// Token错误
	if account == nil {
		logger.Debug("Token错误")
		return
	}

	logger.Debug("Token登录成功")

	// 账户信息
	resp.Data["account"] = account.ToJsonMap()

	// 默认选中的服务器
	serverList := Data.GetServerList()
	if account.LatelyServer == "" {
		resp.Data["server"] = serverList[0].ToJsonMap()
	} else {
		serverId := strings.Split(account.LatelyServer, ",")[0]
		for _, server := range serverList {
			if string(server.Id) == serverId {
				resp.Data["server"] = serverList[0].ToJsonMap()
				break
			}
		}
	}
	resp.Code = messages.RC_Success
}

func M_TokenLogin() *TokenLoginEvent {
	return &TokenLoginEvent{}
}
