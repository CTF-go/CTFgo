package apiinit

import (
	u "CTFgo/api/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/login", u.Login)
		//v1.GET("/register", u.register)
	}
	return r
}
