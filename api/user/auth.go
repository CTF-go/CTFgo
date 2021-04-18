/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

import (
	cfg "CTFgo/configs"
	i "CTFgo/databases/init"
	"CTFgo/logs"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type user struct {
	id          int    //用户id，唯一，自增
	token       string //用户token，唯一，API鉴权使用
	username    string //用户名，唯一
	password    string //用户密码，md5(程序启动时生成的随机密钥+原密码）
	email       string //邮箱，唯一
	affiliation string //组织、战队或机构等
	country     string //国家
	hidden      int    //1：隐藏，0：显示
	banned      int    //1：禁止，0：正常
	team_id     int    //队伍id
	created     string //用户注册时间，10位数时间戳
	role        int    //1：管理员，0：普通用户
}

// Login_struct定义接收Login数据的结构体。
type login_struct struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	User   string `form:"username" json:"username" binding:"required"`
	Passwd string `form:"password" json:"password" binding:"required"`
}

//先写通过用户名登录简单的
func Login(c *gin.Context) {
	var json login_struct
	var user user

	//用ShouldBindJSON解析绑定传入的Json数据。
	if err := c.ShouldBindJSON(&json); err != nil {
		logs.WARNING("bindjson error: ", err)
		c.JSON(400, gin.H{"code": "400", "msg": err.Error()})
		return
	}
	//password进行md5加密
	json.Passwd = cfg.MD5(json.Passwd)

	//查询数据
	db := i.DB
	sql_str := "SELECT password FROM user WHERE username = ?;"
	row := db.QueryRow(sql_str, json.User)
	row.Scan(&user.password)

	// 判断用户名密码是否正确
	if json.Passwd != user.password {
		logs.INFO("[" + json.User + "]" + " login error!")
		c.JSON(400, gin.H{"code": "400", "msg": "login error!"})
		return
	}
	logs.INFO("[" + json.User + "]" + " login success!")
	c.JSON(200, gin.H{"code": "200", "msg": "login success!"})
}

func Register(c *gin.Context) {
	var json login_struct
	if err := c.ShouldBindJSON(&json); err != nil {
		logs.INFO("success")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
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
