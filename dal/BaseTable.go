package dal

import (
	"reflect"
	"strings"
	"time"
)

type BaseTable struct {
	subclass interface{}
}

func (e *BaseTable) Init(subclass interface{}) {
	e.subclass = subclass
}

func (e *BaseTable) GetId() (id int64) {
	if e.subclass != nil {
		t := reflect.TypeOf(e.subclass)
		fields := t.Elem()
		fieldNum := fields.NumField()
	__For:
		for i := 0; i < fieldNum; i++ {
			field := fields.Field(i)
			tag := field.Tag.Get("db")
			titles := strings.Split(tag, ",")
			for _, v := range titles {
				if v == "pk" {
					id = e.GetSubClassPropByName(field.Name).Int()
					break __For
				}
			}
		}
	}
	return
}

func (e *BaseTable) GetTableName() (name string) {
	if e.subclass != nil {
		subclass := reflect.TypeOf(e.subclass).Elem()
		name = strings.ToLower(subclass.Name())
	}
	return
}

func (e *BaseTable) GetSubClassPropByName(name string) (val reflect.Value) {
	if e.subclass != nil {
		val = reflect.ValueOf(e.subclass).Elem().FieldByName(name)
	}
	return
}

func (e *BaseTable) ToJsonMap() (result map[string]interface{}) {
	result = make(map[string]interface{})
	if e.subclass != nil {
		t := reflect.TypeOf(e.subclass)
		fields := t.Elem()
		fieldNum := fields.NumField()
		for i := 0; i < fieldNum; i++ {
			field := fields.Field(i)
			tag := field.Tag.Get("json")
			if tag == "" {
				continue
			}
			val := e.GetSubClassPropByName(field.Name).Interface()
			switch val.(type) {
			case time.Time:
				result[tag] = val.(time.Time).Unix()
			default:
				result[tag] = val
			}
		}
	}
	return
}
