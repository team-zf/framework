package messages

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/team-zf/framework/utils"
)

const (
	messageHeader uint32 = 0x00001100
	messageMaxLen uint16 = ^uint16(0)
)

type WebSocketMessageHandle struct {
	routes map[uint32]interface{}
}

func (e *WebSocketMessageHandle) Marshal(data interface{}) ([]byte, error) {
	buff := &bytes.Buffer{}
	in_data, err := json.Marshal(data)
	tempbuff := make([]byte, 4)
	pklen := uint32(len(in_data)+8) | messageHeader
	binary.LittleEndian.PutUint32(tempbuff, pklen)
	buff.Write(tempbuff)
	buff.Write(in_data)
	return buff.Bytes(), err
}

func (e *WebSocketMessageHandle) Unmarshal(buff []byte) (data interface{}, err error) {
	pklen := binary.LittleEndian.Uint32(buff[:4])
	pklen = pklen ^ messageHeader
	if pklen != uint32(len(buff)) {
		err = fmt.Errorf("MsgLen Error: %d", pklen)
		return
	}
	buff = buff[4:]

	// 从Body中取得JsonMap
	jsmap := make(map[string]interface{})
	err = json.Unmarshal(buff, &jsmap)
	if err != nil {
		return
	}

	// 从JsonMap中取得Cmd
	var cmd uint32
	cmd, err = utils.NewStringAny(jsmap["cmd"]).ToUint32()
	if err != nil {
		return
	}

	// 通过Cmd取得对应的路由
	data, err = e.GetRoute(cmd)
	if err != nil {
		return
	}

	// 将Body转为路由格式
	err = json.Unmarshal(buff, data)
	return
}

func (e *WebSocketMessageHandle) SetRoute(cmd uint32, msg interface{}) {
	e.routes[cmd] = msg
}

func (e *WebSocketMessageHandle) GetRoute(cmd uint32) (msg interface{}, err error) {
	if msget, ok := e.routes[cmd]; ok {
		msg = utils.ReflectNew(msget)
	} else {
		err = fmt.Errorf("Not Exist Cmd: %d.", cmd)
	}
	return
}

func (e *WebSocketMessageHandle) CheckMaxLenVaild(buff []byte) (msglen uint32, ok bool) {
	pklen := binary.LittleEndian.Uint32(buff[:4])
	pklen = pklen ^ messageHeader
	if pklen > uint32(messageMaxLen) {
		return 0, false
	}
	if pklen > uint32(len(buff)) {
		return pklen, false
	}
	return pklen, true
}

func NewWebSocketMessageHandle() *WebSocketMessageHandle {
	return &WebSocketMessageHandle{
		routes: make(map[uint32]interface{}),
	}
}
