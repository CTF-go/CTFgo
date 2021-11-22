package apiUser

import (
	. "CTFgo/api/types"
	cfg "CTFgo/configs"
	"CTFgo/logs"

	"github.com/gin-gonic/gin"
)

// GetAllChallenges 获取所有题目
func GetAllChallenges(c *gin.Context) {
	var challenges []ChallengeResponse

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
	if matched := checkCategory(category); !matched {
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong category!"})
		return
	}

	var challenges []ChallengeResponse
	if err := getChallengesByCategory(&challenges, category); err != nil {
		logs.WARNING("get challenges error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get challenges failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": challenges})
}

// getAllChallenges 操作数据库获取所有题目
func getAllChallenges(challenges *[]ChallengeResponse) error {
	command := "SELECT id, name, score, description, category, tags, hints FROM challenge WHERE visible=1;"
	rows, err := db.Query(command)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var c ChallengeResponse
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

// getChallengesByCategory 操作数据库获取指定类型题目
func getChallengesByCategory(challenges *[]ChallengeResponse, category string) error {
	command := "SELECT id, name, score, description, tags, hints FROM challenge WHERE visible=1 AND category=?;"
	rows, err := db.Query(command, category)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var c ChallengeResponse
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

// checkCategory检查类别是否正确
func checkCategory(c string) bool {
	for _, category := range cfg.ChallengeCategories {
		if category == c {
			return true
		}
	}
	return false
}
