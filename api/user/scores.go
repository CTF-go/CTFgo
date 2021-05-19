/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

import (
	cfg "CTFgo/configs"
	"CTFgo/logs"

	"github.com/gin-gonic/gin"
)

//scores_struct 定义返回得分情况结构体。
type scores_struct struct {
	ID    int
	User  string
	Score int
}

//Specified_score 获取指定id用户得分。
func Specified_score(c *gin.Context) {
	var maxid string
	var scores int
	id := c.Params.ByName("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "msg": "Get score error, need id!"})
		return
	}
	sql_str := "SELECT seq FROM sqlite_sequence WHERE name = 'scores' LIMIT 1;"
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
	sql_str = "SELECT scores FROM scores WHERE id = ? LIMIT 1;"
	row = db.QueryRow(sql_str, id)
	row.Scan(&scores)
	c.JSON(200, gin.H{"code": 200, "data": scores})
}

//All_scores 按降序排列返回所有当前得分及对应用户和用户ID。
func All_scores(c *gin.Context) {
	var user scores_struct
	var users []scores_struct
	sql_str := "SELECT * FROM scores ORDER BY scores DESC;"
	rows, err := db.Query(sql_str)
	if err != nil {
		logs.WARNING("get all scores error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get all score error!"})
		return
	}
	// 循环读取数据
	for rows.Next() {
		rows.Scan(&user.ID, &user.User, &user.Score)
		users = append(users, user)
	}
	c.JSON(200, gin.H{"code": 200, "data": users})
}
