package messages

/**
 * 服务器响应代码
 * 注意: Http和WebSocket通用
 */

const (
	RC_NsqRequestFailure uint32 = 100 // 请求失败
	RC_NotResult         uint32 = 101 // 不回复
	RC_Timeout           uint32 = 102 // 超时
)

const (
	RC_Success         uint32 = 200 // 成功
	RC_User_STATUS_NOT uint32 = 201 // 用户状态错误
)

const (
	RC_Param_Error uint32 = 300 // 参数错误
)

const (
	RC_NotLogic uint32 = 400 // 没有逻辑处理它
	RC_NotCmd   uint32 = 404 // 没有事件处理这个消息
)

const (
	RC_LOGIC_ERROR   uint32 = 500 // 逻辑处理错误
	RC_User_DB_Error uint32 = 501 // 数据库错误
	RC_Config_Error  uint32 = 502 // 配置表错误
)
