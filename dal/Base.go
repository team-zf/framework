package dal

import (
	"database/sql"
	"github.com/team-zf/framework/utils"
	"reflect"
	"strings"
)

type ITable interface {
	GetId() int64
	GetTableName() string
}

type IConnDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// 生成Insert&Update语句, SQL
func MarshalModSql(o ITable) string {
	t := reflect.TypeOf(o)
	fields := t.Elem()
	fieldNum := fields.NumField()

	result := utils.NewStringBuilderCap(fieldNum*40 + 60)
	result.Append("insert into ")
	result.Append(o.GetTableName())
	result.Append("(")

	insetTemp := utils.NewStringBuilderCap(fieldNum * 10)
	updateTemp := utils.NewStringBuilderCap(fieldNum * 20)

FieldFor:
	for i := 0; i < fieldNum; i++ {
		field := fields.Field(i)
		tag := field.Tag.Get("db")
		titles := strings.Split(tag, ",")
		name := field.Name
		needMod := true

		if tag == "" {
			continue
		}

		// 遍历所有标题
		for _, v := range titles {
			switch v {
			case "!mod", "pk":
				needMod = false
			case "!":
				continue FieldFor
			default:
				name = v
			}
		}
		// 非首位补逗号分隔
		if !insetTemp.IsEmpty() {
			result.Append(",")
			insetTemp.Append(",")
		}
		result.Append(name)
		insetTemp.Append("?")
		// 加入非主键的字段
		if needMod {
			// 非首位补逗号分隔
			if !updateTemp.IsEmpty() {
				updateTemp.Append(",")
			}
			updateTemp.Append(name)
			updateTemp.Append("=values(")
			updateTemp.Append(name)
			updateTemp.Append(")")
		}
	}

	result.Append(") values (")
	result.Append(insetTemp.ToString())
	result.Append(") on duplicate key update ")
	result.Append(updateTemp.ToString())
	result.Append(";")
	return result.ToString()
}

// 生成Delete语句, SQL
func MarshalDelSql(o ITable, wheres ...string) string {
	t := reflect.TypeOf(o)
	fields := t.Elem()
	fieldNum := fields.NumField()

	result := utils.NewStringBuilderCap(fieldNum*30 + 30)
	where := utils.NewStringBuilderCap(30)

	result.Append("delete from ")
	result.Append(o.GetTableName())

	// 自定义条件
	if len(wheres) > 0 {
		for _, v := range wheres {
			if !where.IsEmpty() {
				where.Append(" and ")
			}
			where.Append(v)
			where.Append("=?")
		}
	}

	// 加入条件
	if !where.IsEmpty() {
		result.Append(" where ")
		result.Append(where.ToString())
	}

	result.Append(";")
	return result.ToString()
}

// 生成Select语句, SQL
func MarshalGetSql(o ITable, wheres ...string) string {
	t := reflect.TypeOf(o)
	fields := t.Elem()
	fieldNum := fields.NumField()

	result := utils.NewStringBuilderCap(fieldNum*30 + 30)
	where := utils.NewStringBuilderCap(30)

	result.Append("select ")

	isColumn := false
FieldFor:
	for i := 0; i < fieldNum; i++ {
		field := fields.Field(i)
		tag := field.Tag.Get("db")
		titles := strings.Split(tag, ",")
		name := field.Name

		if tag == "" {
			continue
		}

		for _, v := range titles {
			switch v {
			case "!mod":
				continue
			case "!":
				continue FieldFor
			default:
				if name == field.Name {
					name = v
				}
			}
		}
		if isColumn {
			result.Append(",")
		} else {
			isColumn = true
		}
		result.Append(name)
	}

	result.Append(" from ")
	result.Append(o.GetTableName())

	// 自定义条件
	if len(wheres) > 0 {
		for _, key := range wheres {
			if !where.IsEmpty() {
				where.Append(" and ")
			}
			where.Append(key)
			where.Append("=?")
		}
	}

	// 加入条件
	if !where.IsEmpty() {
		result.Append(" where ")
		result.Append(where.ToString())
	}

	result.Append(";")
	return result.ToString()
}
