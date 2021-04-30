/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

import (
	cfg "CTFgo/configs"

	"github.com/gorilla/sessions"
)

//sessions存储于文件系统
var Store = sessions.NewFilesystemStore(cfg.Session_dir, []byte(cfg.Token()))
