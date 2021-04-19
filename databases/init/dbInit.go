/*
Package dbInit实现初始化数据库接口功能。
*/
package dbInit

import (
	"CTFgo/databases/sqlite3"
	"database/sql"
)

var DB *sql.DB

func init() {
	//从configs获取用户选择的数据库，调用不同的数据库接口，后面再实现。
	DB = sqlite3.Sqlite_conn()
}
