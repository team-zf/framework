package modules

import (
	"flag"
	"fmt"
	"github.com/team-zf/framework/config"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/tables"
	"github.com/team-zf/framework/utils"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type App struct {
	PStatusTime               time.Duration // 打印状态的时间
	config                    *config.AppConfig
	debug                     bool
	parse                     bool
	confPath                  string
	logDir                    string
	tableDir                  string
	started                   bool
	modules                   []IModule
	event_ConfigurationLoaded func(app IApp, conf *config.AppConfig)
	event_TablesLoaded        func(app IApp)
	event_Startup             func(app IApp)
	event_Stoped              func(app IApp)
}

func (e *App) Init() IApp {
	if e.parse {
		confPath := flag.String("c", "", "配置文件路径")
		logDir := flag.String("l", "", "日志文件目录")
		tableDir := flag.String("t", "", "数据表目录")
		debug := flag.String("d", "", "是否启动调试模式")
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
		if debug != nil && *debug != "" {
			e.debug = strings.ToLower(*debug) == "true"
		}
	}

	if e.logDir != "" {
		utils.Mkdir(e.logDir)
	}
	e.loadConfig()
	//e.loadTables()
	return e
}

func (e *App) loadConfig() {
	if _, err := os.Open(e.confPath); err != nil {
		panic(fmt.Sprintf("未找到服务器配置文件 %s", e.confPath))
	}

	if conf, err := config.LoadConfig(e.confPath); err == nil {
		e.config = conf
	} else {
		panic(err)
	}

	logger.Init(e.debug, e.logDir, e.config.Logger)
	if e.event_ConfigurationLoaded != nil {
		e.event_ConfigurationLoaded(e, e.config)
	}
}

func (e *App) loadTables() {
	if e.tableDir == "" {
		return
	}
	tables.LoadTables(e.tableDir, e.config.Table)
	if e.event_TablesLoaded != nil {
		e.event_TablesLoaded(e)
	}
}

func (e *App) Run(mds ...IModule) IApp {
	if len(mds) > 0 {
		e.AddModule(mds...)
	}

	e.started = true
	for _, md := range e.modules {
		md.Start()
	}
	if e.event_Startup != nil {
		e.event_Startup(e)
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
	if e.event_Stoped != nil {
		e.event_Stoped(e)
	}
	for i := len(e.modules) - 1; i >= 0; i-- {
		md := e.modules[i]
		md.Stop()
	}
	logger.Close()
	return e
}

func (e *App) AddModule(mds ...IModule) IApp {
	e.modules = append(e.modules, mds...)
	for _, md := range mds {
		md.Init()
	}
	return e
}

func (e *App) OnConfigurationLoaded(fn func(app IApp, conf *config.AppConfig)) {
	e.event_ConfigurationLoaded = fn
}

func (e *App) OnTablesLoaded(fn func(app IApp)) {
	e.event_TablesLoaded = fn
}

func (e *App) OnStartup(fn func(app IApp)) {
	e.event_Startup = fn
}

func (e *App) OnStoped(fn func(app IApp)) {
	e.event_Stoped = fn
}

func (e *App) GetConfig() *config.AppConfig {
	return e.config
}

func (e *App) Debug() bool {
	return e.debug
}

func NewApp(opts ...AppOptions) *App {
	result := &App{
		PStatusTime: 10 * time.Second,
		logDir:      "./logs",
		confPath:    "./server.json",
		tableDir:    "",
		modules:     make([]IModule, 0),
		started:     false,
	}
	for _, opt := range opts {
		opt(result)
	}
	return result
}

func AppSetDebug(v bool) AppOptions {
	return func(app IApp) {
		app.(*App).debug = v
	}
}

func AppSetParse(v bool) AppOptions {
	return func(app IApp) {
		app.(*App).parse = v
	}
}

func AppSetLogDir(v string) AppOptions {
	return func(app IApp) {
		app.(*App).logDir = v
	}
}

func AppSetTableDir(v string) AppOptions {
	return func(app IApp) {
		app.(*App).tableDir = v
	}
}

func AppSetConfPath(v string) AppOptions {
	return func(app IApp) {
		app.(*App).confPath = v
	}
}

func AppSetPStatusTime(v time.Duration) AppOptions {
	return func(app IApp) {
		app.(*App).PStatusTime = v
	}
}
