package apiInit

import (
	u "CTFgo/api/user"

	"github.com/gin-gonic/gin"
)

/*
注册api，对应功能实现在CTFgo/api/下
v1是路由组，启动后访问127.0.0.1:8080/v1/login
*/
func SetupRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.POST("/login", u.Login)
		v1.POST("/register", u.Register)
	}
	return r
}
