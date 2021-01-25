package modules

import (
	"github.com/team-zf/framework/config"
)

// IModule 模块接口
type IModule interface {
	// Init 初始化
	Init()
	// Start 启动
	Start()
	// Stop 停止
	Stop()
	// PrintStatus 打印状态
	PrintStatus() string
}

type IApp interface {
	Init() IApp
	Run(mds ...IModule) IApp
	AddModule(mds ...IModule) IApp
	OnConfigurationLoaded(fn func(app IApp, conf *config.AppConfig))
	OnTablesLoaded(fn func(app IApp))
	OnStartup(fn func(app IApp))
	OnStoped(fn func(app IApp))
	GetConfig() *config.AppConfig
	Debug() bool
}

type IWebSocketModule interface {
	IModule
	SetRoute(cmd uint32, msg interface{})
	GetRoute(cmd uint32) (msg interface{}, err error)
}

type AppOptions func(app IApp)

type ModOptions func(mod IModule)
