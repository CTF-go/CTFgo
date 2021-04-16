package logs

//日志功能函数

import (
	c "CTFgo/configs"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var params gin.LogFormatterParams

//Save_log在CTFgo退出时将run.log重命名为当前时间.log，保存在CTFgo/logs文件夹。
func Save_log() {
	err := os.Rename(c.Current_log_path, c.Save_log_path)
	if err != nil {
		fmt.Printf("the rename operation failed %q\n", err)
	} else {
		fmt.Println("save logs success!")
	}
}

func info_output(msg string, err error) {
	_, _ = fmt.Fprintf(gin.DefaultWriter, "[INFO] [%s] %s \n",
		c.Times(),
		msg,
	)
	return
}

func warning_output(msg string, err error) {
	_, _ = fmt.Fprintf(gin.DefaultWriter, "[WARNING] [%s] %s - %s\n",
		c.Times(),
		msg,
		err,
	)
	return
}

func error_output(msg string, err error) {
	_, _ = fmt.Fprintf(gin.DefaultWriter, "[ERROR] [%s] %s - %s\n",
		c.Times(),
		msg,
		err,
	)
	os.Exit(1)
	return
}

/*恶意攻击日志后面再实现，单独存日志文件
func attack_output(params gin.LogFormatterParams, msg string, err error) {
	_, _ = fmt.Fprintf(gin.DefaultWriter, "[ATTACK] %s - [%s] \"%s %s %s %3d %s \"%s\" %s\"\n",
		params.ClientIP,
		params.TimeStamp.Format("2006/01/02 15:04:05"),
		params.Method,
		params.Path,
		params.Request.Proto,
		params.StatusCode,
		params.Latency,
		params.Request.UserAgent(),
		params.ErrorMessage,
	)
	return
}
*/

func INFO(msg string, err error) {
	info_output(msg, err)
	return
}

func WARNING(msg string, err error) {
	warning_output(msg, err)
	return
}

func ERROR(msg string, err error) {
	error_output(msg, err)
	return
}

/*
func ATTACK(msg string, err error) {
	attack_output(params, msg, err)
	return
}
*/
