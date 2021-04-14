package apiInit

import (
	u "CTFgo/api/user"
	c "CTFgo/configs"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var Log_path *os.File

// SetupRouter用于注册api，对应功能实现在CTFgo/api下
func SetupRouter() *gin.Engine {
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DisableConsoleColor()
	err := os.MkdirAll(c.Log_dir, 0777)
	if err != nil {
		fmt.Printf("create logs dir error, err:%v\n", err)
	}
	// 创建run.log
	Log_path, _ = os.Create(c.Current_log_path)
	// 将log输出到控制台和文件
	gin.DefaultWriter = io.MultiWriter(Log_path, os.Stdout)
	c := gin.LoggerConfig{
		Output:    gin.DefaultWriter,
		// 需要跳过记录log的Api
		SkipPaths: []string{"/test"},
		// log格式
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				params.ClientIP,
				params.TimeStamp.Format(time.RFC1123),
				params.Method,
				params.Path,
				params.Request.Proto,
				params.StatusCode,
				params.Latency,
				params.Request.UserAgent(),
				params.ErrorMessage,
			)
		},
	}
	r := gin.New()
	r.Use(gin.LoggerWithConfig(c))
	r.Use(gin.Recovery())
	v1 := r.Group("/v1")
	{
		v1.POST("/login",    u.Login)
		v1.POST("/register", u.Register)
	}
	return r
}
