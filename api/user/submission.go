package apiUser

import (
	//admin "CTFgo/api/admin"
	cfg "CTFgo/configs"
	"CTFgo/logs"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// Submission 表示一次flag提交记录
type Submission struct {
	ID          int    `json:"id"`
	UserID      int    `json:"uid"`
	ChallengeID int    `json:"cid"`
	Flag        string `json:"flag"`
	Time        int64  `json:"submitted_at"`
}

// Solve 表示一次正确的flag提交记录
type Solve struct {
	ID          int   `json:"id"`
	UserID      int   `json:"uid"`
	ChallengeID int   `json:"cid"`
	Time        int64 `json:"solved_at"`
}

// SubmitFlag 提交一个flag
func SubmitFlag(c *gin.Context) {
	var request submissionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Request format wrong!"})
		return
	}

	// 获取UserID
	session, err := Store.Get(c.Request, cfg.SessionID)
	if err != nil {
		c.JSON(200, gin.H{"code": 400, "msg": "Get CTFGOSESSID error"})
		return
	}
	var user User
	user, ok := session.Values["user"].(User)
	if !ok {
		c.JSON(200, gin.H{"code": 400, "msg": "No session"})
		return
	}

	// 检查题目是否存在
	if !isChallengeExisted(request.Cid) {
		c.JSON(400, gin.H{"code": 400, "msg": "Challenge does not exist!"})
		return
	}

	// Submission记录
	solvedTime := time.Now().Unix()
	submission := &Submission{
		UserID:      user.ID,
		ChallengeID: request.Cid,
		Flag:        request.Flag,
		Time:        solvedTime,
	}
	err = addSubmission(submission)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "Record submission failure!"})
		return
	}

	// 是否已经解出该题
	if hasAlreadySolved(user.ID, request.Cid) {
		c.JSON(400, gin.H{"code": 400, "msg": "Already solved!"})
		return
	}

	// 获取flag
	flag, err := getFlag(request.Cid)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "Get flag failure!"})
		return
	}
	// 对比flag
	if request.Flag != flag {
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong flag!"})
	} else {
		// Solve记录
		solve := &Solve{
			UserID:      user.ID,
			ChallengeID: request.Cid,
			Time:        solvedTime,
		}
		err = addSolve(solve)
		if err != nil {
			c.JSON(400, gin.H{"code": 400, "msg": "Record solve failure!"})
			return
		}
		// 加分数
		err = addScore(user.Username, request.Cid)
		if err != nil {
			c.JSON(400, gin.H{"code": 400, "msg": "Add Score failure!"})
			return
		}

		logs.INFO(fmt.Sprintf("[%s] user solved [%d].", user.Username, request.Cid))
		c.JSON(200, gin.H{"code": 200, "msg": "Correct flag!"})
	}
}

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

// GetAllSolves 获取所有正确的flag提交记录
func GetAllSolves(c *gin.Context) {
	var solves []Solve

	if err := getAllSolves(&solves); err != nil {
		logs.WARNING("get solves error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get all solves failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": solves})
}

// GetSolvesByUid 根据用户id获取正确的flag提交记录
func GetSolvesByUid(c *gin.Context) {
	uid, err := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		logs.WARNING("wrong uid error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong uid!"})
		return
	}

	if uid == 1 {
		c.JSON(400, gin.H{"code": 400, "msg": "Not allowed!"})
		return
	}

	var solves []Solve
	if err := getSolvesByUid(&solves, int(uid)); err != nil {
		logs.WARNING("get specified solves error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get specified solves failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": solves})
}

// GetSolvesByCid 根据题目id获取正确的flag提交记录
func GetSolvesByCid(c *gin.Context) {
	cid, err := strconv.ParseInt(c.Param("cid"), 10, 64)
	if err != nil {
		logs.WARNING("wrong cid error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong cid!"})
		return
	}

	var solves []Solve
	if err := getSolvesByCid(&solves, int(cid)); err != nil {
		logs.WARNING("get specified solves error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get specified solves failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": solves})
}

// hasAlreadySolved 检查某道题是否已经被某用户解出
func hasAlreadySolved(uid int, cid int) (exists bool) {
	command := "SELECT EXISTS(SELECT 1 FROM solve WHERE uid=? AND cid=?);"
	if err := db.QueryRow(command, uid, cid).Scan(&exists); err != nil {
		logs.WARNING("query or scan error", err)
		return false
	}
	return exists
}

