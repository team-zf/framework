package main

import (
	"github.com/team-zf/framework/Route"
	"github.com/team-zf/framework/app"
	"github.com/team-zf/framework/modules"
	"time"
)

func main() {
	app := CreateApp(
		app.AppSetDebug(true),
		app.AppSetParse(true),
		app.AppSetTableDir("./JSON"),
		app.AppSetPStatusTime(3*time.Second),
	)
	app.AddModule(modules.NewHttpModule(
		modules.HttpSetRoute(Route.HttpRoute),
	))
	app.AddModule(modules.NewWebSocketModule())
	app.Run()
}

func CreateApp(opts ...app.Options) app.IApp {
	return app.NewDefaultApp(opts...)
}
