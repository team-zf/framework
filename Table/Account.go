package Table

import (
	"github.com/team-zf/framework/dal"
	"time"
)

type Account struct {
	dal.BaseTable

	Id           int64     `db:"id,pk,get,!mod"`
	UserName     string    `db:"username"`
	PassWord     string    `db:"password"`
	Email        string    `db:"email"`
	Phone        string    `db:"phone"`
	Token        string    `db:"token"json:"token"`
	Status       int       `db:"status"json:"status"`
	LatelyServer string    `db:"lately_server"`
	CreateTime   time.Time `db:"create_time"`
}

func NewAccount() *Account {
	result := new(Account)
	result.BaseTable.Init(result)
	return result
}
