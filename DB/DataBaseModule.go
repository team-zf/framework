package DB

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/messages"
	"github.com/team-zf/framework/modules"
	"github.com/team-zf/framework/utils/threads"
	"runtime"
	"sync/atomic"
	"time"
)

type DataBaseModule struct {
	name      string
	dsn       string
	db        *sql.DB
	chanNum   int                              // 通道缓存空间
	timeout   time.Duration                    // 超时时长
	logicList map[int]*DataBaseThread          // 子逻辑列表
	keyList   []int                            // Key列表, 用来间隔遍历
	chanList  chan []messages.IDataBaseMessage // 消息信通
	getNum    int64                            // 收到的总消息数
	saveNum   int64                            // 保存次数
	thgo      *threads.ThreadGo                // 子协程管理
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
	e.chanList = make(chan []messages.IDataBaseMessage, e.chanNum)
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
		"\r\n\t\t%s的状态:\t%d/%d/%d\t(Logic/Save/Request)",
		e.name,
		len(e.logicList),
		atomic.LoadInt64(&e.saveNum),
		atomic.LoadInt64(&e.getNum))
}

func (e *DataBaseModule) Handle() {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	loop := 0
	for {
		select {
		case msgs, ok := <-e.chanList:
			if !ok {
				for _, lth := range e.logicList {
					lth.Stop()
				}
				return
			}
			if len(msgs) == 0 {
				continue
			}
			upmd := msgs[0]
			lth, ok := e.logicList[upmd.DBThreadID()]
			if !ok {
				// 新开一个协程
				lth = NewDataBaseThread(
					upmd.DBThreadID(),
					e.chanNum,
					e.db,
				)
				e.logicList[lth.DBThreadID] = lth
				e.keyList = append(e.keyList, lth.DBThreadID)
				lth.Start(e)
			}
			lth.AddMsg(msgs)
		case <-t.C:
			if len(e.keyList) == 0 {
				continue
			}
			loop = loop % len(e.keyList)
			keyid := e.keyList[loop]
			if lth, ok := e.logicList[keyid]; ok {
				if lth.GetMsgNum() == 0 && time.Now().Sub(lth.upTime) > e.timeout {
					lth.Stop()
					delete(e.logicList, keyid)
					e.keyList = append(e.keyList[:loop], e.keyList[loop+1:]...)
				}
			}
			loop++
		}
	}
}

func (e *DataBaseModule) AddMsg(msgs ...messages.IDataBaseMessage) {
	atomic.AddInt64(&e.getNum, 1)
	e.chanList <- msgs
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
		name:      "DataBase",
		chanNum:   1024,
		timeout:   2 * time.Minute,
		logicList: make(map[int]*DataBaseThread, runtime.NumCPU()*10),
		keyList:   make([]int, 0),
		thgo:      threads.NewThreadGo(),
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

type DataBaseThread struct {
	DBThreadID int                                  // 协程ID
	upDataList map[string]messages.IDataBaseMessage // 缓存要更新的数据
	chanList   chan []messages.IDataBaseMessage     // 收要更新的数据
	Conndb     *sql.DB                              // 数据库连接对象
	upTime     time.Time                            // 更新时间
	cancel     context.CancelFunc                   // 关闭
}

func (e *DataBaseThread) Start(mod *DataBaseModule) {
	e.cancel = mod.thgo.SubGo(
		func(ctx context.Context) {
			e.Handle(ctx, mod)
		},
	)
}

func (e *DataBaseThread) Stop() {
	e.cancel()
	close(e.chanList)
}

func (e *DataBaseThread) Handle(ctx context.Context, mod *DataBaseModule) {
	tk := time.NewTimer(time.Second)
	defer tk.Stop()
	isruned := false
trheadhandle:
	for {
		select {
		case msg, ok := <-e.chanList:
			{
				if !ok {
					e.Save()
					atomic.AddInt64(&mod.saveNum, 1)
					break trheadhandle
				}
				if len(msg) == 0 {
					continue
				}
				for _, data := range msg {
					e.upDataList[data.GetDataKey()] = data
				}
				if isruned {
					tk.Reset(time.Second)
					isruned = false
				}

			}
		case <-tk.C:
			{
				if len(e.upDataList) > 0 {
					e.Save()
					atomic.AddInt64(&mod.saveNum, 1)
					e.upDataList = make(map[string]messages.IDataBaseMessage)
				}
				isruned = true
			}
		}
	}
}

func (e *DataBaseThread) AddMsg(msgs []messages.IDataBaseMessage) {
	e.upTime = time.Now()
	e.chanList <- msgs
}

func (e *DataBaseThread) Save() {
	if tx, err := e.Conndb.Begin(); err == nil {
		threads.Try(
			func() {
				for _, data := range e.upDataList {
					if err = data.SaveDB(tx); err != nil {
						panic(errors.New(fmt.Sprintf("keyid: %d; DataKey: %s", data.DBThreadID(), data.GetDataKey())))
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
}

func (e *DataBaseThread) GetMsgNum() int {
	return len(e.chanList) + len(e.upDataList)
}

func NewDataBaseThread(id, channum int, db *sql.DB) *DataBaseThread {
	return &DataBaseThread{
		DBThreadID: id,
		upDataList: make(map[string]messages.IDataBaseMessage),
		chanList:   make(chan []messages.IDataBaseMessage, channum),
		Conndb:     db,
	}
}
