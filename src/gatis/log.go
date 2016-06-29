package gatis

import (
	"github.com/astaxie/beego/logs"
)

var Log *logs.BeeLogger

func init() {
	Log = logs.NewLogger(1 << 10)
	Log.SetLogger("console", "")
}