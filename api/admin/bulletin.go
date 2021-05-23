package apiAdmin

import (
	u "CTFgo/api/user"
	i "CTFgo/databases/init"
	"CTFgo/logs"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

var db *sql.DB = i.DB

// NewBulletin 新增一个公告
func NewBulletin(c *gin.Context) {
	var request newBulletinRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Request format wrong!"})
		return
	}

	bulletin := &u.Bulletin{
		Title:   request.Title,
		Content: request.Content,
	}
	if err := addBulletin(bulletin); err != nil {
		logs.WARNING("add bulletin to database error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Add bulletin failure!"})
		return
	}

	logs.INFO(fmt.Sprintf("new bulletin [%s] added success!", bulletin.Title))
	c.JSON(200, gin.H{"code": 200, "msg": "Add bulletin success!"})
}

// EditBulletin 修改一个公告
func EditBulletin(c *gin.Context) {
	var request editBulletinRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Request format wrong!"})
		return
	}

	if !isBulletinExisted(request.ID) {
		c.JSON(400, gin.H{"code": 400, "msg": fmt.Sprintf("bulletin [%d] does not exsit", request.ID)})
		return
	}

	bulletin := &u.Bulletin{
		ID:      request.ID,
		Title:   request.Title,
		Content: request.Content,
	}
	if err := updateBulletin(bulletin); err != nil {
		logs.WARNING("update bulletin error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Update bulletin failure!"})
		return
	}

	logs.INFO(fmt.Sprintf("update bulletin [%d] success!", request.ID))
	c.JSON(200, gin.H{"code": 200, "msg": "Update bulletin success!"})
}

// DeleteBulletin 删除一个公告
func DeleteBulletin(c *gin.Context) {
	var request deleteBulletinRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Request format wrong!"})
		return
	}

	if !isBulletinExisted(request.ID) {
		c.JSON(400, gin.H{"code": 400, "msg": fmt.Sprintf("bulletin [%d] does not exsit", request.ID)})
		return
	}

	if err := deleteBulletin(request.ID); err != nil {
		logs.WARNING("delete bulletin error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Delete bulletin failure!"})
		return
	}

	logs.INFO(fmt.Sprintf("delete bulletin [id=%d] success!", request.ID))
	c.JSON(200, gin.H{"code": 200, "msg": "Delete bulletin success!"})
}

// GetAllBulletins 获取所有的公告
func GetAllBulletins(c *gin.Context) {
	var bulletins []u.Bulletin

	if err := getAllBulletins(&bulletins); err != nil {
		logs.WARNING("get bulletins error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get all bulletins failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": bulletins})
}

// addBulletin 操作数据库新增一个公告
func addBulletin(b *u.Bulletin) error {
	command := "INSERT INTO bulletin (title,content) VALUES (?,?);"
	res, err := db.Exec(command, b.Title, b.Content)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		err := errors.New("0 rows affected")
		return err
	}
	return nil
}

// updateBulletin 操作数据库更新一个公告
func updateBulletin(b *u.Bulletin) error {
	command := "UPDATE bulletin SET title=?, content=?, updated_at=(datetime('now','localtime')) WHERE id=?;"
	res, err := db.Exec(command, b.Title, b.Content, b.ID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		err := errors.New("0 rows affected")
		return err
	}
	return nil
}

// deleteBulletin 操作数据库删除一个公告
func deleteBulletin(id int) error {
	command := "DELETE FROM bulletin WHERE id=?;"
	res, err := db.Exec(command, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		err := errors.New("0 rows affected")
		return err
	}
	return nil
}

// isBulletinExisted 检查数据库中是否存在某条公告
func isBulletinExisted(id int) (exists bool) {
	command := "SELECT EXISTS(SELECT 1 FROM bulletin WHERE id = ?);"
	if err := db.QueryRow(command, id).Scan(&exists); err != nil {
		logs.WARNING("query or scan error", err)
		return false
	}
	return exists
}

func getAllBulletins(bulletins *[]u.Bulletin) error {
	command := "SELECT id, title, content, created_at, updated_at FROM bulletin;"
	rows, err := db.Query(command)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var b u.Bulletin
		err = rows.Scan(&b.ID, &b.Title, &b.Content, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return err
		}
		*bulletins = append(*bulletins, b)
	}
	return rows.Err()
}
