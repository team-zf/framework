package tables

import "github.com/team-zf/framework/utils"

type Map struct {
	key   string
	value *utils.String
}

func NewMap(key string, value *utils.String) *Map {
	return &Map{
		key:   key,
		value: value,
	}
}
