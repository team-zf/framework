package Network

import "github.com/team-zf/framework/utils"

type WebSocketDDM struct {
	_Rules map[string][]string
	_Data  map[string]interface{}
	_Del   map[string]interface{}
	_Mod   map[string]interface{}
}

func (e *WebSocketDDM) Data(key string, val interface{}) {
	if e._Data == nil {
		e._Data = make(map[string]interface{})
	}
	e._Data[key] = val
}

func (e *WebSocketDDM) Del(rule []string, args ...map[string]interface{}) {
	if e._Del == nil {
		e._Del = make(map[string]interface{})
	}

	key := rule[0]
	pk := rule[1]

	for _, arg := range args {
		if e._Del[key] == nil {
			e._Del[key] = make([]int, 0)
		}
		id := utils.NewStringAny(arg[pk]).ToIntV()
		e._Del[key] = append(e._Del[key].([]int), id)
	}

	if e._Rules == nil {
		e._Rules = make(map[string][]string)
	}
	if _, ok := e._Rules[key]; !ok {
		e._Rules[key] = rule
	}
}

func (e *WebSocketDDM) Mod(rule []string, args ...map[string]interface{}) {
	if e._Mod == nil {
		e._Mod = make(map[string]interface{})
	}

	key := rule[0]
	pk := rule[1]
	dk := rule[2]

	for _, arg := range args {
		if e._Mod[key] == nil {
			e._Mod[key] = make(map[string]interface{})
		}
		if pk != "" {
			m := e._Mod[key].(map[string]interface{})
			k := utils.NewStringAny(arg[pk]).ToString()
			if dk != "" {
				m[k] = arg[dk]
			} else {
				m[k] = arg
			}
		} else {
			e._Mod[key] = arg
		}
	}

	if e._Rules == nil {
		e._Rules = make(map[string][]string)
	}
	if _, ok := e._Rules[key]; !ok {
		e._Rules[key] = rule
	}
}

func (e *WebSocketDDM) Join(args ...*WebSocketDDM) *WebSocketDDM {
	for _, arg := range args {
		e.__JoinMod(e, arg)
		e.__JoinDel(e, arg)
	}
	return e
}

func (e *WebSocketDDM) __JoinMod(a *WebSocketDDM, b *WebSocketDDM) {
	for key, m := range b._Mod {
		rule := b._Rules[key]

		if a._Mod[key] == nil {
			a._Mod[key] = make(map[string]interface{})
		}

		if rule[1] != "" {
			ajm := a._Mod[key].(map[string]interface{})
			bjm := m.(map[string]interface{})
			for k, v := range bjm {
				ajm[k] = v
			}
		} else {
			a._Mod[key] = m
		}
	}
}

func (e *WebSocketDDM) __JoinDel(a *WebSocketDDM, b *WebSocketDDM) {
	for key, m := range b._Del {
		if a._Del[key] == nil {
			a._Del[key] = make([]int, 0)
		}

		for _, v := range m.([]int) {
			a._Del[key] = append(a._Del[key].([]int), v)
		}
	}
}

func (e *WebSocketDDM) ToJsonMap() map[string]interface{} {
	m := make(map[string]interface{})
	if e._Data != nil {
		m["data"] = e._Data
	}
	if e._Mod != nil {
		m["mod"] = e._Mod
	}
	if e._Del != nil {
		m["del"] = e._Del
	}
	return m
}
