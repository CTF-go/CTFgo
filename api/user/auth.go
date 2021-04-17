/*
Package apiUser包实现用户和其他普通API接口功能。
*/
package apiUser

import (
	"CTFgo/logs"

	"github.com/gin-gonic/gin"
)

// Login_struct定义接收Login数据的结构体。
type Login_struct struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	User   string `form:"username" json:"username" binding:"required"`
	Passwd string `form:"password" json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var form Login_struct
	// Bind()默认解析并绑定form格式
	// 根据请求头中content-type自动推断
	if err := c.Bind(&form); err != nil {
		logs.WARNING("warning", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 判断用户名密码是否正确
	//等一下实现sqlite，先起个架构
	if form.User != "root" || form.Passwd != "admin" {
		logs.INFO("login success!", nil)
		c.JSON(400, gin.H{"status": "304"})
		return
	}
	c.JSON(200, gin.H{"status": "200"})
}

func Register(c *gin.Context) {
	var form Login_struct
	// Bind()默认解析并绑定form格式
	// 根据请求头中content-type自动推断
	if err := c.Bind(&form); err != nil {
		logs.ERROR("warning", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
}
