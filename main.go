package main

import (
	"github.com/team-zf/framework/Control"
	"github.com/team-zf/framework/Route"
	"github.com/team-zf/framework/config"
	"github.com/team-zf/framework/modules"
	"time"
)

func init() {
}

func main() {
	Control.App = CreateApp(
		modules.AppSetDebug(true),
		modules.AppSetParse(true),
		modules.AppSetTableDir("./JSON"),
		modules.AppSetPStatusTime(3*time.Second),
	)
	Control.App.OnConfigurationLoaded(func(app modules.IApp, conf *config.AppConfig) {
		// 载入数据库模块
		Control.DB = modules.NewDataBaseModule(
			modules.DataBaseSetConf(conf.MySql),
		)
		app.AddModule(Control.DB)
	})
	Control.App.Init()

	Control.App.Run(
		modules.NewHttpModule(
			modules.HttpSetRoute(Route.HttpRoute),
		),
		modules.NewWebSocketModule(),
	)
}

func CreateApp(opts ...modules.AppOptions) modules.IApp {
	return modules.NewApp(opts...)
}
