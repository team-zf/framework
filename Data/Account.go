package Data

import (
	"crypto/md5"
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
	"github.com/team-zf/framework/Control"
	"github.com/team-zf/framework/Table"
	"github.com/team-zf/framework/dal"
	"math/rand"
	"time"
)

func GetAccountByToken(token string) (account *Table.Account) {
	account = Table.NewAccount()
	sqlstr := dal.MarshalGetSql(account, "token")
	row := Control.DbModule.QueryRow(sqlstr, token)
	if row.Scan(
		&account.Id,
		&account.UserName,
		&account.PassWord,
		&account.Email,
		&account.Phone,
		&account.Token,
		&account.Status,
		&account.LatelyServer,
		&account.CreateTime,
	) != nil {
		account = nil
	}
	return
}

func GetAccountByUserName(userName string) (account *Table.Account) {
	account = Table.NewAccount()
	sqlstr := dal.MarshalGetSql(account, "username")
	row := Control.DbModule.QueryRow(sqlstr, userName)
	if row.Scan(
		&account.Id,
		&account.UserName,
		&account.PassWord,
		&account.Email,
		&account.Phone,
		&account.Token,
		&account.Status,
		&account.LatelyServer,
		&account.CreateTime,
	) != nil {
		account = nil
	}
	return
}

func RefreshAccountToken(account *Table.Account) bool {
	token, ok := GenerateToken()
	if ok {
		sqlstr := `update account set token = ? where username = ?;`
		_, err := Control.DbModule.Exec(sqlstr, token, account.UserName)
		if err == nil {
			account.Token = token
			return true
		}
	}
	return false
}

func RegisterAccount(account *Table.Account) bool {
	sqlstr := `insert into account(id, username, password, email, phone, token, status, lately_server, create_time) values (?, ?, ?, ?, ?, ?, ?, ?, ?);`
	_, err := Control.DbModule.Exec(sqlstr,
		account.Id,
		account.UserName,
		account.PassWord,
		account.Email,
		account.Phone,
		account.Token,
		account.Status,
		account.LatelyServer,
		account.CreateTime)
	return err == nil
}

func GenerateAccountId() (accountId int64, ok bool) {
	retryCount := 0
	rand.Seed(time.Now().Unix())
	for {
		retryCount++
		// 重试次数达到100次以后直接宣布生成失败
		if retryCount > 100 {
			return 0, false
		}
		accountId = rand.Int63n(8999999) + 1000000 // 随机生成一个七位数ID
		sqlstr := `select username from account where id = ?;`
		row := Control.DbModule.QueryRow(sqlstr, accountId)
		username := ""
		if row.Scan(&username) != nil {
			break
		}
	}
	return accountId, true
}

func GenerateToken() (token string, ok bool) {
	retryCount := 0
	rand.Seed(time.Now().Unix())
	for {
		retryCount++
		// 重试次数达到100次以后直接宣布生成失败
		if retryCount > 100 {
			return "", false
		}
		h := md5.New()
		h.Write([]byte(uuid.NewV1().String()))
		token = hex.EncodeToString(h.Sum(nil))
		sqlstr := `select id from account where token = ?;`
		row := Control.DbModule.QueryRow(sqlstr, token)
		id := 0
		if row.Scan(&id) != nil {
			break
		}
	}
	return token, true
}
