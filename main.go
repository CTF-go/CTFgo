package main

import (
	i "CTFgo/api/init"
	"CTFgo/logs"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

//exitfunc用于执行CTFgo退出前释放资源等一些操作。
func exitfunc() {
	logs.Save_log()
	fmt.Println("CTFgo has stopped")
	os.Exit(0)
}

//main执行启动CTFgo及关闭CTFgo相关操作。
func main() {
	r := i.SetupRouter()
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill，实现优雅退出
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("exit: ", s)
				exitfunc()
			default:
				fmt.Println("other", s)
			}
		}
	}()
	if err := r.Run(); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
