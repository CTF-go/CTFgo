package apiUser

import (
	"CTFgo/logs"

	"github.com/gin-gonic/gin"
)

// GetCategories 获取所有题目分类。
func GetCategories(c *gin.Context) {
	var categories []string

	if err := getAllCategories(&categories); err != nil {
		logs.WARNING("get categories error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get categories failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": categories})
}

// getAllCategories 操作数据库所有题目分类。
func getAllCategories(categories *[]string) error {
	sql := "SELECT category FROM category;"
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var category string
		err = rows.Scan(&category)
		if err != nil {
			return err
		}
		*categories = append(*categories, category)
	}
	return rows.Err()
}

// CheckCategory 检查类别是否正确
func CheckCategory(c string) bool {
	var categories []string
	if err := getAllCategories(&categories); err != nil {
		logs.WARNING("get categories error", err)
		return false
	}

	for _, category := range categories {
		if category == c {
			return true
		}
	}
	return false
}
