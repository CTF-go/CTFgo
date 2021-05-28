package apiAdmin

import (
	cfg "CTFgo/configs"
	i "CTFgo/databases/init"
	"CTFgo/logs"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

var db *sql.DB = i.DB

// Notice 定义一个公告
type Notice struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int32  `json:"created_at"`
	UpdatedAt int32  `json:"updated_at"`
}

// NewNotice 新增一个公告
func NewNotice(c *gin.Context) {
	var request newNoticeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Request format wrong!"})
		return
	}

	notice := &Notice{
		Title:     request.Title,
		Content:   request.Content,
		CreatedAt: cfg.Timestamp(),
		UpdatedAt: cfg.Timestamp(),
	}
	if err := addNotice(notice); err != nil {
		logs.WARNING("add notice to database error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Add notice failure!"})
		return
	}

	logs.INFO(fmt.Sprintf("new notice [%s] added success!", notice.Title))
	c.JSON(200, gin.H{"code": 200, "msg": "Add notice success!"})
}

// EditNotice 修改一个公告
func EditNotice(c *gin.Context) {
	var request editNoticeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Request format wrong!"})
		return
	}

	if !isNoticeExisted(request.ID) {
		c.JSON(400, gin.H{"code": 400, "msg": "notice does not exist"})
		return
	}

	notice := &Notice{
		ID:        request.ID,
		Title:     request.Title,
		Content:   request.Content,
		UpdatedAt: cfg.Timestamp(),
	}
	if err := updateNotice(notice); err != nil {
		logs.WARNING("update notice error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Update notice failure!"})
		return
	}

	logs.INFO(fmt.Sprintf("update notice [%d] success!", request.ID))
	c.JSON(200, gin.H{"code": 200, "msg": "Update notice success!"})
}

// DeleteNotice 删除一个公告
func DeleteNotice(c *gin.Context) {
	var request deleteNoticeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Request format wrong!"})
		return
	}

	if !isNoticeExisted(request.ID) {
		c.JSON(400, gin.H{"code": 400, "msg": "notice does not exist"})
		return
	}

	if err := deleteNotice(request.ID); err != nil {
		logs.WARNING("delete notice error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Delete notice failure!"})
		return
	}

	logs.INFO(fmt.Sprintf("delete notice [id=%d] success!", request.ID))
	c.JSON(200, gin.H{"code": 200, "msg": "Delete notice success!"})
}

// GetAllNotices 获取所有的公告
func GetAllNotices(c *gin.Context) {
	var notices []Notice

	if err := getAllNotices(&notices); err != nil {
		logs.WARNING("get notices error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get all notices failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": notices})
}

// addNotice 操作数据库新增一个公告
func addNotice(b *Notice) error {
	command := "INSERT INTO notice (title,content,created_at,updated_at) VALUES (?,?,?,?);"
	res, err := db.Exec(command, b.Title, b.Content, b.CreatedAt, b.UpdatedAt)
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

// updateNotice 操作数据库更新一个公告
func updateNotice(b *Notice) error {
	command := "UPDATE notice SET title=?, content=?, updated_at=? WHERE id=?;"
	res, err := db.Exec(command, b.Title, b.Content, b.UpdatedAt, b.ID)
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

func getAllNotices(notices *[]Notice) error {
	command := "SELECT id, title, content, created_at, updated_at FROM notice;"
	rows, err := db.Query(command)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var b Notice
		err = rows.Scan(&b.ID, &b.Title, &b.Content, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return err
		}
		*notices = append(*notices, b)
	}
	return rows.Err()
}
