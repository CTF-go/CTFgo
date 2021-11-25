package apiUser

import (
	. "CTFgo/api/types"
	cfg "CTFgo/configs"
	"CTFgo/logs"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// SubmitStudentInfo 实现校内用户提交学号等信息接口。
func SubmitStudentInfo(c *gin.Context) {
	var request SubmitStudentInfoRequest
	var count = 0

	//用ShouldBindJSON解析绑定传入的Json数据。
	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error: ", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Bind json error!"})
		return
	}

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
	// 首先判断该队伍是否已经提交过信息
	if isUserInfoExisted(user.ID) {
		c.JSON(400, gin.H{"code": 400, "msg": "Info has already submitted!"})
		return
	}

	// 判断学生1信息是否为空，学生1必须，其他学生信息可为空，但必须按顺序
	if !isStudentEmpty(&request.Student1) {
		if !checkStudentAll(&request.Student1) {
			c.JSON(200, gin.H{"code": 2001, "msg": "Student1 info has format error!"})
			return
		}
		count++
	} else {
		c.JSON(200, gin.H{"code": 2000, "msg": "Student1 info is empty!"})
		return
	}

	// 判断学生2信息是否为空
	if !isStudentEmpty(&request.Student2) {
		if !checkStudentAll(&request.Student2) {
			c.JSON(200, gin.H{"code": 2002, "msg": "Student2 info has format error!"})
			return
		}
		count++
	}
	// 判断学生3信息是否为空
	if !isStudentEmpty(&request.Student3) {
		if !checkStudentAll(&request.Student3) {
			c.JSON(200, gin.H{"code": 2003, "msg": "Student3 info has format error!"})
			return
		}
		count++
	}
	// 判断学生4信息是否为空
	if !isStudentEmpty(&request.Student4) {
		if !checkStudentAll(&request.Student4) {
			c.JSON(200, gin.H{"code": 2004, "msg": "Student4 info has format error!"})
			return
		}
		count++
	}

	// 向数据库插入用户
	var command string
	var res sql.Result
	switch count {
	case 1:
		command = "INSERT INTO studentinfo (team_id, username, student_id, qq) VALUES (?,?,?,?)"
		res, err = db.Exec(command, user.ID, request.Student1.Username, request.Student1.StudentID, request.Student1.QQ)
	case 2:
		command = "INSERT INTO studentinfo (team_id, username, student_id, qq) VALUES (?,?,?,?),(?,?,?,?)"
		res, err = db.Exec(command, user.ID, request.Student1.Username, request.Student1.StudentID, request.Student1.QQ, user.ID, request.Student2.Username, request.Student2.StudentID, request.Student2.QQ)
	case 3:
		command = "INSERT INTO studentinfo (team_id, username, student_id, qq) VALUES (?,?,?,?),(?,?,?,?),(?,?,?,?)"
		res, err = db.Exec(command, user.ID, request.Student1.Username, request.Student1.StudentID, request.Student1.QQ, user.ID, request.Student2.Username, request.Student2.StudentID, request.Student2.QQ, user.ID, request.Student3.Username, request.Student3.StudentID, request.Student3.QQ)
	case 4:
		command = "INSERT INTO studentinfo (team_id, username, student_id, qq) VALUES (?,?,?,?),(?,?,?,?),(?,?,?,?),(?,?,?,?)"
		res, err = db.Exec(command, user.ID, request.Student1.Username, request.Student1.StudentID, request.Student1.QQ, user.ID, request.Student2.Username, request.Student2.StudentID, request.Student2.QQ, user.ID, request.Student3.Username, request.Student3.StudentID, request.Student3.QQ, user.ID, request.Student4.Username, request.Student4.StudentID, request.Student4.QQ)
	}
	if err != nil {
		logs.WARNING("submit student info insert error: ", err)
		c.JSON(200, gin.H{"code": 2010, "msg": "submit insert error!"})
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		err := errors.New("0 rows affected")
		logs.WARNING("submit student info insert error: ", err)
		c.JSON(400, gin.H{"code": 400, "msg": "submit insert error!"})
		return
	}
	err = disableHiddenStatus(user.ID)
	if err != nil {
		logs.WARNING("update hidden status error: ", err)
		c.JSON(400, gin.H{"code": 400, "msg": "update hidden status error!"})
	}
	logs.INFO(fmt.Sprintf("[%d] submit %d students info success!", user.ID, count))
	c.JSON(200, gin.H{"code": 200, "msg": "Submit student info success!"})
}

