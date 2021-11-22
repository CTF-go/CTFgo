package apiUser

import (
	. "CTFgo/api/types"
	"CTFgo/logs"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllSolves 获取所有正确的解题情况记录
func GetAllSolves(c *gin.Context) {
	var solves []SolveResponse

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

	var solves []SolveResponse
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

	var solves []SolveResponse
	if err := getSolvesByCid(&solves, int(cid)); err != nil {
		logs.WARNING("get specified solves error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get specified solves failure!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": solves})
}

// getAllSolves 操作数据库获取所有正确的提交记录
func getAllSolves(solves *[]SolveResponse) error {
	command := "SELECT s.id, s.uid, s.cid, u.username, c.name, s.submitted_at FROM solve AS s, user AS u, challenge AS c WHERE u.id != 1 AND s.uid=u.id AND s.cid=c.id;"
	rows, err := db.Query(command)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var s SolveResponse
		err = rows.Scan(&s.ID, &s.Uid, &s.Cid, &s.Username, &s.ChallengeName, &s.SubmittedAt)
		if err != nil {
			return err
		}
		*solves = append(*solves, s)
	}
	return rows.Err()
}

// getSolvesByUid 操作数据库根据用户id获取正确的flag提交记录
func getSolvesByUid(solves *[]SolveResponse, uid int) error {
	command := "SELECT s.id, s.uid, s.cid, u.username, c.name, s.submitted_at FROM solve AS s, user AS u, challenge AS c WHERE s.uid=? AND u.id=s.uid AND c.id=s.cid;"
	rows, err := db.Query(command, uid)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var s SolveResponse
		err = rows.Scan(&s.ID, &s.Uid, &s.Cid, &s.Username, &s.ChallengeName, &s.SubmittedAt)
		if err != nil {
			return err
		}
		*solves = append(*solves, s)
	}
	return rows.Err()
}

// getSolvesByCid 操作数据库根据题目id获取正确的提交记录
func getSolvesByCid(solves *[]SolveResponse, cid int) error {
	command := "SELECT s.id, s.uid, s.cid, u.username, c.name, s.submitted_at FROM solve AS s, user AS u, challenge AS c WHERE s.cid=? AND u.id=s.uid AND c.id=s.cid;"
	rows, err := db.Query(command, cid)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var s SolveResponse
		err = rows.Scan(&s.ID, &s.Uid, &s.Cid, &s.Username, &s.ChallengeName, &s.SubmittedAt)
		if err != nil {
			return err
		}
		*solves = append(*solves, s)
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
