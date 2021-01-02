package tables

import "github.com/team-zf/framework/utils"

type Row struct {
	key   int
	items map[string]interface{}
}

func NewRow() *Row {
	return nil
}

func (e *Row) GetDataToString(name string) string {
	if e.items[name] == nil {
		return ""
	}
	return utils.NewStringAny(e.items[name]).ToString()
}

func (e *Row) GetDataInt(name string) int {
	if e.items[name] == nil {
		return -1
	}
	return utils.NewStringAny(e.items[name]).ToIntV()
}

func (e *Row) GetDataInt64(name string) int64 {
	if e.items[name] == nil {
		return -1
	}
	return utils.NewStringAny(e.items[name]).ToInt64V()
}

func (e *Row) GetDataList(name string) *List {
	if e.items[name] == nil {
		return NewList(nil)
	}
	return NewList(e.items[name].([]interface{}))
}

func (e *Row) GetDataMap(name string) *Map {
	if e.items[name] == nil {
		return nil
	}
	arr := e.items[name].([]interface{})
	key := utils.NewStringAny(arr[0]).ToString()
	val := utils.NewStringAny(arr[1])
	return NewMap(key, val)
}
