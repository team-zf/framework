package main

import (
	"github.com/team-zf/framework/Route"
	"github.com/team-zf/framework/modules"
	"time"
)

func main() {
	app := CreateApp(
		modules.AppSetDebug(true),
		modules.AppSetParse(true),
		modules.AppSetTableDir("./JSON"),
		modules.AppSetPStatusTime(3*time.Second),
	)
	app.AddModule(modules.NewHttpModule(
		modules.HttpSetRoute(Route.HttpRoute),
	))
	app.AddModule(modules.NewWebSocketModule())
	app.Run()
}

func CreateApp(opts ...modules.AppOptions) modules.IApp {
	return modules.NewApp(opts...)
}
