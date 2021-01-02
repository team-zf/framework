package utils

import (
	"strconv"
	"strings"
)

type String struct {
	value string
}

func NewString(v string) *String {
	return &String{value: v}
}

func NewStringInt(v int) *String {
	return &String{value: strconv.Itoa(v)}
}

func NewStringInt64(v int64) *String {
	return &String{value: strconv.FormatInt(v, 10)}
}

func NewStringFloat64(v float64) *String {
	return &String{value: strconv.FormatFloat(v, 'f', -1, 64)}
}

func NewStringBool(v bool) *String {
	return &String{value: strconv.FormatBool(v)}
}

func NewStringAny(v interface{}) *String {
	var str *String
	switch v.(type) {
	case string:
		str = NewString(v.(string))
	case int:
		str = NewStringInt(v.(int))
	case int8, int16, int32, int64:
		str = NewStringInt64(v.(int64))
	case float32, float64:
		str = NewStringFloat64(v.(float64))
	case bool:
		str = NewStringBool(v.(bool))
	case uint, uint8, uint16, uint32, uint64:
		str = &String{value: strconv.FormatUint(v.(uint64), 10)}
	}
	return str
}

func (e *String) Clear() *String {
	var newStr string
	e.value = newStr
	return e
}

func (e *String) ToString() string {
	return e.value
}

func (e *String) ToInt() (int, error) {
	return strconv.Atoi(e.value)
}

func (e *String) ToIntV() int {
	v, err := e.ToInt()
	if err != nil {
		panic(err)
	}
	return v
}

func (e *String) ToInt64() (int64, error) {
	return strconv.ParseInt(e.value, 10, 64)
}

func (e *String) ToInt64V() int64 {
	v, err := e.ToInt64()
	if err != nil {
		panic(err)
	}
	return v
}

func (e *String) ToUint() (uint, error) {
	v, err := e.ToInt()
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

func (e *String) ToUintV() uint {
	v := e.ToIntV()
	return uint(v)
}

func (e *String) ToUint32() (uint32, error) {
	v, err := e.ToInt()
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

func (e *String) ToUint32V() uint32 {
	v := e.ToIntV()
	return uint32(v)
}

func (e *String) ToUint64() (uint64, error) {
	v, err := e.ToInt()
	if err != nil {
		return 0, err
	}
	return uint64(v), nil
}

func (e *String) ToUint64V() uint64 {
	v := e.ToIntV()
	return uint64(v)
}

func (e *String) ToFloat() (float64, error) {
	return strconv.ParseFloat(e.value, 64)
}

func (e *String) ToFloatV() float64 {
	v, err := e.ToFloat()
	if err != nil {
		panic(err)
	}
	return v
}

func (e *String) ToBoolV() bool {
	v, err := strconv.ParseBool(e.ToString())
	if err != nil {
		panic(err)
	}
	return v
}

func (e *String) ToArray() []string {
	return strings.Split(e.value, "")
}

func (e *String) ToLower() *String {
	return NewString(strings.ToLower(e.value))
}

func (e *String) ToUpper() *String {
	return NewString(strings.ToUpper(e.value))
}

func (e *String) Len() int {
	return len(e.value)
}

func (e *String) StartsWith(s string) bool {
	return e.SubstrEnd(len(s)).ToString() == s
}

func (e *String) EndsWith(s string) bool {
	return e.SubstrBegin(e.Len()-len(s)).ToString() == s
}

func (e *String) Trim() *String {
	return NewString(strings.Trim(e.value, " "))
}

func (e *String) Replace(old, new string) *String {
	return NewString(strings.Replace(e.value, old, new, -1))
}

func (e *String) Substr(beginIndex, endIndex int) *String {
	return NewString(e.value[beginIndex:endIndex])
}

func (e *String) SubstrBegin(beginIndex int) *String {
	return e.Substr(beginIndex, e.Len())
}

func (e *String) SubstrEnd(endIndex int) *String {
	return e.Substr(0, endIndex)
}
