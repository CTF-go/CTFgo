package apiAdmin

import (
	u "CTFgo/api/user"
	"github.com/gin-gonic/gin"
)

// NewChallenge 新增一个题目
func NewChallenge(c *gin.Context) {

}

// EditChallenge 修改一个题目
func EditChallenge(c *gin.Context) {

}

// DeleteChallenge 删除一个题目
func DeleteChallenge(c *gin.Context) {

}

// GetAllChallenges 获取所有题目
func GetAllChallenges(c *gin.Context) {

}

// addChallenge 操作数据库新增一个题目
func addChallenge(c *u.Challenge) error {
	return nil
}

// updateChallenge 操作数据库更新一个题目
func updateChallenge(c *u.Challenge) error {
	return nil
}

// deleteChallenge 操作数据库删除一个题目
func deleteChallenge(id int) error {
	return nil
}

// getAllChallenges 操作数据库获取所有题目
func getAllChallenges() error {
	return nil
}
