package modules

import (
	"bytes"
	"context"
	"fmt"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/messages"
	"github.com/team-zf/framework/model"
	"github.com/team-zf/framework/utils"
	"github.com/team-zf/framework/utils/threads"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

type WebSocketModule struct {
	ipPort      string                  // HTTP监听的地址
	timeout     time.Duration           // 超时时长
	RouteHandle messages.IMessageHandle // 消息路由
	getnum      int64                   // 收到的总消息数
	runing      int64                   // 当前在处理的消息数
	connlen     int64                   // 连接数
	httpServer  *http.Server            // HTTP请求的对象
	thgo        *threads.ThreadGo       // 协程管理器
	frame       byte                    // websocket PayloadType

	event_WebSocketOnline func(wsmd *model.WebSocketModel) // 连接成功事件, 可以用来获取一些连接的信息; 比如IP
}

func (e *WebSocketModule) Init() {
	e.httpServer = &http.Server{
		Addr:         e.ipPort,
		WriteTimeout: e.timeout,
	}
	// 还可以加别的参数，已后再加，有需要再加
	mux := http.NewServeMux()
	// 这个是主要的逻辑
	mux.Handle("/", websocket.Handler(e.Handle))
	e.httpServer.Handler = mux
}

func (e *WebSocketModule) Start() {
	e.thgo.Go(func(ctx context.Context) {
		logger.Notice("WebSocket Module Start.")
		err := e.httpServer.ListenAndServe()
		if err != nil {
			if err != http.ErrServerClosed {
				logger.Error("Server closed unexpecteed; %v", err)
			}
		}
	})
}

func (e *WebSocketModule) Stop() {
	if err := e.httpServer.Close(); err != nil {
		logger.Error("Close WebSocket Module; %v", err)
	}
	e.thgo.CloseWait()
	logger.Notice("WebSocket Module Stop.")
}

func (e *WebSocketModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n\t\tWebSocket Module\t:%d/%d/%d\t(connum/getmsg/runing)",
		atomic.LoadInt64(&e.connlen),
		atomic.LoadInt64(&e.getnum),
		atomic.LoadInt64(&e.runing))
}

func (e *WebSocketModule) Handle(conn *websocket.Conn) {
	conn.PayloadType = e.frame

	e.thgo.Wg.Add(1)
	defer e.thgo.Wg.Done()
	defer conn.Close()

	// 发给下面的连接对象，可以自定义一些信息和回调
	wsmd := new(model.WebSocketModel)
	wsmd.Conn = conn
	wsmd.KeyID = -1
	if e.event_WebSocketOnline != nil {
		e.event_WebSocketOnline(wsmd)
	}
	atomic.AddInt64(&e.connlen, 1)
	logger.Info("WebSocket Client Open: %+v .", wsmd.KeyID, wsmd.ConInfo)

	// 发消息来说明这个用户掉线了
	defer func() {
		atomic.AddInt64(&e.connlen, -1)
		logger.Info("WebSocket Client Closeing: %+v .", wsmd.KeyID, wsmd.ConInfo)
		// 用来处理发生连接关闭的时候，要处理的事
		if wsmd.CloseFun != nil {
			wsmd.CloseFun(wsmd)
		}
		logger.Info("WebSocket Client Close: %+v .", wsmd.KeyID, wsmd.ConInfo)
	}()

	// 心跳
	runchan := make(chan bool, 8)
	e.thgo.Go(
		func(ctx context.Context) {
			timeout := time.NewTimer(e.timeout)
			defer timeout.Stop()
			defer conn.Close()
			for {
				select {
				case <-ctx.Done():
					return
				case <-timeout.C:
					return
				case ok := <-runchan:
					if ok {
						timeout.Reset(e.timeout)
					} else {
						return
					}
				}
			}
		},
	)

	e.thgo.Try(
		func(ctx context.Context) {
			buf := &bytes.Buffer{}
			for {
				rdbuff := make([]byte, 10240)
				n, err := conn.Read(rdbuff)
				if err != nil {
					if err == io.EOF {
						runchan <- false
					}
					break
				}
				buf.Write(rdbuff[:n])
				buff := buf.Bytes()
				if msglen, ok := e.RouteHandle.CheckMaxLenVaild(buff); ok {
					buff = buf.Next(int(msglen))
				} else if msglen == 0 { // 消息长度异常
					break
				} else {
					continue
				}

				msg, err := e.RouteHandle.Unmarshal(buff)
				if err != nil {
					logger.Info("WebSocket RouteHandle Unmarshal Error: %s", err.Error())
					return
				}
				handle, ok := msg.(messages.IWebSocketMessageHandle)
				if !ok {
					logger.Info("Not is WebSocket Msg: %+v", msg)
					return
				} else {
					logger.Info("WebSocket Get Msg: %+v", msg)
				}

				runchan <- true
				atomic.AddInt64(&e.getnum, 1)
				atomic.AddInt64(&e.runing, 1)
				utils.QueueRun(
					// 检测参数
					e.queueParse(handle, wsmd),
					// 处理回调
					e.queueCall(handle, wsmd),
				)
				atomic.AddInt64(&e.runing, -1)
			}
		},
		func(err error) {},
	)
}

