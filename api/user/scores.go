/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

import (
	cfg "CTFgo/configs"

	"github.com/gin-gonic/gin"
)

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