// isUserInfoExisted 检查数据库中是否存在某队伍的信息。
func isUserInfoExisted(id int) (exists bool) {
	command := "SELECT EXISTS(SELECT 1 FROM studentinfo WHERE team_id = ?) OR EXISTS(SELECT 1 FROM othersinfo WHERE team_id = ?);"
	if err := db.QueryRow(command, id, id).Scan(&exists); err != nil {
		logs.WARNING("query or scan error", err)
		return false
	}
	return exists
}

// isStudentEmpty 返回学生x信息是否为空。
func isStudentEmpty(stu *StudentInfo) bool {
	return (stu.Username == "" || stu.StudentID == "" || stu.QQ == "")
}

// checkStudentID 判断学号是否符合格式。
func checkStudentID(username string) bool {
	if !(utf8.RuneCountInString(username) > 4) || !(utf8.RuneCountInString(username) < 21) {
		return false
	}
	pattern := `^[\w]+$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(username)
}

// checkQQ 判断QQ号是否符合格式。
func checkQQ(username string) bool {
	if !(utf8.RuneCountInString(username) > 1) || !(utf8.RuneCountInString(username) < 16) {
		return false
	}
	pattern := `^[\d]+$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(username)
}

// checkStudentAll 判断所有格式问题。
func checkStudentAll(stu *StudentInfo) bool {
	// 限制传入用户名为中文、数字、大小写字母下划线和横杠，1到10位
	if !checkUsername(stu.Username) {
		return false
	}
	// 限制学号为数字大小写字母，长度5到20
	if !checkStudentID(stu.StudentID) {
		return false
	}
	// 限制QQ号为纯数字，长度2到15
	if !checkQQ(stu.QQ) {
		return false
	}
	return true
}

// SubmitOthersInfo 实现校外用户提交信息接口。
func SubmitOthersInfo(c *gin.Context) {
	var request SubmitOthersInfoRequest
	var count = 0

	//用ShouldBindJSON解析绑定传入的Json数据。
	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error: ", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Bind json error!"})
		return
	}

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
	// 首先判断该队伍是否已经提交过信息
	if isUserInfoExisted(user.ID) {
		c.JSON(400, gin.H{"code": 400, "msg": "Info has already submitted!"})
		return
	}

	// 判断用户1信息是否为空，用户1必须，其他用户信息可为空，但必须按顺序
	if !isOthersEmpty(&request.Others1) {
		if !checkOthersAll(&request.Others1) {
			c.JSON(200, gin.H{"code": 2001, "msg": "Others1 info has format error!"})
			return
		}
		count++
	} else {
		c.JSON(200, gin.H{"code": 2000, "msg": "Others1 info is empty!"})
		return
	}

	// 判断用户2信息是否为空
	if !isOthersEmpty(&request.Others2) {
		if !checkOthersAll(&request.Others2) {
			c.JSON(200, gin.H{"code": 2002, "msg": "Others2 info has format error!"})
			return
		}
		count++
	}
	// 判断用户3信息是否为空
	if !isOthersEmpty(&request.Others3) {
		if !checkOthersAll(&request.Others3) {
			c.JSON(200, gin.H{"code": 2003, "msg": "Others3 info has format error!"})
			return
		}
		count++
	}
	// 判断用户4信息是否为空
	if !isOthersEmpty(&request.Others4) {
		if !checkOthersAll(&request.Others4) {
			c.JSON(200, gin.H{"code": 2004, "msg": "Others4 info has format error!"})
			return
		}
		count++
	}

	// 向数据库插入用户
	var command string
	var res sql.Result
	switch count {
	case 1:
		command = "INSERT INTO othersinfo (team_id, username, email, qq) VALUES (?,?,?,?)"
		res, err = db.Exec(command, user.ID, request.Others1.Username, request.Others1.Email, request.Others1.QQ)
	case 2:
		command = "INSERT INTO othersinfo (team_id, username, email, qq) VALUES (?,?,?,?),(?,?,?,?)"
		res, err = db.Exec(command, user.ID, request.Others1.Username, request.Others1.Email, request.Others1.QQ, user.ID, request.Others2.Username, request.Others2.Email, request.Others2.QQ)
	case 3:
		command = "INSERT INTO othersinfo (team_id, username, email, qq) VALUES (?,?,?,?),(?,?,?,?),(?,?,?,?)"
		res, err = db.Exec(command, user.ID, request.Others1.Username, request.Others1.Email, request.Others1.QQ, user.ID, request.Others2.Username, request.Others2.Email, request.Others2.QQ, user.ID, request.Others3.Username, request.Others3.Email, request.Others3.QQ)
	case 4:
		command = "INSERT INTO othersinfo (team_id, username, email, qq) VALUES (?,?,?,?),(?,?,?,?),(?,?,?,?),(?,?,?,?)"
		res, err = db.Exec(command, user.ID, request.Others1.Username, request.Others1.Email, request.Others1.QQ, user.ID, request.Others2.Username, request.Others2.Email, request.Others2.QQ, user.ID, request.Others3.Username, request.Others3.Email, request.Others3.QQ, user.ID, request.Others4.Username, request.Others4.Email, request.Others4.QQ)
	}
	if err != nil {
		logs.WARNING("submit others info insert error: ", err)
		c.JSON(200, gin.H{"code": 2010, "msg": "submit insert error!"})
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		err := errors.New("0 rows affected")
		logs.WARNING("submit student info insert error: ", err)
		c.JSON(400, gin.H{"code": 400, "msg": "submit insert error!"})
		return
	}
	err = disableHiddenStatus(user.ID)
	if err != nil {
		logs.WARNING("update hidden status error: ", err)
		c.JSON(400, gin.H{"code": 400, "msg": "update hidden status error!"})
	}
	logs.INFO(fmt.Sprintf("[%d] submit %d others info success!", user.ID, count))
	c.JSON(200, gin.H{"code": 200, "msg": "Submit others info success!"})
}

