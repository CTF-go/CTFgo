package apiAdmin

import (
	. "CTFgo/api/types"
	u "CTFgo/api/user"
	cfg "CTFgo/configs"

	"github.com/gin-gonic/gin"
)

// AuthRequired 用于管理员权限控制的中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var user User

		session, err := u.Store.Get(c.Request, cfg.SessionID)
		if err != nil {
			c.JSON(200, gin.H{"code": 400, "msg": "Get CTFGOSESSID error"})
			c.Abort()
			return
		}

		user, ok := session.Values["user"].(User)
		if !ok {
			c.JSON(200, gin.H{"code": 400, "msg": "No session"})
			c.Abort()
			return
		}

		if user.Role != 1 {
			c.JSON(200, gin.H{"code": 400, "msg": "Unauthorized access!"})
			c.Abort()
			return
		}

		c.Next()
	}
}
