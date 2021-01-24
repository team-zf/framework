package Table

import (
	"encoding/json"
	"github.com/team-zf/framework/utils"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

func LoadTable(filePath string, st interface{}) ([]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buff, _ := ioutil.ReadAll(file)
	datas := make([]interface{}, 0)
	json.Unmarshal(buff, &datas)

	rows := make([]interface{}, 0)
	for _, data := range datas {
		elem := utils.ReflectNew(st)
		__LoadRow(data.([]interface{}), elem)
		rows = append(rows, elem)
	}
	return rows, nil
}

func __LoadRow(data []interface{}, st interface{}) {
	fields := reflect.TypeOf(st).Elem()
	values := reflect.ValueOf(st).Elem()
	fieldNum := fields.NumField()

	for i := 0; i < fieldNum; i++ {
		tag := fields.Field(i).Tag.Get("ST")
		val := values.Field(i)

		// 非静态表字段
		if tag == "" {
			continue
		}

		// 主键
		if strings.ToUpper(tag) == "PK" {
			val.SetInt(utils.NewStringAny(data[0]).ToInt64V())
		} else {
			field := __GetFieldByTag(data[1].([]interface{}), tag)
			switch val.Interface().(type) {
			case int:
				if field == nil {
					val.SetInt(-1)
				} else {
					val.SetInt(utils.NewStringAny(field).ToInt64V())
				}
			case string:
				if field == nil {
					val.SetString("")
				} else {
					val.SetString(utils.NewStringAny(field).ToString())
				}
			case bool:
				if field == nil {
					val.SetBool(false)
				} else {
					str := utils.NewStringAny(field).ToString()
					val.SetBool(str == "1" || strings.ToUpper(str) == "TRUE")
				}
			case *List:
				if field == nil {
					val.Set(reflect.ValueOf(NewList(nil)))
				} else {
					val.Set(reflect.ValueOf(NewList(field.([]interface{}))))
				}
			case *Map:
				if field == nil {
					//val.Set(reflect.ValueOf(nil))
				} else {
					val.Set(reflect.ValueOf(NewMap(field.([]interface{}))))
				}
			}
		}
	}
}

func __GetFieldByTag(fields []interface{}, tag string) interface{} {
	for _, field := range fields {
		arr := field.([]interface{})
		key := utils.NewStringAny(arr[0]).ToString()
		if key == tag {
			return arr[1]
		}
	}
	return nil
}
