/*
Package apiInit实现初始化API接口功能。
*/
package apiInit

import (
	u "CTFgo/api/user"
	c "CTFgo/configs"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

// SetupRouter用于注册api，对应功能实现在CTFgo/api下。
func SetupRouter() *gin.Engine {
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DisableConsoleColor()
	err := os.MkdirAll(c.Log_dir, 0777)
	if err != nil {
		fmt.Printf("create logs dir error, err:%v\n", err)
	}
	// 创建run.log
	log_path, _ := os.Create(c.Current_log_path)
	// 将log输出到控制台和文件
	gin.DefaultWriter = io.MultiWriter(log_path, os.Stdout)
	c := gin.LoggerConfig{
		Output: gin.DefaultWriter,
		// 需要跳过记录log的API
		SkipPaths: []string{"/test"},
		// log格式
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("[GIN] [%s] %s - \"%s %s %s %3d %s \"%s\" %s\"\n",
				params.TimeStamp.Format("2006/01/02 15:04:05"),
				params.ClientIP,
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
	r.Use(cors())
	v1 := r.Group("/v1")
	{
		v1.POST("/login", u.Login)
		v1.POST("/register", u.Register)
		v1.GET("/ping", u.Ping)
	}
	return r
}

//暂时跨域。
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
		}
		// 处理请求
		c.Next()
	}
}
