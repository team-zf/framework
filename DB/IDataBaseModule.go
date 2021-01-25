package DB

import (
	"database/sql"
	"github.com/team-zf/framework/modules"
)

type IDataBaseModule interface {
	modules.IModule
	AddMsg(msgs ...IDataBaseMessage)
	GetDB() *sql.DB
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}
