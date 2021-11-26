/*
Package apiInit实现初始化API接口功能。
*/
package apiInit

import (
	admin "CTFgo/api/admin"
	u "CTFgo/api/user"
	cfg "CTFgo/configs"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

// SetupAPI用于注册api，对应功能实现在CTFgo/api下。
func SetupAPI() *gin.Engine {
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DisableConsoleColor()
	err := os.MkdirAll(cfg.LOG_DIR, 0777)
	if err != nil {
		fmt.Printf("create logs dir error, err:%v\n", err)
	}
	// 创建run.log
	log_path, err := os.Create(cfg.CURRENT_LOG_PATH)
	if err != nil {
		fmt.Printf("create logs file error, err:%v\n", err)
	}
	// 将log输出到控制台和文件
	gin.DefaultWriter = io.MultiWriter(log_path, os.Stdout)
	c := gin.LoggerConfig{
		Output: gin.DefaultWriter,
		// 需要跳过记录log的API
		SkipPaths: []string{"/js/app.f2e52bbc.js", "/js/app.f2e52bbc.js.map"},
		// log格式
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("[GIN] [%s] [%s] - [%s] [%3d] %s %s %s \"%s\" %s\n",
				params.TimeStamp.Format("2006/01/02 15:04:05"),
				params.ClientIP,
				params.Method,
				params.StatusCode,
				params.Path,
				params.Request.Proto,
				params.Latency,
				params.Request.UserAgent(),
				params.ErrorMessage,
			)
		},
	}
	r := gin.New()
	r.Use(gin.LoggerWithConfig(c))
	r.Use(gin.Recovery())
	r.Use(Cors())

	api := r.Group("/v1")

	// 公共api，无需登陆即可访问
	public := api.Group("")
	{
		// CTFgo初始化
		public.POST("/install", u.Install)

		// 用户注册
		public.POST("/register", u.Register)
		// 用户登录
		public.POST("/login", u.Login)
		// 用户注销
		public.GET("/logout", u.Logout)

		// 获取验证码图片base64
		public.GET("/captcha", u.Captcha)

		// 获取指定id用户可公开信息
		public.GET("/info/:id", u.GetInfoByUserID)

		// 获取指定id用户分数
		public.GET("/score/:id", u.GetScoreByUserID)
		// 获取所有用户分数，降序排列
		public.GET("/scores/all", u.GetAllScores)

		// 获取所有公告
		public.GET("/notice/all", u.GetAllNotices)

	}

	// 普通用户api，需要用户登陆且Role=0才能访问
	personal := api.Group("/user")
	personal.Use(u.AuthRequired())
	{
		// 获取当前用户信息
		personal.GET("/session", u.Session)
		// 修改用户信息
		// personal.POST("/updateinfo", u.UpdateInfo)

		// 获取题目分类
		personal.GET("/categories", u.GetCategories)

		// 获取所有题目信息
		personal.GET("/challenges/all", u.GetAllChallenges)
		// 获取指定类别的题目信息
		personal.GET("/challenges/:category", u.GetChallengesByCategory)

		// 提交flag
		personal.POST("/submitflag", u.SubmitFlag)
		// 获取所有正确的flag提交记录
		personal.GET("/solves/all", u.GetAllSolves)
		// 获取指定用户正确的flag提交记录
		personal.GET("/solves/uid/:uid", u.GetSolvesByUid)
		// 获取指定题目正确的flag提交记录
		personal.GET("/solves/cid/:cid", u.GetSolvesByCid)
		// 获取当前用户正确flag提交记录（即解题记录）按时间从早到晚排序
		personal.GET("/solves/self", u.GetSelfSolves)

		// 获取当前用户分数、排名
		personal.GET("/score/self", u.GetSelfScoreAndRank)

		// 校内提交学号等信息接口
		personal.POST("/submit/studentinfo", u.SubmitStudentInfo)
		// 校外提交联系方式等信息接口
		personal.POST("/submit/othersinfo", u.SubmitOthersInfo)

		// 获取校内用户提交的相关信息
		personal.GET("/info/submit/self", u.GetStudentsAndOthersInfo)
	}

	// 管理员api，需要用户登陆且Role=1才能访问
	manager := api.Group("/admin")
	manager.Use(admin.AuthRequired())
	{
		// 创建新题目
		manager.POST("/challenge", admin.NewChallenge)
		// 更改题目
		manager.PUT("/challenge/:id", admin.EditChallenge)
		// 删除题目
		manager.DELETE("/challenge/:id", admin.DeleteChallenge)
		// 修改某个id的题目为可见
		manager.PUT("/challenge/visible/:id", admin.MakeChallengeVisibleByID)
		// 批量修改所有题目可见
		manager.PUT("/challenge/visible/all", admin.MakeAllChallengeVisible)

		// 创建新公告
		manager.POST("/notice", admin.NewNotice)
		// 删除公告
		manager.DELETE("/notice/:id", admin.DeleteNotice)

		// 获取所有提交记录
		manager.GET("/submissions/all", admin.GetAllSubmissions)
		// 获取指定用户的提交记录
		manager.GET("/submissions/uid/:uid", admin.GetSubmissionsByUid)
		// 获取指定题目的提交记录
		manager.GET("/submissions/cid/:cid", admin.GetSubmissionsByCid)
	}

	return r
}

// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin) //"http://172.20.10.10:8081")
			c.Header("Access-Control-Allow-Headers", "Content-Category, AccessToken, X-CSRF-Token, Authorization, Token, Content-Type")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Category")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.JSON(200, "ok")
		}
		// 处理请求
		c.Next()
	}
}
