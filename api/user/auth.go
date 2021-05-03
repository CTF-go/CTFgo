/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

import (
	cfg "CTFgo/configs"
	i "CTFgo/databases/init"
	"CTFgo/logs"
	"database/sql"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

var db *sql.DB = i.DB

//user定义用户结构体。
type Users struct {
	ID          int    `json:"id"`          //用户id，唯一，自增
	Token       string `json:"token"`       //用户token，唯一，API鉴权使用
	Username    string `json:"username"`    //用户名，唯一
	Password    string `json:"password"`    //用户密码md5值，md5(原密码）
	Email       string `json:"email"`       //邮箱，唯一
	Affiliation string `json:"affiliation"` //组织、战队或机构等，非必需
	Country     string `json:"country"`     //国家，非必需
	Hidden      int    `json:"hidden"`      //1：隐藏，0：显示，默认为0
	Banned      int    `json:"banned"`      //1：禁止，0：正常，默认为1，邮箱激活后为0
	Team_id     int    `json:"team_id"`     //队伍id，在团队模式下必须，个人模式非必需
	Created     string `json:"created"`     //用户注册时间，10位数时间戳
	Role        int    `json:"role"`        //1：管理员，0：普通用户，默认为0
}

// login_struct定义接收Login数据的结构体。
type login_struct struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	User   string `form:"username" json:"username" binding:"required"`
	Passwd string `form:"password" json:"password" binding:"required"`
}

// register_struct定义接收Login数据的结构体。
type register_struct struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	User   string `form:"username" json:"username" binding:"required"`
	Passwd string `form:"password" json:"password" binding:"required"`
	Email  string `form:"email" json:"email" binding:"required"`
}

//Login实现用户名或邮箱登录。
func Login(c *gin.Context) {
	var json login_struct
	var user Users

	//用ShouldBindJSON解析绑定传入的Json数据。
	if err := c.ShouldBindJSON(&json); err != nil {
		logs.WARNING("bindjson error: ", err)
		c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	//判断传入的是用户名还是邮箱，字符串中匹配到@字符则为邮箱，返回索引，匹配不到返回-1
	if strings.Index(json.User, "@") != -1 {
		//判断为邮箱，验证邮箱格式
		if !email_verify(json.User) {
			c.JSON(400, gin.H{"code": 400, "msg": "Email format error!"})
			return
		}
		//查询数据
		sql_str := "SELECT * FROM user WHERE email = ? LIMIT 1;"
		row := db.QueryRow(sql_str, json.User)
		row.Scan(&user.ID, &user.Token, &user.Username, &user.Password, &user.Email, &user.Affiliation, &user.Country, &user.Hidden, &user.Banned, &user.Team_id, &user.Created, &user.Role)
	} else {
		//判断为用户名，验证用户名格式
		if !name_verify(json.User) {
			c.JSON(400, gin.H{"code": 400, "msg": "Username format error!"})
			return
		}
		//查询数据
		sql_str := "SELECT * FROM user WHERE username = ? LIMIT 1;"
		row := db.QueryRow(sql_str, json.User)
		row.Scan(&user.ID, &user.Token, &user.Username, &user.Password, &user.Email, &user.Affiliation, &user.Country, &user.Hidden, &user.Banned, &user.Team_id, &user.Created, &user.Role)
	}

	//password进行md5加密
	json.Passwd = cfg.MD5(json.Passwd)
	//判断密码是否正确
	if json.Passwd != user.Password {
		logs.INFO("[" + json.User + "]" + " login error!")
		c.JSON(200, gin.H{"code": 400, "msg": "Login error!"})
		return
	}
	session, err := Store.Get(c.Request, "CTFGOSESSID")
	if err != nil {
		session.Save(c.Request, c.Writer)
		c.JSON(400, gin.H{"code": 400, "msg": "Get CTFGOSESSID error"})
		return
	}
	session.Options.HttpOnly = true
	session.Options.Secure = true
	session.Options.MaxAge = 86400 //86400秒:有效期一天
	session.Values["id"] = user.ID
	session.Values["token"] = user.Token
	session.Values["username"] = user.Username
	session.Values["email"] = user.Email
	session.Values["affiliation"] = user.Affiliation
	session.Values["country"] = user.Country
	session.Values["hidden"] = user.Hidden
	session.Values["banned"] = user.Banned
	session.Values["team_id"] = user.Team_id
	session.Values["created"] = user.Created
	session.Values["role"] = user.Role
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		logs.WARNING("can not save session:", err)
	}

	logs.INFO("[" + json.User + "]" + " login success!")
	c.JSON(200, gin.H{"code": 200, "username": user.Username, "msg": "Login success!"})
}

