/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

import (
	. "CTFgo/api/types"
	cfg "CTFgo/configs"

	"encoding/gob"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// sessions 存储于文件系统
var Store *sessions.FilesystemStore

func init() {
	Store = sessions.NewFilesystemStore(cfg.SESSION_DIR, securecookie.GenerateRandomKey(32))

	Store.Options = &sessions.Options{
		Domain: "",
		Path:   "/",
		MaxAge: 24 * 60 * 60, // 1 day
		// SameSite: http.SameSiteNoneMode,
		Secure:   false,
		HttpOnly: false,
	}

	gob.Register(User{})
}
