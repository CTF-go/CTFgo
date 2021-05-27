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
	err := os.MkdirAll(cfg.Log_dir, 0777)
	if err != nil {
		fmt.Printf("create logs dir error, err:%v\n", err)
	}
	// 创建run.log
	log_path, err := os.Create(cfg.Current_log_path)
	if err != nil {
		fmt.Printf("create logs file error, err:%v\n", err)
	}
	// 将log输出到控制台和文件
	gin.DefaultWriter = io.MultiWriter(log_path, os.Stdout)
	c := gin.LoggerConfig{
		Output: gin.DefaultWriter,
		// 需要跳过记录log的API
		SkipPaths: []string{"/test"},
		// log格式
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("[GIN] [%s] %s - \"%s %s %s %3d %s \"%s\" %s\"\n",
				params.TimeStamp.Format("2006/01/02 15:04:05"),
				params.ClientIP,
				params.Method,
				params.Path,
				params.Request.Proto,
				params.StatusCode,
				params.Latency,
				params.Request.UserAgent(),
				params.ErrorMessage,
			)
		},
	}
	r := gin.New()
	r.Use(gin.LoggerWithConfig(c))
	r.Use(gin.Recovery())
	r.Use(cors())

	api := r.Group("/v1")

	// 公共api，无需登陆即可访问
	public := api.Group("")
	{
		//CTFgo初始化
		public.POST("/install", u.Install)

		//用户注册
		public.POST("/register", u.Register)
		//用户登录
		public.POST("/login", u.Login)
		//用户注销
		public.GET("/logout", u.Logout)

		//获取验证码图片base64
		public.GET("/captcha", u.Captcha)
		//输出验证码图片
		public.GET("/captcha/:img", u.CaptchaServer)

		//获取指定id用户可公开信息
		public.GET("/info/:id", u.GetInfoByUserID)

		//获取指定id用户分数
		public.GET("/score/:id", u.GetScoreByUserID)
		//获取所有用户分数，降序排列。
		public.GET("/scores/all", u.GetAllScores)

		// 获取所有公告
		public.GET("/notice/all", admin.GetAllNotices)
	}

	// 普通用户api，需要用户登陆且Role=0才能访问
	personal := api.Group("/user")
	personal.Use(u.AuthRequired())
	{
		// 获取当前用户信息
		personal.GET("/session", u.Session)
		// 修改用户信息
		personal.POST("/updateinfo", u.UpdateInfo)

		// 获取所有题目信息
		personal.GET("/challenges/all", admin.GetAllChallenges)
		// 获取指定类别的题目信息
		personal.GET("/challenges", admin.GetChallengesByCategory)

		// 提交flag
		personal.POST("/flag", u.SubmitFlag)
		// 获取所有正确的flag提交记录
		personal.GET("solves/all", u.GetAllSolves)
		// 获取指定用户正确的flag提交记录
		personal.GET("solves/uid", u.GetSolvesByUid)
		// 获取指定题目正确的flag提交记录
		personal.GET("solves/cid", u.GetSolvesByCid)
	}

	// 管理者api，需要用户登陆且Role=1才能访问
	manager := api.Group("/admin")
	manager.Use(admin.AuthRequired())
	{
		// 创建新题目
		manager.POST("challenge", admin.NewChallenge)
		// 更改题目
		manager.PATCH("challenge", admin.EditChallenge)
		// 删除题目
		manager.DELETE("challenge", admin.DeleteChallenge)

		// 创建新公告
		manager.POST("notice", admin.NewNotice)
		// 更改公告
		manager.PATCH("notice", admin.EditNotice)
		// 删除公告
		manager.DELETE("notice", admin.DeleteNotice)

		// 获取所有提交记录
		manager.GET("submissions/all", u.GetAllSubmissions)
		// 获取指定用户的提交记录
		manager.GET("submissions/uid", u.GetSubmissionsByUid)
		// 获取指定题目的提交记录
		manager.GET("submissions/cid", u.GetSubmissionsByCid)
	}

	return r
}

//暂时跨域。
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Category, AccessToken, X-CSRF-Token, Authorization, Token, Content-Type")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Category")
		c.Header("Access-Control-Allow-Credentials", "true")
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
		}
		// 处理请求
		c.Next()
	}
}
