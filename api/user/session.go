/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

import (
	cfg "CTFgo/configs"
	"encoding/gob"

	"github.com/gorilla/sessions"
)

//sessions存储于文件系统
var Store *sessions.FilesystemStore

func init() {
	Store = sessions.NewFilesystemStore(cfg.Session_dir, []byte(cfg.Token()))

	Store.Options = &sessions.Options{
		MaxAge: 24 * 60 * 60, // 1 day
		//Secure: true,
		//HttpOnly: true,
	}

	gob.Register(User{})
}