// isOthersEmpty 返回用户x信息是否为空。
func isOthersEmpty(others *OthersInfo) bool {
	return (others.Username == "" || others.Email == "" || others.QQ == "")
}

// checkOthersAll 判断所有格式问题。
func checkOthersAll(others *OthersInfo) bool {
	// 限制传入用户名为中文、数字、大小写字母下划线和横杠，1到10位
	if !checkUsername(others.Username) {
		return false
	}
	// 限制学号为数字大小写字母，长度5到20
	if !checkEmail(others.Email) {
		return false
	}
	// 限制QQ号为纯数字，长度2到15
	if !checkQQ(others.QQ) {
		return false
	}
	return true
}

// GetStudentsAndOthersInfo 获取用户提交过的信息，两个表一起查。
func GetStudentsAndOthersInfo(c *gin.Context) {
	var userinfoSlice []StudentsOrOthersInfoResponse
	var status = 0 // 0: 还未提交信息，1: 已提交校内，2: 已提交校外

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
	if err := getStudentsAndOthersInfo(&userinfoSlice, user.ID); err != nil {
		logs.WARNING("get submitted info error", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Get submitted info failure!"})
		return
	}
	if len(userinfoSlice) != 0 {
		if !checkEmail(userinfoSlice[0].IDOrEmail) {
			status = 1
		} else {
			status = 2
		}
	}
	c.JSON(200, gin.H{"code": 200, "status": status, "data": userinfoSlice})
}

// getStudentsAndOthersInfo 操作数据库根据team_id获取用户信息。
func getStudentsAndOthersInfo(userinfoSlice *[]StudentsOrOthersInfoResponse, uid int) error {
	command := "SELECT username, student_id, qq FROM studentinfo WHERE team_id = ? UNION SELECT username, email, qq FROM othersinfo WHERE team_id = ?;"
	rows, err := db.Query(command, uid, uid)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var userinfo StudentsOrOthersInfoResponse
		err = rows.Scan(&userinfo.Username, &userinfo.IDOrEmail, &userinfo.QQ)
		if err != nil {
			return err
		}
		*userinfoSlice = append(*userinfoSlice, userinfo)
	}
	return rows.Err()
}

func disableHiddenStatus(uid int) error {
	command := "UPDATE user SET hidden = 0 where id = ?;"
	res, err := db.Exec(command, uid)
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
