/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

import (
	. "CTFgo/api/types"
	cfg "CTFgo/configs"
	"CTFgo/logs"

	"github.com/gin-gonic/gin"
)

// GetScoreByUserID 获取指定id用户得分。
func GetScoreByUserID(c *gin.Context) {
	var score int
	id := c.Params.ByName("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "msg": "Need id!"})
		return
	}
	if !cfg.CheckID(id) {
		c.JSON(400, gin.H{"code": 400, "msg": "Format error!"})
		return
	}
	sql := "SELECT score FROM score WHERE id = ? LIMIT 1;"
	row := db.QueryRow(sql, id)
	err := row.Scan(&score)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "ID error!"})
		return
	}
	c.JSON(200, gin.H{"code": 200, "data": score})
}

// GetAllScores 按降序排列返回所有当前得分及对应用户和用户ID。
func GetAllScores(c *gin.Context) {
	var s ScoreResponse
	var scores []ScoreResponse
	sql := "SELECT * FROM score WHERE id != 1;"
	rows, err := db.Query(sql)
	if err != nil {
		logs.WARNING("get all scores error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get all score error!"})
		return
	}
	// 循环读取数据
	for rows.Next() {
		rows.Scan(&s.ID, &s.Username, &s.Score)
		scores = append(scores, s)
	}
	c.JSON(200, gin.H{"code": 200, "data": scores})
}
