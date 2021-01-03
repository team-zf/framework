package Table

import (
	"github.com/team-zf/framework/dal"
	"time"
)

type Server struct {
	dal.BaseTable

	Id          int64     `db:"id,pk"json:"id"`
	Title       string    `db:"title"json:"title"`
	Status      int       `db:"status"json:"status"`
	Host        string    `db:"host"`
	Port        int       `db:"port"`
	ServiceTime time.Time `db:"service_time"json:"service_time"`
}

func NewServer() *Server {
	result := new(Server)
	result.BaseTable.Init(result)
	return result
}
