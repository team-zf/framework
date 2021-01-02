package Route

import (
	"github.com/team-zf/framework/Route/http/account"
	http_server "github.com/team-zf/framework/Route/http/server"
	"github.com/team-zf/framework/messages"
)

var (
	HttpRoute = messages.NewHttpMessageHandle()
)

const (
	Account_Register      uint32 = 1001 // 账户注册
	Account_PasswordLogin uint32 = 1002 // 账户密码登录
	Account_TokenLogin    uint32 = 1003 // 账户Token登录
	Server_List           uint32 = 2001 // 服务器列表
	Server_Select         uint32 = 2002 // 选服
)

func init() {
	HttpRoute.SetRoute(Account_Register, http_account.M_Register())
	HttpRoute.SetRoute(Account_PasswordLogin, http_account.M_PasswordLogin())
	HttpRoute.SetRoute(Account_TokenLogin, http_account.M_TokenLogin())
	HttpRoute.SetRoute(Server_List, http_server.M_List())
	HttpRoute.SetRoute(Server_Select, http_server.M_Select())
}
