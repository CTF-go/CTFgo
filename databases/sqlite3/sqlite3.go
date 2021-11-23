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

//Sqlite_conn实现连接sqlite3数据库，返回(*sql.DB)。
func Sqlite_conn() *sql.DB {
	db, err := sql.Open("sqlite3", c.DB_FILE)
	if err != nil {
		logs.ERROR("sqlite connect error: ", err)
	}
	return db
}
