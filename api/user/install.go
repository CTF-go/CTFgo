/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

import (
	. "CTFgo/api/types"
	cfg "CTFgo/configs"
	"CTFgo/logs"
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

//Install实现初始化数据库等功能。
func Install(c *gin.Context) {
	var request InstallRequest
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
	//判断是否存在CTFgo/databases/ctfgo.db，不存在则执行初始化，存在则不执行
	if _, err := os.Stat(cfg.DB_FILE); os.IsNotExist(err) {
		logs.INFO("ctfgo.db does not exist!")
		_ = os.MkdirAll(cfg.DB_DIR, 0777)
		db, err := sql.Open("sqlite3", cfg.DB_FILE)
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
				"website"	TEXT,
				"hidden"	INTEGER NOT NULL DEFAULT 1,
				"banned"	INTEGER NOT NULL DEFAULT 0,
				"team_id"	INTEGER DEFAULT 0,
				"created"	INTEGER NOT NULL,
				"role"	INTEGER NOT NULL DEFAULT 0,
				PRIMARY KEY("id" AUTOINCREMENT)
			);
			CREATE TABLE "score" (
				"id"	INTEGER NOT NULL UNIQUE,
				"username"	TEXT NOT NULL UNIQUE,
				"score"	INTEGER NOT NULL DEFAULT 0,
				PRIMARY KEY("id" AUTOINCREMENT)
			);
			CREATE TABLE "submission" (
				"id"	INTEGER NOT NULL UNIQUE,
				"uid"	INTEGER NOT NULL,
				"cid"	INTEGER NOT NULL,
				"flag"  TEXT NOT NULL,
				"ip"	TEXT NOT NULL,
				"submitted_at" INTEGER NOT NULL,
				PRIMARY KEY("id" AUTOINCREMENT)
			);
			CREATE TABLE "solve" (
			    "id"  INTEGER NOT NULL UNIQUE,
			    "uid" INTEGER NOT NULL,
			    "cid" INTEGER NOT NULL,
			    "submitted_at" INTEGER NOT NULL,
				PRIMARY KEY("id" AUTOINCREMENT)
			);
			CREATE TABLE "challenge" (
				"id"	INTEGER NOT NULL UNIQUE,
				"name"	TEXT NOT NULL UNIQUE,
				"score"	INTEGER NOT NULL,
				"flag"	TEXT,
				"description"	TEXT,
				"attachment"	TEXT,
				"category"	TEXT NOT NULL,
				"tags"	TEXT,
				"hints"	TEXT,
				"visible" INTEGER DEFAULT 0,
				PRIMARY KEY("id" AUTOINCREMENT)
			);
			CREATE TABLE "notice" (
			    "id"	INTEGER NOT NULL UNIQUE,
			    "title" TEXT NOT NULL,
			    "content" TEXT,
				"created_at" INTEGER NOT NULL,
				PRIMARY KEY("id" AUTOINCREMENT)
			);
			CREATE TABLE "category" (
				"id"	INTEGER NOT NULL UNIQUE,
				"category"	TEXT NOT NULL UNIQUE,
				PRIMARY KEY("id" AUTOINCREMENT)
			);
			BEGIN;
			INSERT INTO "category" VALUES (1, 'Web');
			INSERT INTO "category" VALUES (2, 'Pwn');
			INSERT INTO "category" VALUES (3, 'Reverse');
			INSERT INTO "category" VALUES (4, 'Crypto');
			INSERT INTO "category" VALUES (5, 'Misc');
			COMMIT;
			CREATE TABLE "studentinfo" (
				"team_id"     INTEGER NOT NULL,
				"username"    TEXT    NOT NULL,
				"student_id"  TEXT    NOT NULL UNIQUE,
				"qq" TEXT NOT NULL UNIQUE
			);
			CREATE TABLE "othersinfo" (
				"team_id"  INTEGER NOT NULL,
				"username" TEXT    NOT NULL,
				"email"    TEXT NOT NULL UNIQUE,
				"qq" TEXT NOT NULL UNIQUE
			);
			`
		_, err = db.Exec(table_sql)
		if err != nil {
			logs.ERROR("db init error:", err)
		}
		logs.INFO("db init success!")
		sql1 := "INSERT INTO user (token,username,password,email,affiliation,country,website,hidden,banned,team_id,created,role) VALUES (?,?,?,?,?,?,?,?,?,?,?,?);"
		_, err1 := db.Exec(sql1, cfg.Token(), request.Username, cfg.MD5(request.Password), request.Email, "", "", "", 1, 0, 0, cfg.Timestamp(), 1)
		sql2 := "INSERT INTO score (username,score) VALUES (?,0);"
		_, err2 := db.Exec(sql2, request.Username)
		// // --- for test purposes ---
		// sql1 = "INSERT INTO user (token,username,password,email,affiliation,country,website,hidden,banned,team_id,created,role) VALUES (?,?,?,?,?,?,?,?,?,?,?,?);"
		// _, err1 = db.Exec(sql1, cfg.Token(), "test", cfg.MD5("123456"), "test@gmail.com", "", "", "", 0, 0, 0, cfg.Timestamp(), 0)
		// sql2 = "INSERT INTO score (username,score) VALUES (?,0);"
		// _, err2 = db.Exec(sql2, "test")
		// // --- end ---
		if err1 != nil {
			logs.ERROR("admin insert error", err1)
		}
		if err2 != nil {
			logs.ERROR("admin insert error", err2)
		}
		logs.INFO("Administrator account [" + request.Username + "]" + " register success!")

		//新建sessions文件夹
		err = os.MkdirAll(cfg.SESSION_DIR, 0755)
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
