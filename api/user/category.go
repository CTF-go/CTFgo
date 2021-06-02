package apiUser

import (
	cfg "CTFgo/configs"
	"github.com/gin-gonic/gin"
)

// GetCategories 获取题目分类
func GetCategories(c *gin.Context) {
	c.JSON(200, gin.H{"code": 200, "data": cfg.ChallengeCategories})
}
