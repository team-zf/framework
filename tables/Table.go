package tables

import (
	"encoding/json"
	"github.com/team-zf/framework/config"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/utils"
	"io/ioutil"
	"os"
	"path"
)

var (
	tables map[string]*Table
)

func init() {
	tables = make(map[string]*Table)
}

func LoadTables(dir string, conf *config.TableConfig) {
	s, err := os.Stat(dir)
	if err != nil || !s.IsDir() {
		return
	}

	prefix := "wx_"
	if conf != nil && conf.Prefix != "" {
		prefix = conf.Prefix
	}
	suffix := ".json"
	if conf != nil && conf.Suffix != "" {
		suffix = conf.Suffix
	}

	files, _ := ioutil.ReadDir(dir)
	logger.Notice("开始载入数据表.")
	for _, file := range files {
		fileName := utils.NewString(file.Name())
		v1 := file.IsDir() == false
		v2 := fileName.StartsWith(prefix)
		v3 := fileName.EndsWith(suffix)
		if v1 && v2 && v3 {
			x := len(prefix)
			y := len(fileName.ToString()) - len(suffix)
			name := fileName.Substr(x, y).ToString()
			table, _ := loadTableFile(path.Join(dir, file.Name()))
			tables[name] = table
			logger.Notice("File: %s; Key: %s", file.Name(), name)
		}
	}
	logger.Notice("数据表载入完成.")
}

func loadTableFile(filePath string) (*Table, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, _ := ioutil.ReadAll(file)
	datas := make([]interface{}, 0)
	if err := json.Unmarshal(bytes, &datas); err != nil {
		return nil, err
	}
	return NewTable(datas), nil
}

func GetTable(name string) *Table {
	return tables[name]
}

type Table struct {
	rows []*Row
}

func NewTable(datas []interface{}) *Table {
	rows := make([]*Row, 0)
	for _, data := range datas {
		arr := data.([]interface{})
		key := utils.NewStringAny(arr[0]).ToIntV()
		items := make(map[string]interface{})
		for _, v := range arr[1].([]interface{}) {
			l := v.([]interface{})
			k := utils.NewStringAny(l[0]).ToString()
			items[k] = l[1]
		}
		rows = append(rows, &Row{
			key:   key,
			items: items,
		})
	}
	return &Table{rows: rows}
}

func (e *Table) At(idx int) *Row {
	return e.rows[idx]
}

func (e *Table) All() []*Row {
	return e.rows
}

func (e *Table) ByKey(key int) *Row {
	for _, row := range e.rows {
		if row.key == key {
			return row
		}
	}
	return nil
}

func (e *Table) ByProp(prop string, val interface{}) []*Row {
	if _, ok := val.(int); ok {
		val = utils.NewStringAny(val).ToFloatV()
	}
	arr := make([]*Row, 0)
	for _, row := range e.rows {
		if row.items[prop] == val {
			arr = append(arr, row)
		}
	}
	return arr
}
