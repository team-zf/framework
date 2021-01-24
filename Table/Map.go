package Table

import (
	"github.com/team-zf/framework/utils"
)

type Map struct {
	key   string
	value interface{}
}

func (e *Map) Key() string {
	return e.key
}

func (e *Map) KeyToInt() int {
	return utils.NewStringAny(e.key).ToIntV()
}

func (e *Map) Value() interface{} {
	return e.value
}

func (e *Map) ValueToInt() int {
	return utils.NewStringAny(e.value).ToIntV()
}

func (e *Map) ValueToString() string {
	return utils.NewStringAny(e.value).ToString()
}

func (e *Map) ValueToFloat() float64 {
	return utils.NewStringAny(e.value).ToFloatV()
}

func NewMap(data []interface{}) *Map {
	return &Map{
		key:   utils.NewStringAny(data[0]).ToString(),
		value: data[1],
	}
}
