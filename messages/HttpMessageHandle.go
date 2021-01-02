package messages

import (
	"encoding/json"
	"fmt"
	"github.com/team-zf/framework/utils"
)

type HttpMessageHandle struct {
	routes map[uint32]interface{}
}

func NewHttpMessageHandle() *HttpMessageHandle {
	return &HttpMessageHandle{
		routes: make(map[uint32]interface{}),
	}
}

func (e *HttpMessageHandle) Marshal(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (e *HttpMessageHandle) Unmarshal(buff []byte) (data interface{}, err error) {
	jsmap := make(map[string]interface{})
	if err = json.Unmarshal(buff, &jsmap); err != nil {
		return
	}
	var cmd uint32
	cmd, err = utils.NewStringAny(jsmap["cmd"]).ToUint32()
	if err != nil {
		return
	}
	data, err = e.GetRoute(cmd)
	if err != nil {
		return
	}
	err = json.Unmarshal(buff, data)
	return
}

func (e *HttpMessageHandle) SetRoute(cmd uint32, msg interface{}) {
	e.routes[cmd] = msg
}

func (e *HttpMessageHandle) GetRoute(cmd uint32) (msg interface{}, err error) {
	if msget, ok := e.routes[cmd]; ok {
		msg = utils.ReflectNew(msget)
	} else {
		err = fmt.Errorf("Not Exist Cmd: %d.", cmd)
	}
	return
}

func (e *HttpMessageHandle) CheckMaxLenVaild(buff []byte) (msglen uint32, ok bool) {
	return uint32(len(buff)), true
}
