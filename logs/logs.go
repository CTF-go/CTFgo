/*
Package logs实现日志相关功能函数。
*/
package logs

import (
	c "CTFgo/configs"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

//Save_log在CTFgo退出时将run.log重命名为当前时间.log，保存在CTFgo/logs文件夹。
func Save_log() {
	err := os.Rename(c.CURRENT_LOG_PATH, c.SAVE_LOG_PATH)
	if err != nil {
		fmt.Printf("the rename operation failed %q\n", err)
	} else {
		fmt.Println("[EXIT] Save logs success!")
	}
}

//info_output输出info信息到控制台和日志文件。
func info_output(msg string) {
	_, _ = fmt.Fprintf(gin.DefaultWriter, "[INFO] [%s] %s \n",
		c.Times(),
		msg,
	)
}

//warning_output输出warning和错误信息到控制台和日志文件。
func warning_output(msg string, err error) {
	_, _ = fmt.Fprintf(gin.DefaultWriter, "[WARNING] [%s] %s - %s\n",
		c.Times(),
		msg,
		err,
	)
}

//error_output输出error和错误信息到控制台和日志文件，然后停止CTFgo程序。
func error_output(msg string, err error) {
	_, _ = fmt.Fprintf(gin.DefaultWriter, "[ERROR] [%s] %s - %s\n",
		c.Times(),
		msg,
		err,
	)
	os.Exit(1)
}

// var params gin.LogFormatterParams

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

//INFO函数传参(msg string)，无err，仅输出提示信息。
func INFO(msg string) {
	info_output(msg)
}

//WARNING函数传参(msg string, err error)，输出报错信息但不退出。
func WARNING(msg string, err error) {
	warning_output(msg, err)
}

//ERROR函数传参(msg string, err error)，输出报错信息并退出程序。
func ERROR(msg string, err error) {
	error_output(msg, err)
}

/*
func ATTACK(msg string, err error) {
	attack_output(params, msg, err)
	return
}
*/
