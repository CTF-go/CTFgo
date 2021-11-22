package apiAdmin

import (
	. "CTFgo/api/types"
	"CTFgo/logs"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllSubmissions 获取所有flag提交记录
func GetAllSubmissions(c *gin.Context) {
	var submissions []Submission

	if err := getAllSubmissions(&submissions); err != nil {
		logs.WARNING("get submissions error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get all submissions failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": submissions})
}

// GetSubmissionsByUid 根据用户id获取flag提交记录
func GetSubmissionsByUid(c *gin.Context) {
	uid, err := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		logs.WARNING("wrong uid error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong uid!"})
		return
	}

	var submissions []Submission
	if err := getSubmissionsByUid(&submissions, int(uid)); err != nil {
		logs.WARNING("get specified submissions error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get specified submissions failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": submissions})
}

// GetSubmissionsByCid 根据题目id获取flag提交记录
func GetSubmissionsByCid(c *gin.Context) {
	cid, err := strconv.ParseInt(c.Param("cid"), 10, 64)
	if err != nil {
		logs.WARNING("wrong cid error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong cid!"})
		return
	}

	var submissions []Submission
	if err := getSubmissionsByCid(&submissions, int(cid)); err != nil {
		logs.WARNING("get specified submissions error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get specified submissions failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": submissions})
}

// getAllSubmissions 操作数据库获取所有提交记录
func getAllSubmissions(submissions *[]Submission) error {
	command := "SELECT id, uid, cid, ip, flag, submitted_at FROM submission;"
	rows, err := db.Query(command)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var s Submission
		err = rows.Scan(&s.ID, &s.UserID, &s.ChallengeID, &s.IP, &s.Flag, &s.Time)
		if err != nil {
			return err
		}
		*submissions = append(*submissions, s)
	}
	return rows.Err()
}

// getSubmissionsByUid 操作数据库根据uid获取提交记录
func getSubmissionsByUid(submissions *[]Submission, uid int) error {
	command := "SELECT id, uid, cid, ip, flag, submitted_at FROM submission WHERE uid=?;"
	rows, err := db.Query(command, uid)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var s Submission
		err = rows.Scan(&s.ID, &s.UserID, &s.ChallengeID, &s.IP, &s.Flag, &s.Time)
		if err != nil {
			return err
		}
		*submissions = append(*submissions, s)
	}
	return rows.Err()
}

// getSubmissionsByCid 操作数据库根据cid获取提交记录
func getSubmissionsByCid(submissions *[]Submission, cid int) error {
	command := "SELECT id, uid, cid, ip, flag, submitted_at FROM submission WHERE cid=?;"
	rows, err := db.Query(command, cid)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var s Submission
		err = rows.Scan(&s.ID, &s.UserID, &s.ChallengeID, &s.IP, &s.Flag, &s.Time)
		if err != nil {
			return err
		}
		*submissions = append(*submissions, s)
	}
	return rows.Err()
}
