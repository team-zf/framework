package Network

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/team-zf/framework/utils"
)

const (
	ROUTEHANDLE_HEADER uint32 = 0x00001100
	ROUTEHANDLE_MAXLEN uint16 = ^uint16(0)
)

type WebSocketRouteHandle struct {
	routes map[uint32]IWebSocketRoute
}

func (e *WebSocketRouteHandle) Marshal(data interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	buff, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	temp := make([]byte, 4)
	msglen := uint32(len(buff)+8) | ROUTEHANDLE_HEADER
	binary.LittleEndian.PutUint32(temp, msglen)
	buffer.Write(temp)
	buffer.Write(buff)
	return buffer.Bytes(), nil
}

func (e *WebSocketRouteHandle) Unmarshal(buff []byte) (interface{}, error) {
	msglen := binary.LittleEndian.Uint32(buff[:4]) ^ ROUTEHANDLE_HEADER
	if msglen != uint32(len(buff)) {
		return nil, fmt.Errorf("MsgLen Error: %d", msglen)
	}
	buff = buff[4:]

	// 从Body中取得JsonMap
	jsmap := make(map[string]interface{})
	err := json.Unmarshal(buff, &jsmap)
	if err != nil {
		return nil, err
	}

	// 从JsonMap中取得Cmd
	cmd, err := utils.NewStringAny(jsmap["cmd"]).ToUint32()
	if err != nil {
		return nil, err
	}

	route, err := e.GetRoute(cmd)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(buff, route)
	return route, nil
}

func (e *WebSocketRouteHandle) CheckMaxLenVaild(buff []byte) (uint32, bool) {
	msglen := binary.LittleEndian.Uint32(buff[:4]) ^ ROUTEHANDLE_HEADER
	if msglen > uint32(ROUTEHANDLE_MAXLEN) {
		return 0, false
	}
	if msglen > uint32(len(buff)) {
		return msglen, false
	}
	return msglen, true
}

func (e *WebSocketRouteHandle) GetRoute(cmd uint32) (IWebSocketRoute, error) {
	if route, ok := e.routes[cmd]; ok {
		newroute := utils.ReflectNew(route).(IWebSocketRoute)
		return newroute, nil
	} else {
		return nil, fmt.Errorf("Not Exist Cmd: %d.", cmd)
	}
}

func (e *WebSocketRouteHandle) SetRoute(cmd uint32, route IWebSocketRoute) {
	e.routes[cmd] = route
}

func NewWebSocketRouteHandle() *WebSocketRouteHandle {
	return &WebSocketRouteHandle{routes: make(map[uint32]IWebSocketRoute)}
}
