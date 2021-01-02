package messages

// 服务器响应代码, Http和WebSocket通用
const (
	RC_NsqRequestFailure = -6 // 请求失败
	RC_NotLogic          = -5 // 没有逻辑处理它
	RC_NotCmd            = -4 // 没有事件处理这个消息
	RC_NotResult         = -3 // 不回复
	RC_LOGIC_ERROR       = -2 // 逻辑处理错误
	RC_Timeout           = -1 // 超时
	RC_Success           = 0  // 成功
	RC_User_STATUS_NOT   = 1  // 用户状态错误
	RC_Param_Error       = 2  // 参数错误
	RC_User_DB_Error     = 3  // 数据库错误
	RC_Config_Error      = 4  // 配置表错误
)
