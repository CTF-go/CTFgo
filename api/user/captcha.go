/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

import (
	"CTFgo/logs"
	"bytes"
	"encoding/base64"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

//Captcha 返回captcha图片的base64值。
func Captcha(c *gin.Context) {
	id := captcha.New()
	b64 := captchaBase64(id)
	if b64 == "" {
		c.JSON(400, gin.H{"code": 400, "msg": "Cannot get captcha!"})
		return
	}
	c.JSON(200, gin.H{"code": 200, "id": id, "data": b64})
}

// captchaVerify 验证验证码id对应的验证码与用户输入的验证码正确与否。
func captchaVerify(id string, solution string) bool {
	if !captcha.VerifyString(id, solution) {
		return false
	} else {
		return true
	}
}

// captchaBase64 返回验证码图片的base64值，验证码图片长240宽80。
func captchaBase64(id string) string {
	var content bytes.Buffer
	err := captcha.WriteImage(&content, id, 240, 80)
	if err != nil {
		logs.WARNING("write captcha image error:", err)
		return ""
	}
	imgb64 := base64.StdEncoding.EncodeToString(content.Bytes())
	return imgb64
}
