package app

import (
	"flag"
	"fmt"
	"github.com/team-zf/framework/config"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/modules"
	"github.com/team-zf/framework/tables"
	"github.com/team-zf/framework/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type DefaultApp struct {
	PStatusTime               time.Duration // 打印状态的时间
	config                    *config.AppConfig
	debug                     bool
	parse                     bool
	confPath                  string
	logDir                    string
	tableDir                  string
	started                   bool
	modules                   []modules.IModule
	event_configurationLoaded func(app IApp)
	event_tablesLoaded        func(app IApp)
	event_startup             func(app IApp)
	event_stoped              func(app IApp)
}

func (e *DefaultApp) Init() {
	if e.parse {
		confPath := flag.String("c", "", "配置文件路径")
		logDir := flag.String("l", "", "日志文件目录")
		tableDir := flag.String("t", "", "数据表目录")
		flag.Parse()

		if confPath != nil && *confPath != "" {
			e.confPath = *confPath
		}
		if logDir != nil && *logDir != "" {
			e.logDir = *logDir
		}
		if tableDir != nil && *tableDir != "" {
			e.tableDir = *tableDir
		}
	}

	if e.logDir != "" {
		utils.Mkdir(e.logDir)
	}
}

func (e *DefaultApp) loadConfig() {
	if _, err := os.Open(e.confPath); err != nil {
		panic(fmt.Sprintf("未找到服务器配置文件 %s", e.confPath))
	}

	if conf, err := config.LoadConfig(e.confPath); err == nil {
		e.config = conf
	} else {
		panic(err)
	}

	logger.Init(e.debug, e.logDir, e.config.Logger)
	if e.event_configurationLoaded != nil {
		e.event_configurationLoaded(e)
	}
}

func (e *DefaultApp) loadTables() {
	if e.tableDir == "" {
		return
	}
	tables.LoadTables(e.tableDir, e.config.Table)
	if e.event_tablesLoaded != nil {
		e.event_tablesLoaded(e)
	}
}

func (e *DefaultApp) Run() {
	e.Init()
	e.loadConfig()
	e.loadTables()

	e.started = true
	for _, md := range e.modules {
		md.Start()
	}
	if e.event_startup != nil {
		e.event_startup(e)
	}

	// 这里要柱塞等关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	t := time.NewTicker(e.PStatusTime)
	defer t.Stop()
Pstatus:
	for {
		select {
		case <-c: //退出
			break Pstatus
		case <-t.C:
			var ps string
			for _, md := range e.modules {
				ps += md.PrintStatus()
			}
			logger.Notice(ps)
		}
	}
	e.started = false
	if e.event_stoped != nil {
		e.event_stoped(e)
	}
	for i := len(e.modules) - 1; i >= 0; i-- {
		md := e.modules[i]
		md.Stop()
	}
	logger.Close()
}

func (e *DefaultApp) AddModule(mds ...modules.IModule) {
	e.modules = append(e.modules, mds...)
	for _, md := range mds {
		md.Init()
	}
}

func (e *DefaultApp) OnConfigurationLoaded(fn func(app IApp)) {
	e.event_configurationLoaded = fn
}

func (e *DefaultApp) OnTablesLoaded(fn func(app IApp)) {
	e.event_tablesLoaded = fn
}

func (e *DefaultApp) OnStartup(fn func(app IApp)) {
	e.event_startup = fn
}

func (e *DefaultApp) OnStoped(fn func(app IApp)) {
	e.event_stoped = fn
}

func NewDefaultApp(opts ...Options) *DefaultApp {
	result := &DefaultApp{
		PStatusTime: 10 * time.Second,
		logDir:      "./logs",
		confPath:    "./server.json",
		tableDir:    "",
		modules:     make([]modules.IModule, 0),
		started:     false,
	}
	for _, opt := range opts {
		opt(result)
	}
	return result
}

func AppSetDebug(v bool) Options {
	return func(app IApp) {
		app.(*DefaultApp).debug = v
	}
}

func AppSetParse(v bool) Options {
	return func(app IApp) {
		app.(*DefaultApp).parse = v
	}
}

func AppSetLogDir(v string) Options {
	return func(app IApp) {
		app.(*DefaultApp).logDir = v
	}
}

func AppSetTableDir(v string) Options {
	return func(app IApp) {
		app.(*DefaultApp).tableDir = v
	}
}

func AppSetConfPath(v string) Options {
	return func(app IApp) {
		app.(*DefaultApp).confPath = v
	}
}

func AppSetPStatusTime(v time.Duration) Options {
	return func(app IApp) {
		app.(*DefaultApp).PStatusTime = v
	}
}
