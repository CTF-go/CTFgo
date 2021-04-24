/*
Package sqlite3实现sqlite3数据库常用操作。
*/
package sqlite3

import (
	c "CTFgo/configs"
	"CTFgo/logs"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

//Create_db用于初次使用时创建数据库，生成ctfgo.db文件。
func Create_db() {
	//后面再写
}

//Sqlite_conn实现连接sqlite3数据库，返回(*sql.DB)。
func Sqlite_conn() *sql.DB {
	db, err := sql.Open("sqlite3", c.DB_file)
	if err != nil {
		logs.ERROR("sqlite connect error: ", err)
	}
	return db
}
