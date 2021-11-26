package apiAdmin

import (
	. "CTFgo/api/types"
	u "CTFgo/api/user"
	"CTFgo/logs"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewChallenge 新增一个题目。
func NewChallenge(c *gin.Context) {
	var request ChallengeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Request format wrong!"})
		return
	}

	if matched := u.CheckCategory(request.Category); !matched {
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong category!"})
		return
	}

	challenge := &Challenge{
		Name:        request.Name,
		Score:       request.Score,
		Flag:        request.Flag,
		Description: request.Description,
		Attachment:  request.Attachment,
		Category:    request.Category,
		Tags:        request.Tags,
		Hints:       request.Hints,
		Visible:     request.Visible,
	}
	if err := addChallenge(challenge); err != nil {
		logs.WARNING("add challenge to database error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Add challenge failure!"})
		return
	}

	logs.INFO(fmt.Sprintf("new challenge [%s] added success!", challenge.Name))
	c.JSON(200, gin.H{"code": 200, "msg": "Add challenge success!"})
}

// EditChallenge 修改一个题目。
// TODO: 判断修改哪个值才修改，其他值为空则不变。
func EditChallenge(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logs.WARNING("wrong id error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong id!"})
		return
	}

	var request ChallengeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Request format wrong!"})
		return
	}

	if !isChallengeExisted(int(id)) {
		c.JSON(400, gin.H{"code": 400, "msg": "Challenge does not exist"})
		return
	}

	if matched := u.CheckCategory(request.Category); !matched {
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong category!"})
		return
	}

	challenge := &Challenge{
		ID:          int(id),
		Name:        request.Name,
		Score:       request.Score,
		Flag:        request.Flag,
		Description: request.Description,
		Category:    request.Category,
		Tags:        request.Tags,
		Hints:       request.Hints,
		Visible:     request.Visible,
	}
	if err := updateChallenge(challenge); err != nil {
		logs.WARNING("update challenge error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Update challenge failure!"})
		return
	}

	logs.INFO(fmt.Sprintf("update challenge [%d] success!", id))
	c.JSON(200, gin.H{"code": 200, "msg": "Update challenge success!"})
}

// DeleteChallenge 删除一个题目。
func DeleteChallenge(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logs.WARNING("wrong id error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong id!"})
		return
	}

	if !isChallengeExisted(int(id)) {
		c.JSON(400, gin.H{"code": 400, "msg": "Challenge does not exist"})
		return
	}

	if err := deleteChallenge(int(id)); err != nil {
		logs.WARNING("delete challenge error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Delete challenge failure!"})
		return
	}

	logs.INFO(fmt.Sprintf("delete challenge [id=%d] success!", id))
	c.JSON(200, gin.H{"code": 200, "msg": "Delete challenge success!"})
}

// MakeChallengeVisibleByID 修改单个题目为可见。
func MakeChallengeVisibleByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logs.WARNING("wrong id error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong id!"})
		return
	}

	if err := makeChallengeVisibleByID(int(id)); err != nil {
		logs.WARNING("make challenge visible error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "make challenge visible failure!"})
		return
	}

	logs.INFO(fmt.Sprintf("make challenge visible [%d] success!", id))
	c.JSON(200, gin.H{"code": 200, "msg": "Make challenge visible success!"})
}

// MakeAllChallengeVisible 批量修改所有题目为可见。
func MakeAllChallengeVisible(c *gin.Context) {
	if err := makeAllChallengeVisible(); err != nil {
		logs.WARNING("make challenge visible error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "make all challenge visible failure!"})
		return
	}

	logs.INFO("make all challenge visible success!")
	c.JSON(200, gin.H{"code": 200, "msg": "Make all challenge visible success!"})
}

// addChallenge 操作数据库新增一个题目。
func addChallenge(c *Challenge) error {
	// 使用逗号分隔字符串
	attachmentString := strings.Join(c.Attachment, ",")
	hintString := strings.Join(c.Hints, ",")
	command := "INSERT INTO challenge (name,score,flag,description,attachment,category,tags,hints,visible) VALUES (?,?,?,?,?,?,?,?,?);"
	res, err := db.Exec(command, c.Name, c.Score, c.Flag, c.Description, attachmentString, c.Category, c.Tags, hintString, c.Visible)
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

// updateChallenge 操作数据库更新一个题目。
func updateChallenge(c *Challenge) error {
	attachmentString := strings.Join(c.Attachment, ",")
	hintString := strings.Join(c.Hints, ",")
	command := "UPDATE challenge SET name=?, score=?, flag=?, description=?, attachment=?, category=?, tags=?, hints=?, visible=? WHERE id=?;"
	res, err := db.Exec(command, c.Name, c.Score, c.Flag, c.Description, attachmentString, c.Category, c.Tags, hintString, c.Visible, c.ID)
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

// deleteChallenge 操作数据库删除一个题目。
func deleteChallenge(id int) error {
	command := "DELETE FROM challenge WHERE id=?;"
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

// isChallengeExisted 检查数据库中是否存在某个题目。
func isChallengeExisted(id int) (exists bool) {
	command := "SELECT EXISTS(SELECT 1 FROM challenge WHERE id = ?);"
	if err := db.QueryRow(command, id).Scan(&exists); err != nil {
		logs.WARNING("query or scan error", err)
		return false
	}
	return exists
}

// makeChallengeVisibleByID 操作数据库更新一个题目为可见。
func makeChallengeVisibleByID(id int) error {
	command := "UPDATE challenge SET visible=1 WHERE id=?;"
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

// makeAllChallengeVisible 操作数据库更新s所有题目为可见。
func makeAllChallengeVisible() error {
	command := "UPDATE challenge SET visible=1;"
	res, err := db.Exec(command)
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
