package http_account

import (
	"github.com/team-zf/framework/Data"
	"github.com/team-zf/framework/Table"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/messages"
	"github.com/team-zf/framework/utils"
	"net/http"
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
	// 该账户已存在
	if Data.GetAccountByUserName(e.UserName) != nil {
		logger.Debug("该账户已存在")
		return
	}

	var id int64
	var token string
	var ok bool

	id, ok = Data.GenerateAccountId()
	// 账户ID生成失败
	if !ok {
		logger.Debug("账户ID生成失败")
		return
	}

	token, ok = Data.GenerateToken()
	// 账户Token生成失败
	if !ok {
		logger.Debug("账户Token生成失败")
		return
	}

	account := Table.NewAccount()
	account.Id = id
	account.UserName = e.UserName
	account.PassWord = e.PassWord
	account.Token = token
	if Data.RegisterAccount(account) {
		logger.Debug("注册成功")
		resp.Code = messages.RC_Success
	} else {
		logger.Debug("注册失败")
		resp.Code = messages.RC_LOGIC_ERROR
	}

	return
}

func M_Register() *RegisterEvent {
	return &RegisterEvent{}
}
