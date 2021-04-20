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
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

var db *sql.DB = i.DB

//user定义用户结构体。
type user struct {
	id          int    //用户id，唯一，自增
	token       string //用户token，唯一，API鉴权使用
	username    string //用户名，唯一
	password    string //用户密码，md5(程序启动时生成的随机密钥+原密码）
	email       string //邮箱，唯一
	affiliation string //组织、战队或机构等，非必需
	country     string //国家，非必需
	hidden      int    //1：隐藏，0：显示，默认为0
	banned      int    //1：禁止，0：正常，默认为1，邮箱激活后为0
	team_id     int    //队伍id，在团队模式下必须，个人模式非必需
	created     string //用户注册时间，10位数时间戳
	role        int    //1：管理员，0：普通用户，默认为0
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

//Install实现初始化数据库，生成随机密钥等功能。
func Install(c *gin.Context) {
	//...
}

//Login实现用户名或邮箱登录。
func Login(c *gin.Context) {
	var json login_struct
	var user user

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
		sql_str := "SELECT password FROM user WHERE email = ?;"
		row := db.QueryRow(sql_str, json.User)
		row.Scan(&user.password)
	} else {
		//判断为用户名，验证用户名格式
		if !name_verify(json.User) {
			c.JSON(400, gin.H{"code": 400, "msg": "Username format error!"})
			return
		}
		//查询数据
		sql_str := "SELECT password FROM user WHERE username = ?;"
		row := db.QueryRow(sql_str, json.User)
		row.Scan(&user.password)
	}

	//password进行md5加密
	json.Passwd = cfg.MD5(json.Passwd)
	//判断密码是否正确
	if json.Passwd != user.password {
		logs.INFO("[" + json.User + "]" + " login error!")
		c.JSON(200, gin.H{"code": 400, "msg": "Login error!"})
		return
	}
	logs.INFO("[" + json.User + "]" + " login success!")
	c.JSON(200, gin.H{"code": 200, "msg": "Login success!"})
}

//Register实现注册功能。
func Register(c *gin.Context) {
	var json register_struct
	var user user
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

//Ping是一些功能测试接口。
func Ping(c *gin.Context) {
	fmt.Println("Beginning 20s")
	for i := 0; i <= 20; i++ {
		fmt.Println(i)
		time.Sleep(1 * 1e9)
	}
	fmt.Println("End of 20s")
}

//email_verify验证是否符合邮箱格式，返回true或false。
func email_verify(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

//name_verify验证用户名是否符合中文数字字母下划线横杠，长度1到10位，返回true或false。
func name_verify(username string) bool {
	if !(utf8.RuneCountInString(username) > 0) || !(utf8.RuneCountInString(username) < 11) {
		return false
	}
	pattern := `[-\w\p{Han}]+`
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
func user_exists(user user, username string) bool {
	sql_str := `SELECT username FROM user WHERE username = ?`
	err := db.QueryRow(sql_str, username).Scan(&user.username)
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
func email_exists(user user, email string) bool {
	sql_str := `SELECT email FROM user WHERE email = ?`
	err := db.QueryRow(sql_str, email).Scan(&user.email)
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
