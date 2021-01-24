package Network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/messages"
	"github.com/team-zf/framework/modules"
	"github.com/team-zf/framework/utils"
	"github.com/team-zf/framework/utils/threads"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	WRITE_TIMEOUT     time.Duration = time.Minute
	HEARTBEAT_TIMEOUT time.Duration = time.Minute
)

type WebSocketModule struct {
	name         string
	addr         string
	httpServer   *http.Server
	routeHandle  *WebSocketRouteHandle
	thgo         *threads.ThreadGo // 协程管理器
	requestCount int64             // 收到的请求总数
	runingCount  int64             // 正在运行的总数
	onlineCount  int64             // 在线总人数
}

func (e *WebSocketModule) Init() {
	e.httpServer = &http.Server{
		Addr:         e.addr,
		WriteTimeout: WRITE_TIMEOUT,
	}
	handler := websocket.Handler(func(conn *websocket.Conn) {
		atomic.AddInt64(&e.onlineCount, 1)
		conn.PayloadType = websocket.BinaryFrame
		e.Handle(conn)
		atomic.AddInt64(&e.onlineCount, -1)
	})
	mux := http.NewServeMux()
	mux.Handle("/", handler)
	e.httpServer.Handler = mux
}

func (e *WebSocketModule) Start() {
	e.thgo.Go(func(ctx context.Context) {
		logger.Notice("%s启动", e.name)
		err := e.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error("%s异常关闭, 原因: %+v", err)
		}
	})
}

func (e *WebSocketModule) Stop() {
	e.httpServer.Close()
	e.thgo.CloseWait()
	logger.Notice("%s已停止", e.name)
}

func (e *WebSocketModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n\t\t%s的状态:\t%d/%d/%d\t(Online/Runing/Request)",
		e.name,
		atomic.LoadInt64(&e.onlineCount),
		atomic.LoadInt64(&e.runingCount),
		atomic.LoadInt64(&e.requestCount))
}

func (e *WebSocketModule) Handle(conn *websocket.Conn) {
	e.thgo.Wg.Add(1)
	defer e.thgo.Wg.Done()
	defer conn.Close()

	agent := new(WebSocketAgent)
	agent.Conn = conn

	// 心跳检测机制
	heartbeat := make(chan bool, 8)
	e.thgo.Go(func(ctx context.Context) {
		timeout := time.NewTimer(HEARTBEAT_TIMEOUT)
		defer timeout.Stop()
		defer conn.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timeout.C:
				return
			case reset := <-heartbeat:
				if reset {
					timeout.Reset(HEARTBEAT_TIMEOUT)
				} else {
					return
				}
			}
		}
	})

	// 消息接收
	e.thgo.Try(func(ctx context.Context) {
		buffer := &bytes.Buffer{}
		for {
			view := make([]byte, 1024*10)
			n, err := conn.Read(view)
			if err != nil {
				if err == io.EOF {
					heartbeat <- false
				}
				break
			}
			buffer.Write(view[:n])
			buff := buffer.Bytes()

			msglen, ok := e.routeHandle.CheckMaxLenVaild(buff)
			if ok { // 消息拼接完成
				buff = buffer.Next(int(msglen))
			} else if msglen > 0 { // 消息拼接未完成
				continue
			} else { // 异常消息的长度为0
				break
			}

			data, err := e.routeHandle.Unmarshal(buff)
			if err != nil {
				logger.Error("%s消息解码失败, 原因: %+v", e.name, err)
				return
			}

			route := data.(IWebSocketRoute)
			logger.Notice("%s收到请求: %s", e.name, route.Header())

			heartbeat <- true
			atomic.AddInt64(&e.requestCount, 1)
			atomic.AddInt64(&e.runingCount, 1)
			e.TryDirectCall(route, agent)
			atomic.AddInt64(&e.runingCount, -1)
		}
	}, func(err error) {
		// 无需处理
	})
}

func (e *WebSocketModule) TryDirectCall(route IWebSocketRoute, agent *WebSocketAgent) {
	utils.QueueRun(
		// 参数解析
		func() bool {
			result := true
			threads.Try(func() {
				route.Parse()
			}, func(err error) {
				result = false
				// 返回参数错误
				agent.SendData(&WebSocketResponse{
					Cmd:  route.GetCmd(),
					Code: messages.RC_Param_Error,
				})
			})
			return result
		},
		// 逻辑运行
		func() bool {
			result := true
			threads.Try(func() {
				code := route.Handle(agent)
				resp := &WebSocketResponse{
					Cmd:  route.GetCmd(),
					Code: code,
				}
				buff, _ := json.Marshal(resp)
				jsmap := make(map[string]interface{})
				json.Unmarshal(buff, &jsmap)
				for k, v := range route.ToJsonMap() {
					if _, ok := jsmap[k]; !ok {
						jsmap[k] = v
					}
				}
				agent.SendData(jsmap)
			}, func(err error) {
				result = false
				// 返回逻辑错误
				agent.SendData(&WebSocketResponse{
					Cmd:  route.GetCmd(),
					Code: messages.RC_LOGIC_ERROR,
				})
			})
			return result
		},
	)
}

func NewWebSocketModule(opts ...modules.ModOptions) *WebSocketModule {
	result := &WebSocketModule{
		name: "WebSocket",
		addr: ":8081",
		thgo: threads.NewThreadGo(),
	}

	for _, opt := range opts {
		opt(result)
	}
	return result
}

func WebSocketSetName(v string) modules.ModOptions {
	return func(mod modules.IModule) {
		mod.(*WebSocketModule).name = v
	}
}

func WebSocketSetAddr(v string) modules.ModOptions {
	return func(mod modules.IModule) {
		mod.(*WebSocketModule).addr = v
	}
}

func WebSocketSetRoute(v *WebSocketRouteHandle) modules.ModOptions {
	return func(mod modules.IModule) {
		mod.(*WebSocketModule).routeHandle = v
	}
}
