package apiUser

import (
	. "CTFgo/api/types"
	"CTFgo/logs"

	"github.com/gin-gonic/gin"
)

// GetAllNotices 获取所有的公告
func GetAllNotices(c *gin.Context) {
	var notices []Notice

	if err := getAllNotices(&notices); err != nil {
		logs.WARNING("get notices error", err)
		c.JSON(200, gin.H{"code": 400, "msg": "Get all notices failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": notices})
}

func getAllNotices(notices *[]Notice) error {
	command := "SELECT id, title, content, created_at FROM notice ORDER BY created_at ASC;"
	rows, err := db.Query(command)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var b Notice
		err = rows.Scan(&b.ID, &b.Title, &b.Content, &b.CreatedAt)
		if err != nil {
			return err
		}
		*notices = append(*notices, b)
	}
	return rows.Err()
}
