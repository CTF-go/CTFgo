/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

import (
	cfg "CTFgo/configs"
	i "CTFgo/databases/init"
	"CTFgo/logs"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

var db *sql.DB = i.DB

// Login 实现用户名或邮箱登录
func Login(c *gin.Context) {
	var request loginRequest
	var user User

	//用ShouldBindJSON解析绑定传入的Json数据。
	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error: ", err)
		c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	//判断传入的是用户名还是邮箱，字符串中匹配到@字符则为邮箱，返回索引，匹配不到返回-1
	if strings.Index(request.Username, "@") != -1 {
		//判断为邮箱，验证邮箱格式
		if !checkEmail(request.Username) {
			c.JSON(400, gin.H{"code": 400, "msg": "Email format error!"})
			return
		}
		//查询数据
		sql_str := "SELECT * FROM user WHERE email = ? LIMIT 1;"
		row := db.QueryRow(sql_str, request.Username)
		row.Scan(&user.ID, &user.Token, &user.Username, &user.Password, &user.Email, &user.Affiliation, &user.Country, &user.Hidden, &user.Banned, &user.TeamID, &user.Created, &user.Role)
	} else {
		//判断为用户名，验证用户名格式
		if !checkUsername(request.Username) {
			c.JSON(400, gin.H{"code": 400, "msg": "Username format error!"})
			return
		}
		//查询数据
		sql_str := "SELECT * FROM user WHERE username = ? LIMIT 1;"
		row := db.QueryRow(sql_str, request.Username)
		row.Scan(&user.ID, &user.Token, &user.Username, &user.Password, &user.Email, &user.Affiliation, &user.Country, &user.Hidden, &user.Banned, &user.TeamID, &user.Created, &user.Role)
	}

	//password进行md5加密
	hashedPassword := cfg.MD5(request.Password)
	//判断密码是否正确
	logs.INFO(fmt.Sprintf("%s %s", hashedPassword, user.Password))
	if hashedPassword != user.Password {
		logs.INFO(fmt.Sprintf("[%s] login error!", user.Username))
		c.JSON(200, gin.H{"code": 400, "msg": "Login error!"})
		return
	}
	// 至此，身份认证完成

	// 设置session
	session, _ := Store.Get(c.Request, cfg.SessionID)
	user.Password = ""
	session.Values["user"] = user

	err := session.Save(c.Request, c.Writer)
	if err != nil {
		logs.WARNING("can not save session:", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Save CTFGOSESSID error"})
		return
	}

	logs.INFO("[" + request.Username + "]" + " login success!")
	c.JSON(200, gin.H{"code": 200, "username": user.Username, "msg": "Login success!"})
}

// Register 实现用户注册
func Register(c *gin.Context) {
	var request registerRequest
	var user User
	//用ShouldBindJSON解析绑定传入的Json数据。
	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error: ", err)
		c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	//限制传入用户名为中文、数字、大小写字母下划线和横杠，1到10位
	if !checkUsername(request.Username) {
		c.JSON(400, gin.H{"code": 400, "msg": "Username format error!"})
		return
	}
	//限制密码长度6到20位
	if !checkPassword(request.Password) {
		c.JSON(400, gin.H{"code": 400, "msg": "Password format error!"})
		return
	}
	//限制传入邮箱符合格式
	if !checkEmail(request.Email) {
		c.JSON(400, gin.H{"code": 400, "msg": "Email format error!"})
		return
	}
	//判断用户名是否已被使用
	if isNameExisted(user, request.Username) {
		c.JSON(200, gin.H{"code": 1000, "msg": "Username has already been used!"})
		return
	}
	//判断邮箱是否已被使用
	if isEmailExisted(user, request.Email) {
		c.JSON(200, gin.H{"code": 1001, "msg": "Email has already been used!"})
		return
	}
	//向数据库插入用户
	sql_str := "INSERT INTO user (token,username,password,email,created) VALUES (?,?,?,?,?);"
	res, err := db.Exec(sql_str, cfg.Token(), request.Username, cfg.MD5(request.Password), request.Email, cfg.Timestamp())
	sql_str2 := "INSERT INTO scores (username,scores) VALUES (?,0);"
	_, err2 := db.Exec(sql_str2, request.Username)
	if err != nil || err2 != nil {
		logs.WARNING("register insert error: ", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Register error!"})
		return
	}
	//id, _ := res.LastInsertId()
	affected, _ := res.RowsAffected()
	if affected == 0 {
		logs.WARNING("register insert error: ", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Register error!"})
		return
	}
	logs.INFO("[" + request.Username + "]" + " register success!")
	c.JSON(200, gin.H{"code": 200, "msg": "Register success!"})
}

// Logout 实现用户注销登陆
func Logout(c *gin.Context) {
	var user User

	session, err := Store.Get(c.Request, cfg.SessionID)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "Get CTFGOSESSID error"})
		return
	}
	user, ok := session.Values["user"].(User)
	if !ok {
		c.JSON(400, gin.H{"code": 400, "msg": "No session"})
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		logs.WARNING("can not save session:", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Save CTFGOSESSID error"})
		return
	}

	logs.INFO(fmt.Sprintf("[%s] logout success!", user.Username))
	c.JSON(200, gin.H{"code": 200, "msg": "Logout success!"})
}

// Session 获取当前用户session信息
func Session(c *gin.Context) {
	var user User

	// 默认在此之前已经通过了中间件的session权限验证
	session, _ := Store.Get(c.Request, cfg.SessionID)
	user = session.Values["user"].(User)

	c.JSON(200, gin.H{"code": 200, "data": user})
}

// UpdateInfo 更新用户信息
func UpdateInfo(c *gin.Context) {
	var request infoRequest
	var user User

	// 用ShouldBindJSON解析绑定传入的Json数据
	if err := c.ShouldBindJSON(&request); err != nil {
		logs.WARNING("bindjson error: ", err)
		c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// 默认在此之前已经通过了中间件的session权限验证
	session, _ := Store.Get(c.Request, cfg.SessionID)
	user = session.Values["user"].(User)

	if request.Username != "" && request.Username != user.Username {
		//限制传入用户名为中文、数字、大小写字母下划线和横杠，1到10位
		if !checkUsername(request.Username) {
			c.JSON(400, gin.H{"code": 400, "msg": "Username format error!"})
			return
		}
		//判断用户名是否已被使用
		if isNameExisted(user, request.Username) {
			c.JSON(200, gin.H{"code": 1000, "msg": "Username has already been used!"})
			return
		}
		//修改用户名
		sql_str := "UPDATE user SET username = ? where id = ?;"
		res, err := db.Exec(sql_str, request.Username, user.ID)
		if err != nil {
			logs.WARNING("update info error: ", err)
			c.JSON(400, gin.H{"code": 400, "msg": "Update info error!"})
			return
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			logs.WARNING("update info error: ", err)
			c.JSON(400, gin.H{"code": 400, "msg": "Update info error!"})
			return
		}

		logs.INFO(fmt.Sprintf("[%s] change username to [%s]", user.Username, request.Username))
		user.Username = request.Username
	}

	if request.Password != "" && request.Password != user.Password {
		//限制密码长度6到20位
		if !checkPassword(request.Password) {
			c.JSON(400, gin.H{"code": 400, "msg": "Password format error!"})
			return
		}
		//修改密码
		newPassword := cfg.MD5(request.Password)
		sql_str := "UPDATE user SET password = ? where id = ?;"
		res, err := db.Exec(sql_str, newPassword, user.ID)
		if err != nil {
			logs.WARNING("update info error: ", err)
			c.JSON(400, gin.H{"code": 400, "msg": "Update info error!"})
			return
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			logs.WARNING("update info error: ", err)
			c.JSON(400, gin.H{"code": 400, "msg": "Update info error!"})
			return
		}

		logs.INFO(fmt.Sprintf("[%s] change password successfully", user.Username))
		user.Password = newPassword
	}

	if request.Email != "" && request.Email != user.Email {
		//限制传入邮箱符合格式
		if !checkEmail(request.Email) {
			c.JSON(400, gin.H{"code": 400, "msg": "Email format error!"})
			return
		}
		//判断邮箱是否已被使用
		if isEmailExisted(user, request.Email) {
			c.JSON(200, gin.H{"code": 1001, "msg": "Email has already been used!"})
			return
		}
		//修改邮箱
		sql_str := "UPDATE user SET email = ? where id = ?;"
		res, err := db.Exec(sql_str, request.Email, user.ID)
		if err != nil {
			logs.WARNING("update info error: ", err)
			c.JSON(400, gin.H{"code": 400, "msg": "Update info error!"})
			return
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			logs.WARNING("update info error: ", err)
			c.JSON(400, gin.H{"code": 400, "msg": "Update info error!"})
			return
		}

		logs.INFO(fmt.Sprintf("[%s] change email from [%s] to [%s]", user.Username, user.Email, request.Email))
		user.Email = request.Email
	}

	// TODO: 队伍能随便改嘛？
	if request.Affiliation != "" && request.Affiliation != user.Affiliation {
		//限制传入名称为中文、数字、大小写字母下划线和横杠，1到10位
		if !checkUsername(request.Affiliation) {
			c.JSON(400, gin.H{"code": 400, "msg": "Affiliation format error!"})
			return
		}
		//修改Affiliation
		sql_str := "UPDATE user SET affiliation = ? where id = ?;"
		res, err := db.Exec(sql_str, request.Affiliation, user.ID)
		if err != nil {
			logs.WARNING("update info error: ", err)
			c.JSON(400, gin.H{"code": 400, "msg": "Update info error!"})
			return
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			logs.WARNING("update info error: ", err)
			c.JSON(400, gin.H{"code": 400, "msg": "Update info error!"})
			return
		}

		logs.INFO(fmt.Sprintf("[%s] change affiliation from [%s] to [%s]", user.Username, user.Affiliation, request.Affiliation))
		user.Affiliation = request.Affiliation
	}

	if request.Country != "" && request.Country != user.Country {
		//限制传入名称为中文、数字、大小写字母下划线和横杠，1到10位
		//暂定，等商量country存储格式后修改过滤
		if !checkUsername(request.Country) {
			c.JSON(400, gin.H{"code": 400, "msg": "Country format error!"})
			return
		}
		//修改Country
		sql_str := "UPDATE user SET country = ? where id = ?;"
		res, err := db.Exec(sql_str, request.Country, user.ID)
		if err != nil {
			logs.WARNING("update info error: ", err)
			c.JSON(400, gin.H{"code": 400, "msg": "Update info error!"})
			return
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			logs.WARNING("update info error: ", err)
			c.JSON(400, gin.H{"code": 400, "msg": "Update info error!"})
			return
		}

		logs.INFO(fmt.Sprintf("[%s] change country from [%s] to [%s]", user.Username, user.Country, request.Country))
		user.Country = request.Country
	}

	// 更新session
	session.Values["user"] = user
	err := session.Save(c.Request, c.Writer)
	if err != nil {
		logs.WARNING("can not save session:", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Save CTFGOSESSID error"})
		return
	}

	if request.Username == "" && request.Password == "" && request.Email == "" && request.Affiliation == "" && request.Country == "" {
		c.JSON(400, gin.H{"code": 400, "msg": "Nothing to be update!"})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "Update userinfo success!"})
}

// checkEmail 验证是否符合邮箱格式，返回true或false
func checkEmail(email string) bool {
	pattern := `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// checkUsername 验证用户名是否符合中文数字字母下划线横杠，长度1到10位，返回true或false
func checkUsername(username string) bool {
	if !(utf8.RuneCountInString(username) > 0) || !(utf8.RuneCountInString(username) < 11) {
		return false
	}
	pattern := `^[-\w\p{Han}]+$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(username)
}

// checkPassword 验证密码是否符合长度6到20位，返回true或false
func checkPassword(password string) bool {
	if !(utf8.RuneCountInString(password) > 5) || !(utf8.RuneCountInString(password) < 21) {
		return false
	}
	return true
}

// isNameExisted 判断用户名是否已经被占用，被占用返回true，未被占用则返回false
func isNameExisted(user User, username string) bool {
	sql_str := "SELECT username FROM user WHERE username = ? LIMIT 1;"
	err := db.QueryRow(sql_str, username).Scan(&user.Username)
	if err != nil {
		//数据库没有该用户名时，返回sql.ErrNoRows错误，即没有占用
		if err != sql.ErrNoRows {
			//发生了一些真正的错误
			logs.WARNING("an error occurred in the judgment process: ", err)
		}
		return false
	}
	//返回err为空时，则说明数据库存在该用户名，即用户名被占用
	return true
}

// isEmailExisted 判断邮箱是否已经被占用，被占用返回true，未被占用则返回false
func isEmailExisted(user User, email string) bool {
	sql_str := "SELECT email FROM user WHERE email = ? LIMIT 1;"
	err := db.QueryRow(sql_str, email).Scan(&user.Email)
	if err != nil {
		// 数据库没有该邮箱时，返回sql.ErrNoRows错误，即没有占用
		if err != sql.ErrNoRows {
			// 发生了一些真正的错误
			logs.WARNING("an error occurred in the judgment process: ", err)
		}
		return false
	}
	// 返回err为空时，则说明数据库存在该邮箱，即邮箱被占用
	return true
}

// AuthRequired 用于普通用户权限控制的中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User

		session, err := Store.Get(c.Request, cfg.SessionID)
		if err != nil {
			c.JSON(200, gin.H{"code": 400, "msg": "Get CTFGOSESSID error"})
			c.Abort()
			return
		}

		user, ok := session.Values["user"].(User)
		if !ok {
			c.JSON(200, gin.H{"code": 400, "msg": "No session"})
			c.Abort()
			return
		}

		if user.Role != 0 {
			c.JSON(200, gin.H{"code": 400, "msg": "Unauthorized access!"})
			c.Abort()
			return
		}

		c.Next()
	}
}