//Register实现注册功能。
func Register(c *gin.Context) {
	var json register_struct
	var user Users
	//用ShouldBindJSON解析绑定传入的Json数据。
	if err := c.ShouldBindJSON(&json); err != nil {
		logs.WARNING("bindjson error: ", err)
		c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	//限制传入用户名为中文、数字、大小写字母下划线和横杠，1到10位
	if !name_verify(json.User) {
		c.JSON(400, gin.H{"code": 400, "msg": "Username format error!"})
		return
	}
	//限制密码长度6到20位
	if !passwd_verify(json.Passwd) {
		c.JSON(400, gin.H{"code": 400, "msg": "Password format error!"})
		return
	}
	//限制传入邮箱符合格式
	if !email_verify(json.Email) {
		c.JSON(400, gin.H{"code": 400, "msg": "Email format error!"})
		return
	}
	//判断用户名是否已被使用
	if user_exists(user, json.User) {
		c.JSON(200, gin.H{"code": 1000, "msg": "Username has already been used!"})
		return
	}
	//判断邮箱是否已被使用
	if email_exists(user, json.Email) {
		c.JSON(200, gin.H{"code": 1001, "msg": "Email has already been used!"})
		return
	}
	//向数据库插入用户
	sql_str := "INSERT INTO user (token,username,password,email,created) VALUES (?,?,?,?,?);"
	res, err := db.Exec(sql_str, cfg.Token(), json.User, cfg.MD5(json.Passwd), json.Email, cfg.Timestamp())
	if err != nil {
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
	logs.INFO("[" + json.User + "]" + " register success!")
	c.JSON(200, gin.H{"code": 200, "msg": "Register success!"})
}

//Logout实现注销登录。
func Logout(c *gin.Context) {
	session, err := Store.Get(c.Request, "CTFGOSESSID")
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "Get CTFGOSESSID error"})
		return
	}
	if session.Values["username"] == nil {
		c.JSON(400, gin.H{"code": 400, "msg": "No session"})
		return
	}
	logout_user := session.Values["username"].(string)
	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)
	logs.INFO("[" + logout_user + "]" + " logout success!")
	c.JSON(200, gin.H{"code": 200, "msg": "Logout success!"})
}

//Session获取当前用户session信息，判断是否有效，即是否处于登录态。
func Session(c *gin.Context) {
	var user Users
	session, err := Store.Get(c.Request, "CTFGOSESSID")
	if err != nil {
		session.Save(c.Request, c.Writer)
		c.JSON(200, gin.H{"code": 400, "msg": "Get CTFGOSESSID error"})
		return
	}
	if session.Values["username"] == nil {
		session.Save(c.Request, c.Writer)
		c.JSON(200, gin.H{"code": 400, "msg": "No session"})
		return
	}
	user.ID = session.Values["id"].(int)
	user.Token = session.Values["token"].(string)
	user.Username = session.Values["username"].(string)
	user.Email = session.Values["email"].(string)
	user.Affiliation = session.Values["affiliation"].(string)
	user.Country = session.Values["country"].(string)
	user.Hidden = session.Values["hidden"].(int)
	user.Banned = session.Values["banned"].(int)
	user.Team_id = session.Values["team_id"].(int)
	user.Created = session.Values["created"].(string)
	user.Role = session.Values["role"].(int)
	c.JSON(200, gin.H{"code": 200, "msg": "here is the user info", "data": user})
}

//email_verify验证是否符合邮箱格式，返回true或false。
func email_verify(email string) bool {
	pattern := `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

//name_verify验证用户名是否符合中文数字字母下划线横杠，长度1到10位，返回true或false。
func name_verify(username string) bool {
	if !(utf8.RuneCountInString(username) > 0) || !(utf8.RuneCountInString(username) < 11) {
		return false
	}
	pattern := `^[-\w\p{Han}]+$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(username)
}

//passwd_verify验证密码是否符合长度6到20位，返回true或false。
func passwd_verify(password string) bool {
	if !(utf8.RuneCountInString(password) > 5) || !(utf8.RuneCountInString(password) < 21) {
		return false
	}
	return true
}

//user_exists判断用户名是否已经被占用，被占用返回true，未被占用则返回false。
func user_exists(user Users, username string) bool {
	sql_str := "SELECT username FROM user WHERE username = ? LIMIT 1;"
	err := db.QueryRow(sql_str, username).Scan(&user.Username)
	if err != nil {
		//数据库没有该用户名时，返回sql.ErrNoRows错误，即没有占用。
		if err != sql.ErrNoRows {
			//发生了一些真正的错误。
			logs.WARNING("an error occurred in the judgment process: ", err)
		}
		return false
	}
	//返回err为空时，则说明数据库存在该用户名，即用户名被占用。
	return true
}

//email_exists判断邮箱是否已经被占用，被占用返回true，未被占用则返回false。
func email_exists(user Users, email string) bool {
	sql_str := "SELECT email FROM user WHERE email = ? LIMIT 1;"
	err := db.QueryRow(sql_str, email).Scan(&user.Email)
	if err != nil {
		//数据库没有该邮箱时，返回sql.ErrNoRows错误，即没有占用。
		if err != sql.ErrNoRows {
			//发生了一些真正的错误。
			logs.WARNING("an error occurred in the judgment process: ", err)
		}
		return false
	}
	//返回err为空时，则说明数据库存在该邮箱，即邮箱被占用。
	return true
}
