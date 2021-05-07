package main

import (
	i "CTFgo/api/init"
	u "CTFgo/api/user"
	cfg "CTFgo/configs"
	"CTFgo/logs"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

//exitfunc用于执行CTFgo退出前释放资源等一些操作。
func exitfunc() {
	logs.Save_log()
	fmt.Println("[EXIT] CTFgo has stopped")
	os.Exit(0)
}

//setup_frontend实现初始化前端路由。
func setup_frontend() *gin.Engine {
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DisableConsoleColor()
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
	r.Static("/css", cfg.Static_path+"/css")
	r.Static("/js", cfg.Static_path+"/js")
	r.Static("/img", cfg.Static_path+"/img")
	r.Static("/fonts", cfg.Static_path+"/fonts")

	r.StaticFile("/home", cfg.Static_path+"/index.html")
	r.StaticFile("/users", cfg.Static_path+"/index.html")
	r.StaticFile("/scoreboard", cfg.Static_path+"/index.html")
	r.StaticFile("/challenges", cfg.Static_path+"/index.html")

	r.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/home"
		r.HandleContext(c)
	})

	//输出验证码图片
	r.GET("/captcha/:img", u.Captcha_server)
	return r
}

//main执行启动CTFgo及关闭CTFgo相关操作。
func main() {
	//前端路由
	front_router := setup_frontend()
	//Listen and Server in 0.0.0.0:8080
	go front_router.Run(":8080")

	r := i.SetupAPI()
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill，实现优雅退出
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("[EXIT] Exit by", s)
				exitfunc()
			default:
				fmt.Println("[EXIT] ", s)
			}
		}
	}()
	if err := r.Run(); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}
	//Listen and Server in 0.0.0.0:8081
	r.Run(":8081")
}
