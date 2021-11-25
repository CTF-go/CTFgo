package apiUser

import (
	. "CTFgo/api/types"
	cfg "CTFgo/configs"
	"CTFgo/logs"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/gin-gonic/gin"
)

// SubmitFlag 提交一个flag。
func SubmitFlag(c *gin.Context) {
	var request SubmissionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Request format wrong!"})
		return
	}

	// 获取UserID
	session, err := Store.Get(c.Request, cfg.SESSION_ID)
	if err != nil {
		c.JSON(200, gin.H{"code": 400, "msg": "Get CTFGOSESSID error"})
		return
	}
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
		IP:          c.ClientIP(),
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
		c.JSON(200, gin.H{"code": 400, "msg": "Wrong flag!"})
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
		err = addUserScore(user.Username, request.Cid)
		if err != nil {
			c.JSON(400, gin.H{"code": 400, "msg": "Add user score failure!"})
			return
		}
		// 题目动态分数
		reducedScore, err := editChallengeScore(request.Cid)
		if err != nil {
			c.JSON(400, gin.H{"code": 400, "msg": "Edit challenge score failure!"})
			return
		}
		err = updateUserScores(reducedScore, request.Cid)
		if err != nil {
			c.JSON(400, gin.H{"code": 400, "msg": "Update user scores failure!"})
		}

		logs.INFO(fmt.Sprintf("[%s] user solved [%d].", user.Username, request.Cid))
		c.JSON(200, gin.H{"code": 200, "msg": "Correct flag!"})
	}
}

// hasAlreadySolved 检查某道题是否已经被某用户解出。
func hasAlreadySolved(uid int, cid int) (exists bool) {
	command := "SELECT EXISTS(SELECT 1 FROM solve WHERE uid=? AND cid=?);"
	if err := db.QueryRow(command, uid, cid).Scan(&exists); err != nil {
		logs.WARNING("query or scan error", err)
		return false
	}
	return exists
}

// addSolve 操作数据库加入一条正确的flag提交记录。
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

// addUserScore 操作数据库为指定用户增加某题的分数。
func addUserScore(username string, cid int) error {
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

// updateUserScores 操作数据库更新解出用户的分数。
func updateUserScores(reducedScore, cid int) error {
	command := "UPDATE score SET score=score-? WHERE EXISTS(SELECT 1 FROM user,solve WHERE user.id=solve.uid AND score.username=user.username AND solve.cid=?);"
	_, err := db.Exec(command, reducedScore, cid)
	return err
}

// editChallengeScore 操作数据库修改指定题目增的动态分数。
func editChallengeScore(cid int) (reducedScore int, err error) {
	var currentScore int
	command := "SELECT score FROM challenge WHERE id=?;"
	if err := db.QueryRow(command, cid).Scan(&currentScore); err != nil {
		logs.WARNING("query challenge score error", err)
		return 0, err
	}

	solverCount, err := getSolverCount(cid)
	if err != nil {
		logs.WARNING("get solverCount error", err)
		return 0, err
	}
	// According to https://github.com/o-o-overflow/scoring-playground
	newScore := int(100 + (1000-100)/(1.0+float64(solverCount)*0.08*math.Log(float64(solverCount))))
	reducedScore = currentScore - newScore

	command = "UPDATE challenge SET score=? WHERE id=?;"
	res, err := db.Exec(command, newScore, cid)
	if err != nil {
		return 0, err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		err = errors.New("0 rows affected")
		return 0, err
	}

	return reducedScore, nil
}

// addSubmission 操作数据库加入一条flag提交记录。
func addSubmission(s *Submission) error {
	command := "INSERT INTO submission (uid, cid, ip, flag, submitted_at) VALUES (?,?,?,?,?);"
	res, err := db.Exec(command, s.UserID, s.ChallengeID, s.IP, s.Flag, s.Time)
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

// isChallengeExisted 检查数据库中是否存在某个题目。
func isChallengeExisted(id int) (exists bool) {
	command := "SELECT EXISTS(SELECT 1 FROM challenge WHERE id = ?);"
	if err := db.QueryRow(command, id).Scan(&exists); err != nil {
		logs.WARNING("query or scan error", err)
		return false
	}
	return exists
}
