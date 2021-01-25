package Network

import (
	"context"
	"fmt"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/messages"
	"github.com/team-zf/framework/modules"
	"github.com/team-zf/framework/utils"
	"github.com/team-zf/framework/utils/threads"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"
)

type HttpModule struct {
	name         string
	ipPort       string
	httpServer   *http.Server
	routeHandle  *HttpRouteHandle
	thgo         *threads.ThreadGo
	timeout      time.Duration
	timeoutFun   func(IHttpRoute, http.ResponseWriter, *http.Request)
	requestCount int64 // 收到的请求总数
	runingCount  int64 // 正在运行的总数
}

func (e *HttpModule) Init() {
	e.httpServer = &http.Server{
		Addr:         e.ipPort,
		WriteTimeout: e.timeout,
	}
	// 还可以加别的参数，已后再加，有需要再加
	mux := http.NewServeMux()
	// 这个是主要的逻辑
	mux.HandleFunc("/", e.Handle)
	e.httpServer.Handler = mux
}

func (e *HttpModule) Start() {
	e.thgo.Go(func(ctx context.Context) {
		logger.Notice("%s启动", e.name)
		err := e.httpServer.ListenAndServe()
		if err != nil {
			if err != http.ErrServerClosed {
				logger.Error("Server closed unexpecteed; %v", err)
			}
		}
	})
}

func (e *HttpModule) Stop() {
	if err := e.httpServer.Close(); err != nil {
		logger.Error("Close Http Module; %v", err)
	}
	e.thgo.CloseWait()
	logger.Notice("%s已停止", e.name)
}

func (e *HttpModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n\t\t%s的状态:\t%d/%d/%d\t(Runing/Request)",
		e.name,
		atomic.LoadInt64(&e.runingCount),
		atomic.LoadInt64(&e.requestCount))
}

func (e *HttpModule) Handle(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")

	e.thgo.Wg.Add(1)
	defer e.thgo.Wg.Done()

	bytes, _ := ioutil.ReadAll(req.Body)
	msg, err := e.routeHandle.Unmarshal(bytes)
	if err != nil {
		logger.Error("%s消息解码失败, 原因: %+v", e.name, err)
		return
	}
	route, _ := msg.(IHttpRoute)
	logger.Notice("%s收到请求: %s", e.name, route.Header())

	atomic.AddInt64(&e.requestCount, 1)
	atomic.AddInt64(&e.runingCount, 1)
	e.TryDirectCall(route, res, req)
	atomic.AddInt64(&e.runingCount, -1)
}

func (e *HttpModule) TryDirectCall(route IHttpRoute, res http.ResponseWriter, req *http.Request) {
	utils.QueueRun(
		func() bool {
			result := true
			threads.Try(func() {
				route.Parse()
			}, func(err error) {
				result = false
				resp := &HttpResponse{
					Code: messages.RC_Param_Error,
				}
				if bytes, err := e.routeHandle.Marshal(resp); err == nil {
					res.Write(bytes)
				}
			})
			return result
		},
		func() bool {
			result := true
			threads.Try(
				func() {
					t := time.NewTimer(e.timeout - 2*time.Second)
					g := threads.NewGoRun(func() {
						threads.Try(
							func() {
								code := route.Handle(req)
								resp := &HttpResponse{
									Code: code,
								}
								if data, err := e.routeHandle.Marshal(resp); err == nil {
									res.Write(data)
								}
							},
							func(err error) {
								logger.Error("%s; 逻辑报错: %+v", route.Header(), err)
								resp := &HttpResponse{
									Code: messages.RC_LOGIC_ERROR,
								}
								// 返回逻辑错误
								if data, err := e.routeHandle.Marshal(resp); err == nil {
									res.Write(data)
								}
							},
						)
					})
					select {
					// 业务逻辑完成
					case <-g.Chanresult:
						t.Stop()
						break
					// 业务逻辑超时
					case <-t.C:
						if e.timeoutFun != nil {
							e.timeoutFun(route, res, req)
						} else {
							e.defaultTimeoutFunc(route, res, req)
						}
						break
					}
				},
				func(err error) {
					result = false
					resp := &HttpResponse{
						Code: messages.RC_Param_Error,
					}
					if bytes, err := e.routeHandle.Marshal(resp); err == nil {
						res.Write(bytes)
					}
				},
			)
			return result
		},
	)
}

func (e *HttpModule) defaultTimeoutFunc(route IHttpRoute, res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("timeout Run!"))
}

func NewHttpModule(opts ...modules.ModOptions) *HttpModule {
	result := &HttpModule{
		name:        "Http",
		ipPort:      ":8080",
		timeout:     30 * time.Second,
		thgo:        threads.NewThreadGo(),
		routeHandle: NewHttpRouteHandle(),
	}
	for _, opt := range opts {
		opt(result)
	}
	return result
}

//设置名称
func HttpSetName(v string) modules.ModOptions {
	return func(mod modules.IModule) {
		mod.(*HttpModule).name = v
	}
}

// 设置Web地址
func HttpSetIpPort(ipPort string) modules.ModOptions {
	return func(mod modules.IModule) {
		mod.(*HttpModule).ipPort = ipPort
	}
}

// 设置超时时间
func HttpSetTimeout(timeout time.Duration) modules.ModOptions {
	return func(mod modules.IModule) {
		mod.(*HttpModule).timeout = timeout * time.Second
	}
}

// 设置超时回调方法
func HttpSetTimeoutFunc(timeoutfunc func(IHttpRoute, http.ResponseWriter, *http.Request)) modules.ModOptions {
	return func(mod modules.IModule) {
		mod.(*HttpModule).timeoutFun = timeoutfunc
	}
}

// 设置路由
func HttpSetRoute(route *HttpRouteHandle) modules.ModOptions {
	return func(mod modules.IModule) {
		mod.(*HttpModule).routeHandle = route
	}
}
