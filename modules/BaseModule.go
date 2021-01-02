package modules

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
	Run()
	AddModule(mds ...IModule)
	OnConfigurationLoaded(fn func(app IApp))
	OnTablesLoaded(fn func(app IApp))
	OnStartup(fn func(app IApp))
	OnStoped(fn func(app IApp))
}

type AppOptions func(app IApp)

type ModOptions func(mod IModule)
