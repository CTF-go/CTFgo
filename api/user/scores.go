/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

import (
	cfg "CTFgo/configs"
	"CTFgo/logs"

	"github.com/gin-gonic/gin"
)

// GetScoreByUserID 获取指定id用户得分。
func GetScoreByUserID(c *gin.Context) {
	var maxid string
	var score int
	id := c.Params.ByName("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "msg": "Get score error, need id!"})
		return
	}
	sql_str := "SELECT seq FROM sqlite_sequence WHERE name = 'score' LIMIT 1;"
	row := db.QueryRow(sql_str)
	row.Scan(&maxid)
	if !cfg.ID_verify(id) {
		c.JSON(400, gin.H{"code": 400, "msg": "Get score error!"})
		return
	}
	if !cfg.Num_compare(id, maxid) {
		c.JSON(400, gin.H{"code": 400, "msg": "Get score error!"})
		return
	}
	sql_str = "SELECT score FROM score WHERE id = ? LIMIT 1;"
	row = db.QueryRow(sql_str, id)
	row.Scan(&score)
	c.JSON(200, gin.H{"code": 200, "data": score})
}

// GetAllScores 按降序排列返回所有当前得分及对应用户和用户ID。
func GetAllScores(c *gin.Context) {
	var s Score
	var scores []Score
	sql_str := "SELECT * FROM score;"
	rows, err := db.Query(sql_str)
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
