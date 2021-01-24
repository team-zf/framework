package Network

type IWebSocketRule interface {
	GetName() string
	GetPrimaryKey() string
	GetDirectKey() string
}
