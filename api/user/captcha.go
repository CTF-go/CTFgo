/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

import (
	"CTFgo/logs"
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

//Captcha 返回captcha图片的base64值。
func Captcha(c *gin.Context) {
	id := captcha.New()
	b64 := captcha_base64(id)
	if b64 == "" {
		c.JSON(400, gin.H{"code": 400, "msg": "Cannot get captcha!"})
		return
	}
	c.JSON(200, gin.H{"code": 200, "id": id, "data": b64})
}

//Captcha_verify 验证验证码id对应的验证码与用户输入的验证码正确与否。
func Captcha_verify(c *gin.Context) {
	var json captchaRequest

	//用ShouldBindJSON解析绑定传入的Json数据。
	if err := c.ShouldBindJSON(&json); err != nil {
		logs.WARNING("bindjson error: ", err)
		c.JSON(400, gin.H{"code": 400, "msg": "Missing parameters or format error!"})
		return
	}
	if !captcha.VerifyString(json.CaptchaID, json.Solution) {
		c.JSON(200, gin.H{"code": 400, "msg": "Captcha verify failed"})
		return
	} else {
		c.JSON(200, gin.H{"code": 200, "msg": "Captcha verify success"})
		return
	}
}

//captcha_base64 返回验证码图片的base64值。
func captcha_base64(id string) string {
	imgurl := "http://127.0.0.1:8081/v1/captcha/" + id + ".png"
	response, err := http.Get(imgurl)
	if err != nil || response.StatusCode != 200 {
		logs.WARNING("get captcha image error", err)
		return ""
	}
	img, err := ioutil.ReadAll(response.Body)
	imgb64 := base64.StdEncoding.EncodeToString([]byte(img))
	return imgb64
}

//Captcha_server 提供验证码图片.
func Captcha_server(c *gin.Context) {
	ServeHTTP(c.Writer, c.Request)
}

//Serve 是captcha包原生函数，移植方便gin使用。
func Serve(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		w.Header().Set("Content-Category", "image/png")
		captcha.WriteImage(&content, id, width, height)
	case ".wav":
		w.Header().Set("Content-Category", "audio/x-wav")
		captcha.WriteAudio(&content, id, lang)
	default:
		return captcha.ErrNotFound
	}

	if download {
		w.Header().Set("Content-Category", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}

//ServeHTTP 是captcha包原生函数，移植方便gin使用。
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir, file := path.Split(r.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || id == "" {
		http.NotFound(w, r)
		return
	}
	if r.FormValue("reload") != "" {
		captcha.Reload(id)
	}
	lang := strings.ToLower(r.FormValue("lang"))
	download := path.Base(dir) == "download"
	if Serve(w, r, id, ext, lang, download, captcha.StdWidth, captcha.StdHeight) == captcha.ErrNotFound {
		http.NotFound(w, r)
	}
}
