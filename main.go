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

// exitfunc用于执行CTFgo退出前释放资源等一些操作。
func exitfunc() {
	logs.Save_log()
	fmt.Println("[EXIT] CTFgo has stopped")
	os.Exit(0)
}

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
	fullName := path.Join(_this_.path, name)
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

	// 刷新后展示相同页面
	r.StaticFile("/home", staticIndex.path)
	r.StaticFile("/users", staticIndex.path)
	r.StaticFile("/notices", staticIndex.path)
	r.StaticFile("/scoreboard", staticIndex.path)
	r.StaticFile("/profile", staticIndex.path)
	// r.StaticFile("/settings", staticIndex.path)
	r.StaticFile("/dashboard", staticIndex.path)

	// challenges子路由动态获取
	r.GET("/challenges/*type", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		indexHTML, _ := staticIndex.staticFS.ReadFile(staticIndex.path)
		c.Writer.Write(indexHTML)
		c.Writer.Header().Add("Accept", "text/html")
		c.Writer.Flush()
	})

	// 创建监听退出chan
	c := make(chan os.Signal)
	// 监听指定信号 ctrl+c kill，实现优雅退出
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
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
