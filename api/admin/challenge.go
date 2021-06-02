package apiAdmin

import (
	cfg "CTFgo/configs"
	"CTFgo/logs"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Challenge 定义一个题目
type Challenge struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Score       int    `json:"score"`
	Flag        string `json:"flag"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Tags        string `json:"tags"`
	Hints       string `json:"hints"`
	Visible     int    `json:"visible"` // 0表示隐藏，1表示可见
}

// NewChallenge 新增一个题目
func NewChallenge(c *gin.Context) {
	var request challengeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Request format wrong!"})
		return
	}

	if matched := checkCategory(request.Category); matched == false {
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong category!"})
		return
	}

	challenge := &Challenge{
		Name:        request.Name,
		Score:       request.Score,
		Flag:        request.Flag,
		Description: request.Description,
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

// EditChallenge 修改一个题目
func EditChallenge(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logs.WARNING("wrong id error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong id!"})
		return
	}

	var request challengeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Request format wrong!"})
		return
	}

	if !isChallengeExisted(int(id)) {
		c.JSON(400, gin.H{"code": 400, "msg": "Challenge does not exist"})
		return
	}

	if matched := checkCategory(request.Category); matched == false {
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

// DeleteChallenge 删除一个题目
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

// GetAllChallenges 获取所有题目
func GetAllChallenges(c *gin.Context) {
	var challenges []challengeResponse

	if err := getAllChallenges(&challenges); err != nil {
		logs.WARNING("get challenges error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get all challenges failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": challenges})
}

// GetChallengesByCategory 获取指定类别的题目
func GetChallengesByCategory(c *gin.Context) {
	category := c.Param("category")
	if matched := checkCategory(category); matched == false {
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong category!"})
		return
	}

	var challenges []challengeResponse
	if err := getChallengesByCategory(&challenges, category); err != nil {
		logs.WARNING("get challenges error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get challenges failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": challenges})
}

// addChallenge 操作数据库新增一个题目
func addChallenge(c *Challenge) error {
	command := "INSERT INTO challenge (name,score,flag,description,category,tags,hints,visible) VALUES (?,?,?,?,?,?,?,?);"
	res, err := db.Exec(command, c.Name, c.Score, c.Flag, c.Description, c.Category, c.Tags, c.Hints, c.Visible)
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

// updateChallenge 操作数据库更新一个题目
func updateChallenge(c *Challenge) error {
	command := "UPDATE challenge SET name=?, score=?, flag=?, description=?, category=?, tags=?, hints=?, visible=?  WHERE id=?;"
	res, err := db.Exec(command, c.Name, c.Score, c.Flag, c.Description, c.Category, c.Tags, c.Hints, c.Visible, c.ID)
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

// deleteChallenge 操作数据库删除一个题目
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

// isChallengeExisted 检查数据库中是否存在某条公告
func isChallengeExisted(id int) (exists bool) {
	command := "SELECT EXISTS(SELECT 1 FROM challenge WHERE id = ?);"
	if err := db.QueryRow(command, id).Scan(&exists); err != nil {
		logs.WARNING("query or scan error", err)
		return false
	}
	return exists
}

// getAllChallenges 操作数据库获取所有题目
func getAllChallenges(challenges *[]challengeResponse) error {
	command := "SELECT id, name, score, description, category, tags, hints FROM challenge WHERE visible=1;"
	rows, err := db.Query(command)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var c challengeResponse
		err = rows.Scan(&c.ID, &c.Name, &c.Score, &c.Description, &c.Category, &c.Tags, &c.Hints)
		if err != nil {
			return err
		}
		solverCount, err := getSolverCount(c.ID)
		if err != nil {
			return err
		}
		c.SolverCount = solverCount
		*challenges = append(*challenges, c)
	}
	return rows.Err()
}

// getChallengesByCategory 操作数据库获取所有题目
func getChallengesByCategory(challenges *[]challengeResponse, category string) error {
	command := "SELECT id, name, score, description, tags, hints FROM challenge WHERE visible=1 AND category=?;"
	rows, err := db.Query(command, category)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var c challengeResponse
		err = rows.Scan(&c.ID, &c.Name, &c.Score, &c.Description, &c.Tags, &c.Hints)
		if err != nil {
			return err
		}
		solverCount, err := getSolverCount(c.ID)
		if err != nil {
			return err
		}
		c.SolverCount = solverCount
		c.Category = category
		*challenges = append(*challenges, c)
	}
	return rows.Err()
}

// getSolverCount 操作数据库获取指定id题目的解出人数
func getSolverCount(id int) (count int, err error) {
	command := "SELECT COUNT(*) FROM solve WHERE cid = ?;"
	if err := db.QueryRow(command, id).Scan(&count); err != nil {
		logs.WARNING("query or scan error", err)
		return 0, err
	}
	return count, nil
}

// checkCategory检查类别是否正确
func checkCategory(c string) bool {
	for _, category := range cfg.ChallengeCategories {
		if category == c {
			return true
		}
	}
	return false
}
