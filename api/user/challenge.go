package apiUser

import (
	. "CTFgo/api/types"
	cfg "CTFgo/configs"
	"CTFgo/logs"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetAllChallenges 获取所有题目。
func GetAllChallenges(c *gin.Context) {
	var challenges []ChallengeResponse

	if err := getAllChallenges(c, &challenges); err != nil {
		logs.WARNING("get challenges error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get all challenges failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": challenges})
}

// GetChallengesByCategory 获取指定类别的题目。
func GetChallengesByCategory(c *gin.Context) {
	category := c.Param("category")
	if matched := CheckCategory(category); !matched {
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong category!"})
		return
	}

	var challenges []ChallengeResponse
	if err := getChallengesByCategory(c, &challenges, category); err != nil {
		logs.WARNING("get challenges error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get challenges failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": challenges})
}

// getAllChallenges 操作数据库获取所有题目。
func getAllChallenges(c *gin.Context, challenges *[]ChallengeResponse) error {
	var attachmentString, hints string
	sql := "SELECT id, name, score, description, attachment, category, tags, hints FROM challenge WHERE visible=1;"
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var challenge ChallengeResponse
		err = rows.Scan(&challenge.ID, &challenge.Name, &challenge.Score, &challenge.Description, &attachmentString, &challenge.Category, &challenge.Tags, &hints)
		if err != nil {
			return err
		}
		// 解析为切片
		challenge.Attachment = strings.Split(attachmentString, ",")
		challenge.Hints = strings.Split(hints, ",")

		solverCount, err := getSolverCount(challenge.ID)
		if err != nil {
			return err
		}
		challenge.SolverCount = solverCount
		session, err := Store.Get(c.Request, cfg.SESSION_ID)
		if err != nil {
			c.JSON(200, gin.H{"code": 400, "msg": "Get CTFGOSESSID error"})
			return err
		}
		user, ok := session.Values["user"].(User)
		if !ok {
			c.JSON(200, gin.H{"code": 400, "msg": "No session"})
			return errors.New("no session")
		}
		challenge.IsSolved = hasAlreadySolved(user.ID, challenge.ID)
		*challenges = append(*challenges, challenge)
	}
	return rows.Err()
}

// getChallengesByCategory 操作数据库获取指定类型题目。
func getChallengesByCategory(c *gin.Context, challenges *[]ChallengeResponse, category string) error {
	var attachmentString, hints string
	sql := "SELECT id, name, score, description, attachment, tags, hints FROM challenge WHERE visible=1 AND category=?;"
	rows, err := db.Query(sql, category)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var challenge ChallengeResponse
		err = rows.Scan(&challenge.ID, &challenge.Name, &challenge.Score, &challenge.Description, &attachmentString, &challenge.Tags, &hints)
		if err != nil {
			return err
		}
		// 解析为切片
		challenge.Attachment = strings.Split(attachmentString, ",")
		challenge.Hints = strings.Split(hints, ",")

		solverCount, err := getSolverCount(challenge.ID)
		if err != nil {
			return err
		}
		challenge.SolverCount = solverCount
		challenge.Category = category
		session, err := Store.Get(c.Request, cfg.SESSION_ID)
		if err != nil {
			c.JSON(200, gin.H{"code": 400, "msg": "Get CTFGOSESSID error"})
			return err
		}
		user, ok := session.Values["user"].(User)
		if !ok {
			c.JSON(200, gin.H{"code": 400, "msg": "No session"})
			return errors.New("no session")
		}
		challenge.IsSolved = hasAlreadySolved(user.ID, challenge.ID)
		*challenges = append(*challenges, challenge)
	}
	return rows.Err()
}
