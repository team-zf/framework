package app

import "github.com/team-zf/framework/modules"

type IApp interface {
	Run()
	AddModule(mds ...modules.IModule)
	OnConfigurationLoaded(fn func(app IApp))
	OnTablesLoaded(fn func(app IApp))
	OnStartup(fn func(app IApp))
	OnStoped(fn func(app IApp))
}

type Options func(app IApp)
