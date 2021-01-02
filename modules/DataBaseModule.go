package modules

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/team-zf/framework/logger"
	"github.com/team-zf/framework/messages"
	"github.com/team-zf/framework/utils/threads"
	"sync/atomic"
	"time"
)

type DataBaseModule struct {
	conn      *sql.DB                          // 数据库连接对象
	chanNum   int                              // 通道缓存空间
	timeout   time.Duration                    // 超时时长
	logicList map[int]*DataBaseThread          //子逻辑列表
	keyList   []int                            // Key列表, 用来间隔遍历
	chanList  chan []messages.IDataBaseMessage // 消息信通
	getNum    int64                            // 收到的总消息数
	saveNum   int64                            // 保存次数
	thgo      *threads.ThreadGo                // 子协程管理
}

func (e *DataBaseModule) Init() {
	e.chanList = make(chan []messages.IDataBaseMessage, e.chanNum)
}

func (e *DataBaseModule) Start() {
	e.thgo.Go(e.Handle)
}

func (e *DataBaseModule) Stop() {
	close(e.chanList)
	e.thgo.CloseWait()
}

func (e *DataBaseModule) PrintStats() string {
	return fmt.Sprintf(
		"\r\n\t\tDataBase Module\t:%d/%d/%d\t(logic/get/save)",
		len(e.logicList),
		atomic.LoadInt64(&e.getNum),
		atomic.LoadInt64(&e.saveNum))
}

func (e *DataBaseModule) Handle(ctx context.Context) {
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
				lth = NewDataBaseThread()
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

type DataBaseThread struct {
	DBThreadID int                                  //协程ID
	upDataList map[string]messages.IDataBaseMessage //缓存要更新的数据
	chanList   chan []messages.IDataBaseMessage     //收要更新的数据
	Conndb     *sql.DB                              //数据库连接对象
	upTime     time.Time                            //更新时间
	cancel     context.CancelFunc                   //关闭
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
}

func (e *DataBaseThread) AddMsg(msgs []messages.IDataBaseMessage) {
	e.upTime = time.Now()
	e.chanList <- msgs
}

func (e *DataBaseThread) Save() {
	if tx, err := e.Conndb.Begin(); err != nil {
		threads.Try(
			func() {
				for _, data := range e.upDataList {
					if err = data.SaveDB(); err != nil {
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

func NewDataBaseThread() *DataBaseThread {
	return &DataBaseThread{}
}