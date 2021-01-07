package dal

import (
	"fmt"
	"runtime"
)

type DalMessage struct {
	UserId  int64
	Table   ITable
	RunFunc func(conn IConnDB) error
}

func (e *DalMessage) DBThreadID() int {
	cpu := runtime.NumCPU() * 10
	return int(e.UserId % int64(cpu))
}

func (e *DalMessage) GetDataKey() string {
	return fmt.Sprintf("%d_%s_%d", e.UserId, e.Table.GetTableName(), e.Table.GetId())
}

func (e *DalMessage) SaveDB(db IConnDB) error {
	return e.RunFunc(db)
}