// addSolve 操作数据库加入一条正确的flag提交记录
func addSolve(s *Solve) error {
	command := "INSERT INTO solve (uid, cid, submitted_at) VALUES (?,?,?);"
	res, err := db.Exec(command, s.UserID, s.ChallengeID, s.Time)
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

// addScore 操作数据库为指定用户增加某题的分数
func addScore(username string, cid int) error {
	var newScore int
	command := "SELECT score FROM challenge WHERE id=?"
	err := db.QueryRow(command, cid).Scan(&newScore)
	if err != nil {
		return err
	}

	command = "UPDATE score SET score=score+? WHERE username=?"
	res, err := db.Exec(command, newScore, username)
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

// addSubmission 操作数据库加入一条flag提交记录
func addSubmission(s *Submission) error {
	command := "INSERT INTO submission (uid, cid, flag, submitted_at) VALUES (?,?,?,?);"
	res, err := db.Exec(command, s.UserID, s.ChallengeID, s.Flag, s.Time)
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

// getFlag 根据题目id获取该题的flag
func getFlag(id int) (flag string, err error) {
	command := "SELECT flag FROM challenge WHERE id=?"
	if err := db.QueryRow(command, id).Scan(&flag); err != nil {
		return "", err
	}
	return flag, nil
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

// getAllSubmissions 操作数据库获取所有提交记录
func getAllSubmissions(submissions *[]Submission) error {
	command := "SELECT id, uid, cid, flag, submitted_at FROM submission;"
	rows, err := db.Query(command)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var s Submission
		err = rows.Scan(&s.ID, &s.UserID, &s.ChallengeID, &s.Flag, &s.Time)
		if err != nil {
			return err
		}
		*submissions = append(*submissions, s)
	}
	return rows.Err()
}

// getSubmissionsByUid 操作数据库根据uid获取提交记录
func getSubmissionsByUid(submissions *[]Submission, uid int) error {
	command := "SELECT id, uid, cid, flag, submitted_at FROM submission WHERE uid=?;"
	rows, err := db.Query(command, uid)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var s Submission
		err = rows.Scan(&s.ID, &s.UserID, &s.ChallengeID, &s.Flag, &s.Time)
		if err != nil {
			return err
		}
		*submissions = append(*submissions, s)
	}
	return rows.Err()
}

// getSubmissionsByCid 操作数据库根据cid获取提交记录
func getSubmissionsByCid(submissions *[]Submission, cid int) error {
	command := "SELECT id, uid, cid, flag, submitted_at FROM submission WHERE cid=?;"
	rows, err := db.Query(command, cid)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var s Submission
		err = rows.Scan(&s.ID, &s.UserID, &s.ChallengeID, &s.Flag, &s.Time)
		if err != nil {
			return err
		}
		*submissions = append(*submissions, s)
	}
	return rows.Err()
}

// getAllSolves 操作数据库获取所有正确的提交记录
func getAllSolves(solves *[]Solve) error {
	command := "SELECT id, uid, cid, submitted_at FROM solve WHERE uid != 1;"
	rows, err := db.Query(command)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var s Solve
		err = rows.Scan(&s.ID, &s.UserID, &s.ChallengeID, &s.Time)
		if err != nil {
			return err
		}
		*solves = append(*solves, s)
	}
	return rows.Err()
}

// getSolvesByUid 操作数据库根据用户id获取正确的flag提交记录
func getSolvesByUid(solves *[]Solve, uid int) error {
	command := "SELECT id, uid, cid, submitted_at FROM solve WHERE uid=?;"
	rows, err := db.Query(command, uid)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var s Solve
		err = rows.Scan(&s.ID, &s.UserID, &s.ChallengeID, &s.Time)
		if err != nil {
			return err
		}
		*solves = append(*solves, s)
	}
	return rows.Err()
}

// getSolvesByCid 操作数据库根据题目id获取正确的提交记录
func getSolvesByCid(solves *[]Solve, cid int) error {
	command := "SELECT id, uid, cid, submitted_at FROM solve WHERE cid=?;"
	rows, err := db.Query(command, cid)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var s Solve
		err = rows.Scan(&s.ID, &s.UserID, &s.ChallengeID, &s.Time)
		if err != nil {
			return err
		}
		*solves = append(*solves, s)
	}
	return rows.Err()
}
