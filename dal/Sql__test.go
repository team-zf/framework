package dal

import (
	"fmt"
	"testing"
	"time"
)

type Account struct {
	BaseTable

	Id         int64     `db:"id,pk"`
	UserName   string    `db:"username,!mod"`
	PassWord   string    `db:"password"`
	CreateTime time.Time `db:"create_time"`
}

func NewAccount() *Account {
	result := new(Account)
	result.BaseTable.Init(result)
	return result
}

func TestSql(t *testing.T) {
	account := NewAccount()
	account.Id = 100
	account.PassWord = "hello"
	fmt.Println(MarshalModSql(account))
	fmt.Println(MarshalGetSql(account))
	fmt.Println(MarshalDelSql(account, "id", "username"))
}
