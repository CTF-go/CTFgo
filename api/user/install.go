/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

import (
	cfg "CTFgo/configs"
	"CTFgo/logs"
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

//Install实现初始化数据库等功能。
func Install(c *gin.Context) {
	var json installRequest
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
	//判断是否存在CTFgo/databases/ctfgo.db，不存在则执行初始化，存在则不执行
	if _, err := os.Stat(cfg.DB_file); os.IsNotExist(err) {
		logs.INFO("ctfgo.db does not exist!")
		err := os.MkdirAll(cfg.DB_dir, 0777)
		db, err := sql.Open("sqlite3", cfg.DB_file)
		if err != nil {
			logs.ERROR("sqlite connect error: ", err)
		}
		logs.INFO("create ctfgo.db success!")
		//后面加其他表的创建
		table_sql := `
			CREATE TABLE "user" (
				"id"	INTEGER NOT NULL UNIQUE,
				"token"	TEXT NOT NULL UNIQUE,
				"username"	TEXT NOT NULL UNIQUE,
				"password"	TEXT NOT NULL,
				"email"	TEXT NOT NULL UNIQUE,
				"affiliation"	TEXT,
				"country"	TEXT,
				"hidden"	INTEGER NOT NULL DEFAULT 0,
				"banned"	INTEGER NOT NULL DEFAULT 1,
				"team_id"	INTEGER,
				"created"	TEXT NOT NULL,
				"role"	INTEGER NOT NULL DEFAULT 0,
				PRIMARY KEY("id" AUTOINCREMENT)
			);
			CREATE TABLE "scores" (
				"id"	INTEGER NOT NULL UNIQUE,
				"username"	TEXT NOT NULL UNIQUE,
				"scores"	INTEGER NOT NULL DEFAULT 0,
				PRIMARY KEY("id" AUTOINCREMENT)
			);
			CREATE TABLE "challenges" (
				"id"	INTEGER NOT NULL UNIQUE,
				"name"	TEXT NOT NULL UNIQUE,
				"score"	INTEGER NOT NULL,
				"flag"	TEXT NOT NULL,
				"max_attempts"	INTEGER NOT NULL DEFAULT 0,
				"description"	TEXT,
				"category"	TEXT,
				"tags"	TEXT,
				"hints"	TEXT,
				"requirements"	TEXT,
				"solves"	TEXT,
				PRIMARY KEY("id" AUTOINCREMENT)
			)
			`
		_, err = db.Exec(table_sql)
		if err != nil {
			logs.ERROR("create user table error:", err)
		}
		logs.INFO("create user table success!")
		sql_str := "INSERT INTO user (token,username,password,email,hidden,banned,created,role) VALUES (?,?,?,?,?,?,?,?);"
		_, err1 := db.Exec(sql_str, cfg.Token(), json.User, cfg.MD5(json.Passwd), json.Email, 1, 0, cfg.Timestamp(), 1)
		sql_str2 := "INSERT INTO scores (id,username,scores) VALUES (1,?,0);"
		_, err2 := db.Exec(sql_str2, json.User)
		if err1 != nil || err2 != nil {
			logs.ERROR("admin insert error", err)
		}
		logs.INFO("Administrator account [" + json.User + "]" + " register success!")
		//新建sessions文件夹
		err = os.MkdirAll(cfg.Session_dir, 0755)
		if err != nil {
			logs.ERROR("create sessions dir error:", err)
		}
		c.JSON(200, gin.H{"code": 200, "msg": "CTFgo installed successfully!"})
		return
	} else {
		logs.INFO("ctfgo.db already exist!")
		c.JSON(200, gin.H{"code": 200, "msg": "CTFgo has been installed, if you need to reinstall, please delete databases/ctfgo.db."})
		return
	}
}
