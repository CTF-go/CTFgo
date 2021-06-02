/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

// loginRequest 定义接收登录数据的结构体。
type loginRequest struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Remember bool   `form:"remember" json:"remember"`
}

// registerRequest 定义接收注册数据的结构体。
type registerRequest struct {
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	Email     string `form:"email" json:"email" binding:"required"`
	CaptchaID string `form:"captchaid" json:"captchaid" binding:"required"`
	Solution  string `form:"solution" json:"solution" binding:"required"`
}

// infoRequest 定义接收用户修改信息的结构体。
type infoRequest struct {
	Username    string `form:"username" json:"username"`
	Password    string `form:"password" json:"password"`
	Email       string `form:"email" json:"email"`
	Affiliation string `form:"affiliation" json:"affiliation"`
	Country     string `form:"country" json:"country"`
	Website     string `form:"website" json:"website"`
}

// installRequest 定义接收installRequest数据的结构体。
type installRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
}

type submissionRequest struct {
	Cid  int    `form:"cid" json:"cid" binding:"required"`
	Flag string `form:"flag" json:"flag" binding:"required"`
}

// scoreResponse 定义返回得分情况结构体。
type scoreResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
}

type solveResponse struct {
	ID            int    `json:"id"`
	Uid           int    `json:"uid"`
	Cid           int    `json:"cid"`
	Username      string `json:"username"`
	ChallengeName string `json:"challenge_name"`
	SubmittedAt   int    `json:"submitted_at"`
}

// publicInfoResponse 定义返回用户公开信息结构体。
type publicInfoResponse struct {
	Username    string `json:"username"`
	Affiliation string `json:"affiliation"`
	Country     string `json:"country"`
	TeamID      int    `json:"team_id"`
}
