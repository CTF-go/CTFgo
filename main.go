package main

import (
	i "CTFgo/api/init"
	"CTFgo/logs"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/gin-gonic/gin"
)

//exitfunc用于执行CTFgo退出前释放资源等一些操作。
func exitfunc() {
	logs.Save_log()
	fmt.Println("[EXIT] CTFgo has stopped")
	os.Exit(0)
}

// TODO: 使用 embed 嵌入静态资源。
//
//setup_frontend实现初始化前端路由。
// func setup_frontend() *gin.Engine {
// 	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
// 	gin.DisableConsoleColor()
// 	c := gin.LoggerConfig{
// 		Output: gin.DefaultWriter,
// 		// 需要跳过记录log的API
// 		SkipPaths: []string{"/test"},
// 		// log格式
// 		Formatter: func(params gin.LogFormatterParams) string {
// 			return fmt.Sprintf("[GIN] [%s] %s - \"%s %s %s %3d %s \"%s\" %s\"\n",
// 				params.TimeStamp.Format("2006/01/02 15:04:05"),
// 				params.ClientIP,
// 				params.Method,
// 				params.Path,
// 				params.Request.Proto,
// 				params.StatusCode,
// 				params.Latency,
// 				params.Request.UserAgent(),
// 				params.ErrorMessage,
// 			)
// 		},
// 	}
// 	r := gin.New()
// 	r.Use(gin.LoggerWithConfig(c))
// 	r.Use(gin.Recovery())
// 	r.Static("/css", cfg.Static_path+"/css")
// 	r.Static("/js", cfg.Static_path+"/js")
// 	r.Static("/img", cfg.Static_path+"/img")
// 	r.Static("/fonts", cfg.Static_path+"/fonts")

// 	r.StaticFile("/home", cfg.Static_path+"/index.html")
// 	r.StaticFile("/users", cfg.Static_path+"/index.html")
// 	r.StaticFile("/scoreboard", cfg.Static_path+"/index.html")
// 	r.StaticFile("/challenges", cfg.Static_path+"/index.html")

// 	r.GET("/", func(c *gin.Context) {
// 		c.Request.URL.Path = "/home"
// 		r.HandleContext(c)
// 	})
// 	return r
// }

//go:embed themes
var themesFS embed.FS

type StaticResource struct {
	// 静态资源
	staticFS embed.FS
	// 设置embed文件到静态资源的相对路径，也就是embed注释里的路径
	path string
}

// 静态资源被访问逻辑
func (_this_ *StaticResource) Open(name string) (fs.File, error) {
	var fullName string
	// fmt.Println(name)
	// if strings.Contains(name, `/`) {
	// 	fullName = path.Join(_this_.path, "static", name)
	// 	fmt.Println(1, fullName)
	// } else {
	fullName = path.Join(_this_.path, name)
	// fmt.Println(2, fullName)
	// }
	file, err := _this_.staticFS.Open(fullName)
	return file, err
}

//main执行启动CTFgo及关闭CTFgo相关操作。
func main() {
	r := i.SetupAPI()

	// 静态资源 themes
	staticIndex := &StaticResource{
		staticFS: themesFS,
		path:     "themes/index.html",
	}
	staticIcon := &StaticResource{
		staticFS: themesFS,
		path:     "themes/favicon.ico",
	}
	staticCss := &StaticResource{
		staticFS: themesFS,
		path:     "themes/css",
	}
	staticJs := &StaticResource{
		staticFS: themesFS,
		path:     "themes/js",
	}
	staticFonts := &StaticResource{
		staticFS: themesFS,
		path:     "themes/fonts",
	}
	staticImg := &StaticResource{
		staticFS: themesFS,
		path:     "themes/img",
	}

	r.StaticFS("/css/", http.FS(staticCss))
	r.StaticFS("/js/", http.FS(staticJs))
	r.StaticFS("/fonts/", http.FS(staticFonts))
	r.StaticFS("/img/", http.FS(staticImg))
	r.StaticFile("/favicon.ico", staticIcon.path)
	// 首页
	r.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		indexHTML, _ := staticIndex.staticFS.ReadFile(staticIndex.path)
		c.Writer.Write(indexHTML)
		c.Writer.Header().Add("Accept", "text/html")
		c.Writer.Flush()
	})

	// r.Any("/themes/*filepath", func(c *gin.Context) {
	// 	staticServer := http.FileServer(http.FS(themesFS))
	// 	staticServer.ServeHTTP(c.Writer, c.Request)
	// })

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
