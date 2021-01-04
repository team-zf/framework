package framework

import (
	"github.com/team-zf/framework/modules"
)

func CreateApp(opts ...modules.AppOptions) modules.IApp {
	return modules.NewApp(opts...)
}
