/*
Package configs用于存放常量和常用变量，不引入CTFgo的包，可以引入系统包。
*/
package configs

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"
)

var Work_dir, Log_dir, Current_log_path, Save_log_path, DB_dir, DB_file, Static_path, Session_dir, SessionID string

//init初始化常量。
func init() {
	Work_dir, _ = os.Getwd()
	Session_dir = Work_dir + "/sessions"
	SessionID = "CTFGOSESSID"
	Log_dir = Work_dir + "/logs"
	Current_log_path = Log_dir + "/run.log"
	Save_log_path = Log_dir + "/" + log_times() + ".log"
	//Key_file = Work_dir + "/security.key"
	DB_dir = Work_dir + "/databases"
	DB_file = DB_dir + "/ctfgo.db"
	Static_path = Work_dir + "/themes/default"

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

//Times用于获取当前时间，格式如2006/01/02 15:04:05。
func Times() string {
	// 东八区，先默认这个，后面再改成动态配置的
	time_zone := time.FixedZone("CST", 8*3600)
	n := time.Now().In(time_zone)
	// 获取时间，格式如2006/01/02 15:04:05
	t := n.Format("2006/01/02 15:04:05")
	return t
}

//Timestamp用于获取当前10位数时间戳。
func Timestamp() int32 {
	// 东八区，先默认这个，后面再改成动态配置的
	time_zone := time.FixedZone("CST", 8*3600)
	t := time.Now().In(time_zone).Unix()
	return int32(t)
}

//MD5进行md5加密。
func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	//将[]byte转成16进制
	md5_str := fmt.Sprintf("%x", has)
	return md5_str
}

//Random生成随机数。
func Random() []byte {
	b := make([]byte, 32)
	//ReadFull从rand.Reader精确地读取len(b)字节数据填充进b
	//rand.Reader是一个全局、共享的密码用强随机数生成器
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		fmt.Printf("random number generation error: %v", err)
	}
	return b
}

//Token生成随机token。
func Token() string {
	b := Random()[:16]
	return fmt.Sprintf("%x", b)
}

//ID_verify 验证id是否为非负正整数。
func ID_verify(id string) bool {
	pattern := `^[1-9]\d*$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(id)
}

//num_compare 判断s1是否小于等于s2，小于等于返回true，大于返回false。
func Num_compare(s1, s2 string) bool {
	if len(s1) == len(s2) {
		if s1 <= s2 {
			return true
		} else {
			return false
		}
	} else if len(s1) < len(s2) {
		return true
	}
	return false
}
