package utils

import (
	"reflect"
)

// 反射创建新对象。
func ReflectNew(target interface{}) interface{} {
	t := reflect.TypeOf(target)
	// 指针类型获取真正type需要调用Elem
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 调用反射创建对象
	newStruc := reflect.New(t)
	return newStruc.Interface()
}
