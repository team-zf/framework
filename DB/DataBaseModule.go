package DB

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/modules"
	"github.com/team-zf/framework/utils/threads"
	"sync/atomic"
	"time"
)

type DataBaseModule struct {
	name         string
	dsn          string
	db           *sql.DB
	thgo         *threads.ThreadGo
	chanList     chan []IDataBaseMessage     // 消息信通
	cacheList    map[string]IDataBaseMessage // 缓存要更新的数据
	requestCount int64                       // 收到的请求总数
	saveCount    int64                       // 保存的总数
}

func (e *DataBaseModule) Init() {
	db, err := sql.Open("mysql", e.dsn)
	if err != nil {
		panic(fmt.Sprintf("Mysql连接失败, 错误原因: %+v", err))
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(600 * time.Second)
	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("Mysql尝Ping失败, 错误原因: %+v", err))
	}
	e.db = db
	e.chanList = make(chan []IDataBaseMessage, 1024)
	e.cacheList = make(map[string]IDataBaseMessage)
}

func (e *DataBaseModule) Start() {
	e.thgo.Go(func(ctx context.Context) {
		logger.Notice("%s启动", e.name)
		e.Handle()
	})
}

func (e *DataBaseModule) Stop() {
	close(e.chanList)
	e.db.Close()
	e.thgo.CloseWait()
	logger.Notice("%s已停止", e.name)
}

func (e *DataBaseModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n\t\t%s的状态:\t%d/%d/%d\t(Cache/Save/Request)",
		e.name,
		len(e.cacheList),
		atomic.LoadInt64(&e.saveCount),
		atomic.LoadInt64(&e.requestCount))
}

func (e *DataBaseModule) Handle() {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		select {
		case msgs, ok := <-e.chanList:
			if ok {
				for _, msg := range msgs {
					e.cacheList[msg.GetDataKey()] = msg
				}
			} else {
				e.Save()
				return
			}
		case <-t.C:
			e.Save()
		}
	}
}

func (e *DataBaseModule) Save() {
	if len(e.cacheList) == 0 {
		return
	}

	atomic.AddInt64(&e.saveCount, 1)
	if tx, err := e.db.Begin(); err == nil {
		threads.Try(
			func() {
				for _, msg := range e.cacheList {
					if err = msg.SaveDB(tx); err != nil {
						str := fmt.Sprintf("DataKey: %s; Error: %+v", msg.GetDataKey(), err)
						panic(errors.New(str))
					}
				}
				tx.Commit()
			},
			func(err error) {
				tx.Rollback()
				logger.Error(err.Error())
			},
		)
	}
	e.cacheList = make(map[string]IDataBaseMessage)
}

func (e *DataBaseModule) AddMsg(msgs ...IDataBaseMessage) {
	if len(msgs) > 0 {
		atomic.AddInt64(&e.requestCount, 1)
		e.chanList <- msgs
	}
}

func (e *DataBaseModule) GetDB() *sql.DB {
	return e.db
}

func (e *DataBaseModule) QueryRow(query string, args ...interface{}) *sql.Row {
	return e.db.QueryRow(query, args...)
}

func (e *DataBaseModule) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return e.db.Query(query, args...)
}

func (e *DataBaseModule) Exec(query string, args ...interface{}) (sql.Result, error) {
	return e.db.Exec(query, args...)
}

func NewDataBaseModule(opts ...modules.ModOptions) *DataBaseModule {
	result := &DataBaseModule{
		name: "DataBase",
		thgo: threads.NewThreadGo(),
	}
	for _, opt := range opts {
		opt(result)
	}
	return result
}

func DataBaseSetName(v string) modules.ModOptions {
	return func(mod modules.IModule) {
		mod.(*DataBaseModule).name = v
	}
}

func DataBaseSetDsn(v string) modules.ModOptions {
	return func(mod modules.IModule) {
		mod.(*DataBaseModule).dsn = v
	}
}
