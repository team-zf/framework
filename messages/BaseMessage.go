package messages

import (
	"github.com/team-zf/framework/dal"
)

type IDataBaseMessage interface {
	// 所在DB协程
	DBThreadID() int
	// 数据表,如果你的表放入时，不是马上保存的，那么后续可以用这个KEY来进行覆盖，
	// 这样就可以实现多次修改一次保存的功能
	// 所以这个字段建议是：用户ID+数据表名+数据主键
	GetDataKey() string
	// 调用方法
	SaveDB(db dal.IConnDB) error
}
