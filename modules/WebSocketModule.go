package modules

import (
	"bytes"
	"context"
	"fmt"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/messages"
	"github.com/team-zf/framework/model"
	"github.com/team-zf/framework/utils/threads"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

type WebSocketModule struct {
	ipPort             string                           // HTTP监听的地址
	timeout            time.Duration                    // 超时时长
	RouteHandle        messages.IMessageHandle          // 消息路由
	webSocketOnlineFun func(conn *model.WebSocketModel) // 连接成功后回调，可以用来获取一些连接的信息，比如IP
	getnum             int64                            // 收到的总消息数
	runing             int64                            // 当前在处理的消息数
	connlen            int64                            // 连接数
	httpServer         *http.Server                     // HTTP请求的对象
	thgo               *threads.ThreadGo                // 协程管理器
	frame              byte                             // websocket PayloadType
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
	ws := new(model.WebSocketModel)
	ws.Conn = conn
	ws.KeyID = -1
	if e.webSocketOnlineFun != nil {
		e.webSocketOnlineFun(ws)
	}
	atomic.AddInt64(&e.connlen, 1)
	logger.Info("WebSocket Client Open: %+v .", ws.KeyID, ws.ConInfo)

	// 发消息来说明这个用户掉线了
	defer func() {
		atomic.AddInt64(&e.connlen, -1)
		logger.Info("WebSocket Client Closeing: %+v .", ws.KeyID, ws.ConInfo)
		// 用来处理发生连接关闭的时候，要处理的事
		if ws.CloseFun != nil {
			ws.CloseFun(ws)
		}
		logger.Info("WebSocket Client Close: %+v .", ws.KeyID, ws.ConInfo)
	}()

	// 用来处理超时
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
		listen:
			for {
				rdbuff := make([]byte, 10240)
				n, err := conn.Read(rdbuff)
				if err != nil {
					if err == io.EOF {
						runchan <- false
					}
					break listen
				}
				buf.Write(rdbuff[:n])
				buff := buf.Bytes()
				if msglen, ok := e.RouteHandle.CheckMaxLenVaild(buff); ok {
					buff = buf.Next(int(msglen))
				} else {
					if msglen == 0 {
						// 消息长度异常
						break listen
					}
					continue
				}

				msg, err := e.RouteHandle.Unmarshal(buff)
				if err != nil {
					logger.Info("WebSocket RouteHandle Unmarshal Error: %s", err.Error())
					return
				}
				modmsg, ok := msg.(messages.IWebSocketMessageHandle)
				if !ok {
					logger.Info("Not is WebSocket Msg: %+v", msg)
					return
				} else {
					logger.Info("WebSocket Get Msg: %+v", msg)
				}

				runchan <- true
				atomic.AddInt64(&e.getnum, 1)
				e.thgo.Try(func(ctx context.Context) {
					atomic.AddInt64(&e.runing, 1)
					modmsg.WebSocketDirectCall(ws)
				}, nil, func() {
					atomic.AddInt64(&e.runing, -1)
				})

			}
		}, nil,
	)
}

func NewWebSocketModule(opts ...ModOptions) *WebSocketModule {
	result := &WebSocketModule{
		ipPort:             ":8081",
		timeout:            60 * time.Second,
		getnum:             0,
		runing:             0,
		connlen:            0,
		thgo:               threads.NewThreadGo(),
		RouteHandle:        messages.NewHttpMessageHandle(),
		webSocketOnlineFun: nil,
		frame:              websocket.BinaryFrame, //因为我们用的是路由是二进制的方式，所以这里要用这个值
	}

	for _, opt := range opts {
		opt(result)
	}
	return result
}
