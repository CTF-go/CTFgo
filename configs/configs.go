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

var (
	WORK_DIR, LOG_DIR, CURRENT_LOG_PATH, SAVE_LOG_PATH, DB_DIR, DB_FILE, SESSION_DIR, SESSION_ID string
	START_TIME, END_TIME                                                                         int64
)

// init初始化常量。
func init() {
	WORK_DIR, _ = os.Getwd()
	fmt.Println("CTFgo Work dir is:", WORK_DIR)
	SESSION_DIR = WORK_DIR + "/sessions"
	SESSION_ID = "CTFGOSESSID"
	LOG_DIR = WORK_DIR + "/logs"
	CURRENT_LOG_PATH = LOG_DIR + "/run.log"
	SAVE_LOG_PATH = LOG_DIR + "/" + log_times() + ".log"
	DB_DIR = WORK_DIR + "/databases"
	DB_FILE = DB_DIR + "/ctfgo.db"
	START_TIME = 1637974800
	END_TIME = 1638104400
}

// log_times设置日志文件名，格式如2021-4-15-14_55。
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

// HexBotSayTime。
func HexBotSayTime() string {
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
	second := n.Second()
	// 获取时间，格式如2021/4/15 14:55:22
	t := fmt.Sprintf("%d/%d/%d %d:%d:%d", year, month, day, hour, minute, second)
	return t
}

// Times 用于获取当前时间，格式如2006/01/02 15:04:05。
func Times() string {
	// 东八区，先默认这个，后面再改成动态配置的
	time_zone := time.FixedZone("CST", 8*3600) // 8*3600 = 8h
	n := time.Now().In(time_zone)
	// 获取时间，格式如2006/01/02 15:04:05
	t := n.Format("2006/01/02 15:04:05")
	return t
}

// Timestamp 用于获取当前10位数时间戳。
func Timestamp() int {
	// time_zone := time.FixedZone("UTC", 0)
	// t := time.Now().In(time_zone).Unix()
	t := time.Now().Unix()
	return int(t)
}

// MD5 进行md5加密。
func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	//将[]byte转成16进制
	md5_str := fmt.Sprintf("%x", has)
	return md5_str
}

// Random 生成随机数。
func Random() []byte {
	b := make([]byte, 32)
	//ReadFull从rand.Reader精确地读取len(b)字节数据填充进b
	//rand.Reader是一个全局、共享的密码用强随机数生成器
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		fmt.Printf("random number generation error: %v", err)
	}
	return b
}

// Token 生成随机token。
func Token() string {
	b := Random()[:16]
	return fmt.Sprintf("%x", b)
}

// CheckID 验证id是否为非负正整数。
func CheckID(id string) bool {
	if id == "1" {
		return false
	}
	pattern := `^[1-9]\d*$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(id)
}

/*
// NumCompare 判断s1是否小于等于s2，小于等于返回true，大于返回false。
func NumCompare(s1, s2 string) bool {
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
*/
