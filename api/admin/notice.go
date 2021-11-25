package apiAdmin

import (
	. "CTFgo/api/types"
	cfg "CTFgo/configs"
	"CTFgo/logs"
	"errors"
	"fmt"

	"strconv"

	"github.com/gin-gonic/gin"
)

// NewNotice 新增一个公告
func NewNotice(c *gin.Context) {
	var request NoticeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Request format wrong!"})
		return
	}

	notice := &Notice{
		Title:     request.Title,
		Content:   request.Content,
		CreatedAt: cfg.Timestamp(),
	}
	if err := addNotice(notice); err != nil {
		logs.WARNING("add notice to database error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Add notice failure!"})
		return
	}

	logs.INFO(fmt.Sprintf("new notice [%s] added success!", notice.Title))
	c.JSON(200, gin.H{"code": 200, "msg": "Add notice success!"})
}

// DeleteNotice 删除一个公告
func DeleteNotice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logs.WARNING("wrong id error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong id!"})
		return
	}

	if !isNoticeExisted(int(id)) {
		c.JSON(400, gin.H{"code": 400, "msg": "Notice does not exist!"})
		return
	}

	if err := deleteNotice(int(id)); err != nil {
		logs.WARNING("delete notice error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Delete notice failure!"})
		return
	}

	logs.INFO(fmt.Sprintf("delete notice [id=%d] success!", id))
	c.JSON(200, gin.H{"code": 200, "msg": "Delete notice success!"})
}

// addNotice 操作数据库新增一个公告
func addNotice(b *Notice) error {
	command := "INSERT INTO notice (title,content,created_at) VALUES (?,?,?);"
	res, err := db.Exec(command, b.Title, b.Content, b.CreatedAt)
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

// deleteNotice 操作数据库删除一个公告
func deleteNotice(id int) error {
	command := "DELETE FROM notice WHERE id=?;"
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

// isNoticeExisted 检查数据库中是否存在某条公告
func isNoticeExisted(id int) (exists bool) {
	command := "SELECT EXISTS(SELECT 1 FROM notice WHERE id = ?);"
	if err := db.QueryRow(command, id).Scan(&exists); err != nil {
		logs.WARNING("query or scan error", err)
		return false
	}
	return exists
}