func (e *WebSocketModule) queueParse(handle messages.IWebSocketMessageHandle, wsmd *model.WebSocketModel) func() bool {
	return func() bool {
		result := true
		threads.Try(
			func() {
				handle.Parse()
			},
			func(err error) {
				result = false
				resp := &messages.WebSocketResponse{
					Code: messages.RC_Param_Error,
				}
				if data, err := e.RouteHandle.Marshal(resp); err == nil {
					wsmd.SendByte(data)
				}
			},
		)
		return result
	}
}

func (e *WebSocketModule) queueCall(handle messages.IWebSocketMessageHandle, wsmd *model.WebSocketModel) func() bool {
	return func() bool {
		result := true
		threads.Try(
			func() {
				t := time.NewTimer(5 * time.Second)
				g := threads.NewGoRun(func() {
					resp := &messages.WebSocketResponse{}
					threads.Try(
						func() {
							handle.WebSocketDirectCall(wsmd, resp)
						},
						func(err error) {
							logger.Error("%s; 逻辑报错: %+v", handle.Header(), err)
							resp.Code = messages.RC_LOGIC_ERROR
						},
						func() {
							if data, err := e.RouteHandle.Marshal(resp); err == nil {
								wsmd.SendByte(data)
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
					logger.Debug("逻辑超时: %s", handle.Header())
					break
				}
			},
			func(err error) {
				result = false
				logger.Error("HttpModule queueCall Error: %+v", err)
				// 如果出异常了，跑这里
				wsmd.SendByte([]byte("catch!"))
			},
		)
		return result
	}
}

func NewWebSocketModule(opts ...ModOptions) *WebSocketModule {
	result := &WebSocketModule{
		ipPort:      ":8081",
		timeout:     60 * time.Second,
		getnum:      0,
		runing:      0,
		connlen:     0,
		thgo:        threads.NewThreadGo(),
		RouteHandle: messages.NewHttpMessageHandle(),
		frame:       websocket.BinaryFrame, //因为我们用的是路由是二进制的方式，所以这里要用这个值
	}

	for _, opt := range opts {
		opt(result)
	}
	return result
}

func WebSocketSetIpPort(v string) ModOptions {
	return func(mod IModule) {
		mod.(*WebSocketModule).ipPort = v
	}
}

func WebSocketSetTimeout(v time.Duration) ModOptions {
	return func(mod IModule) {
		mod.(*WebSocketModule).timeout = v
	}
}

func WebSocketSetRoute(v messages.IMessageHandle) ModOptions {
	return func(mod IModule) {
		mod.(*WebSocketModule).RouteHandle = v
	}
}
