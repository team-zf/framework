package tables

import (
	"errors"
	"github.com/team-zf/framework/utils"
)

type List struct {
	items []interface{}
}

func NewList(datas []interface{}) *List {
	items := make([]interface{}, 0)
	for _, item := range datas {
		switch item.(type) {
		case []interface{}:
			arr := item.([]interface{})
			key := utils.NewStringAny(arr[0]).ToString()
			val := utils.NewStringAny(arr[1])
			items = append(items, NewMap(key, val))
		case interface{}:
			items = append(items, item)
		}
	}
	return &List{items: items}
}

func (e *List) At(idx int) interface{} {
	return e.items[idx]
}

func (e *List) Rnd() interface{} {
	return e.items[utils.Range(len(e.items))]
}

func (e *List) ToIntArray() ([]int, error) {
	arr := make([]int, 0)
	for _, item := range e.items {
		str := utils.NewStringAny(item)
		v, err := str.ToInt()
		if err != nil {
			return nil, err
		}
		arr = append(arr, v)
	}
	return arr, nil
}

func (e *List) ToMapArray() ([]*Map, error) {
	arr := make([]*Map, 0)
	for _, item := range e.items {
		switch item.(type) {
		case *Map:
			arr = append(arr, item.(*Map))
		default:
			return nil, errors.New("data type not *Map")
		}
	}
	return arr, nil
}
