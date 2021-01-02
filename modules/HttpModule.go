package modules

import (
	"context"
	"fmt"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/messages"
	"github.com/team-zf/framework/utils"
	"github.com/team-zf/framework/utils/threads"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"
)

type HttpModule struct {
	ipPort      string                                                                // 监听地址
	timeout     time.Duration                                                         // 超时时长
	timeoutFun  func(messages.IHttpMessageHandle, http.ResponseWriter, *http.Request) // 超时回调，把超时的消息传入
	RouteHandle messages.IMessageHandle                                               // 消息路由
	getnum      int64                                                                 // 收到的总消息数
	runing      int64                                                                 // 当前在处理的消息数
	httpServer  *http.Server                                                          // HTTP请求的对象
	thgo        *threads.ThreadGo                                                     // 协程管理器
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
		logger.Notice("Web Module Start.")
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
		logger.Error("Close Web Module; %v", err)
	}
	e.thgo.CloseWait()
	logger.Notice("Web Module Stop.")
}

func (e *HttpModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n\t\thttp Module\t:%d/%d\t(get/runing)",
		atomic.LoadInt64(&e.getnum),
		atomic.LoadInt64(&e.runing))
}

func (e *HttpModule) Handle(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")

	e.thgo.Wg.Add(1)
	defer e.thgo.Wg.Done()
	atomic.AddInt64(&e.getnum, 1)
	atomic.AddInt64(&e.runing, 1)
	defer atomic.AddInt64(&e.runing, -1)

	bytes, _ := ioutil.ReadAll(req.Body)
	msg, err := e.RouteHandle.Unmarshal(bytes)
	if err != nil {
		logger.Warn("http RouteHandle Unmarshal Error: %v", err)
		return
	}
	handle, ok := msg.(messages.IHttpMessageHandle)
	if !ok {
		logger.Notice("Not is http Msg: %+v", msg)
		return
	} else {
		logger.Notice("http Get Msg: %s", handle.Header())
	}

	utils.QueueRun(
		// 检测参数
		e.queueParse(handle, res, req),
		// 处理回调
		e.queueCall(handle, res, req),
	)
}

func (e *HttpModule) queueParse(handle messages.IHttpMessageHandle, res http.ResponseWriter, req *http.Request) func() bool {
	return func() bool {
		result := true
		threads.Try(
			func() {
				handle.Parse()
			},
			func(err error) {
				result = false
				resp := messages.NewHttpResponse(
					messages.ResponseSetCode(messages.RC_Param_Error),
				)
				if bytes, err := e.RouteHandle.Marshal(resp); err == nil {
					res.Write(bytes)
				}
			},
		)
		return result
	}
}

func (e *HttpModule) queueCall(handle messages.IHttpMessageHandle, res http.ResponseWriter, req *http.Request) func() bool {
	return func() bool {
		result := true
		threads.Try(
			func() {
				t := time.NewTimer(e.timeout - 2*time.Second)
				g := threads.NewGoRun(func() {
					resp := messages.NewHttpResponse()
					threads.Try(
						func() {
							handle.HttpDirectCall(req, resp)
						},
						func(err error) {
							logger.Error("%s; Logic Error: %+v", handle.Header(), err)
							resp.Code = messages.RC_LOGIC_ERROR
						},
						func() {
							if bytes, err := e.RouteHandle.Marshal(resp); err == nil {
								res.Write(bytes)
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
					logger.Debug("http Timeout msg: %+v", handle.GetCmd())
					if e.timeoutFun != nil {
						e.timeoutFun(handle, res, req)
					} else {
						e.defaultTimeoutFunc(handle, res, req)
					}
					break
				}
			},
			func(err error) {
				result = false
				logger.Error("HttpModule queueCall Error: %+v", err)
				// 如果出异常了，跑这里
				res.Write([]byte("catch!"))
			},
		)
		return result
	}
}

func (e *HttpModule) defaultTimeoutFunc(msg messages.IHttpMessageHandle, res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("timeout Run!"))
}

func NewHttpModule(opts ...Options) *HttpModule {
	result := &HttpModule{
		ipPort:      ":8080",
		timeout:     30 * time.Second,
		getnum:      0,
		runing:      0,
		thgo:        threads.NewThreadGo(),
		RouteHandle: messages.NewHttpMessageHandle(),
	}
	for _, opt := range opts {
		opt(result)
	}
	return result
}

// 设置Web地址
func HttpSetIpPort(ipPort string) Options {
	return func(mod IModule) {
		mod.(*HttpModule).ipPort = ipPort
	}
}

// 设置超时时间
func HttpSetTimeout(timeout time.Duration) Options {
	return func(mod IModule) {
		mod.(*HttpModule).timeout = timeout * time.Second
	}
}

// 设置超时回调方法
func HttpSetTimeoutFunc(timeoutfunc func(messages.IHttpMessageHandle, http.ResponseWriter, *http.Request)) Options {
	return func(mod IModule) {
		mod.(*HttpModule).timeoutFun = timeoutfunc
	}
}

// 设置路由
func HttpSetRoute(route messages.IMessageHandle) Options {
	return func(mod IModule) {
		mod.(*HttpModule).RouteHandle = route
	}
}
