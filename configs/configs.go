/*
Package configs用于存放常量和常用变量，尽量不引入CTFgo的包，可以引入系统包。
*/
package configs

import (
	"fmt"
	"os"
	"time"
)

var Work_dir, Log_dir, Current_log_path, Save_log_path string

//init初始化常量。
func init() {
	Work_dir, _ = os.Getwd()
	Log_dir = Work_dir + "/logs"
	Current_log_path = Log_dir + "/run.log"
	Save_log_path = Log_dir + "/" + log_times() + ".log"
}

//log_times设置日志文件名，格式如2021-4-15-14_55。
func log_times() string {
	// 东八区，先默认这个，后面再改成动态配置的
	time_zone := time.FixedZone("CST", 8*3600)
	n := time.Now().In(time_zone)
	// 年
	year := n.Year()
	// 月
	month := n.Month()
	// 日
	day := n.Day()
	// 时
	hour := n.Hour()
	// 分
	minute := n.Minute()
	// 秒
	//second := n.Second()
	// 获取时间，格式如2021-4-15-14_55
	t := fmt.Sprintf("%d-%d-%d-%d_%d", year, month, day, hour, minute)
	return t
}

//Times()用于获取当前时间，格式如2006/01/02 15:04:05。
func Times() string {
	// 东八区，先默认这个，后面再改成动态配置的
	time_zone := time.FixedZone("CST", 8*3600)
	n := time.Now().In(time_zone)
	// 获取时间，格式如2006/01/02 15:04:05
	t := n.Format("2006/01/02 15:04:05")
	return t
}
