package http_server

import (
	"github.com/team-zf/framework/Data"
	"github.com/team-zf/framework/Table"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/messages"
	"github.com/team-zf/framework/utils"
	"net/http"
	"strings"
)

type ListEvent struct {
	messages.HttpMessage

	Token string // Token值
}

func (e *ListEvent) Parse() {
	e.Token = utils.NewStringAny(e.Params["token"]).ToString()
}

func (e *ListEvent) HttpDirectCall(req *http.Request, resp *messages.HttpResponse) {
	account := Data.GetAccountByToken(e.Token)

	// Token错误
	if account == nil {
		logger.Debug("Token错误")
		return
	}

	serverList := Data.GetServerList()
	serverIds := strings.Split(account.LatelyServer, ",")

	latelyList := make([]*Table.Server, 0)
	for _, sid := range serverIds {
		for _, server := range serverList {
			if string(server.Id) == sid {
				latelyList = append(latelyList, server)
				break
			}
		}
	}

	// 最近登录的服务器列表
	if len(latelyList) > 0 {
		mlist := make([]map[string]interface{}, 0)
		for _, item := range latelyList {
			mlist = append(mlist, item.ToJsonMap())
		}
		resp.Data["lately"] = mlist
	}
	// 所有服务器列表
	if len(serverList) > 0 {
		mlist := make([]map[string]interface{}, 0)
		for _, item := range serverList {
			mlist = append(mlist, item.ToJsonMap())
		}
		resp.Data["list"] = mlist
	}

	resp.Code = messages.RC_Success
	return
}

func M_List() *ListEvent {
	return &ListEvent{}
}
