package Data

import (
	"github.com/team-zf/framework/Control"
	"github.com/team-zf/framework/Table"
	"github.com/team-zf/framework/dal"
)

func GetServerList() (serverList []*Table.Server) {
	serverList = make([]*Table.Server, 0)
	sqlstr := dal.MarshalGetSql(Table.NewServer())
	rows, err := Control.DbModule.Query(sqlstr)
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			server := Table.NewServer()
			rows.Scan(
				&server.Id,
				&server.Title,
				&server.Status,
				&server.Host,
				&server.Port,
				&server.ServiceTime,
			)
			serverList = append(serverList, server)
		}
	}
	return
}

func GetServerById(serverId int64) (server *Table.Server) {
	server = Table.NewServer()
	sqlstr := dal.MarshalGetSql(server, "id")
	row := Control.DbModule.QueryRow(sqlstr, serverId)
	if row.Scan(
		&server.Id,
		&server.Title,
		&server.Status,
		&server.Host,
		&server.Port,
		&server.ServiceTime,
	) != nil {
		server = nil
	}
	return
}
