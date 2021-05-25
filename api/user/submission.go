package apiUser

import (
	//admin "CTFgo/api/admin"
	cfg "CTFgo/configs"
	"CTFgo/logs"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// Submission 表示一次flag提交记录
type Submission struct {
	ID          int
	UserID      int
	ChallengeID int
	Flag        string
	Time        int64
}

// Solve 表示一次正确的flag提交记录
type Solve struct {
	ID          int
	UserID      int
	ChallengeID int
	Time        int64
}

// SubmitFlag 提交一个flag
func SubmitFlag(c *gin.Context) {
	var request submissionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Request format wrong!"})
		return
	}

	// 检查题目是否存在
	if !isChallengeExisted(request.ChallengeID) {
		c.JSON(400, gin.H{"code": 400, "msg": "Challenge does not exist!"})
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

	// Submission记录
	solvedTime := time.Now().Unix()
	submission := &Submission{
		UserID:      user.ID,
		ChallengeID: request.ChallengeID,
		Flag:        request.Flag,
		Time:        solvedTime,
	}
	err = recordSubmission(submission)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "Record submission failure!"})
		return
	}

	// 获取flag
	flag, err := getFlag(request.ChallengeID)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "Get Flag failure!"})
		return
	}
	// 对比flag
	if request.Flag != flag {
		c.JSON(400, gin.H{"code": 400, "msg": "Wrong Flag!"})
	} else {
		// Solve记录
		solve := &Solve{
			UserID:      user.ID,
			ChallengeID: request.ChallengeID,
			Time:        solvedTime,
		}
		err = recordSolve(solve)
		if err != nil {
			c.JSON(400, gin.H{"code": 400, "msg": "Record solve failure!"})
			return
		}
		// 加分数
		err = addScore(user.Username, request.ChallengeID)
		if err != nil {
			c.JSON(400, gin.H{"code": 400, "msg": "Add Score failure!"})
		}

		logs.INFO(fmt.Sprintf("[%s] user solved [%d].", user.Username, request.ChallengeID))
		c.JSON(200, gin.H{"code": 200, "msg": "Correct Flag!"})
	}
}

// TODO: implement all functions below

func GetAllSubmissions() {

}

func GetSubmissionsByUid() {

}

func GetSubmissionsByCid() {

}

func GetAllSolves() {

}

func GetSolvesByUid() {

}

func GetSolvesByCid() {

}

func recordSolve(s *Solve) error {
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

func recordSubmission(s *Submission) error {
	command := "INSERT INTO submission (uid, cid, Flag, submitted_at) VALUES (?,?,?,?);"
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

func getFlag(id int) (flag string, err error) {
	command := "SELECT Flag FROM challenge WHERE id=?"
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
