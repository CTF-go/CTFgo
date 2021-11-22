package types

// LoginRequest 定义接收登录数据的结构体。
type LoginRequest struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Remember bool   `form:"remember" json:"remember"`
}

// RegisterRequest 定义接收注册数据的结构体。
type RegisterRequest struct {
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	Email     string `form:"email" json:"email" binding:"required"`
	CaptchaID string `form:"captchaid" json:"captchaid" binding:"required"`
	Solution  string `form:"solution" json:"solution" binding:"required"`
}

// InfoRequest 定义接收用户修改信息的结构体。
type InfoRequest struct {
	Username    string `form:"username" json:"username"`
	Password    string `form:"password" json:"password"`
	Email       string `form:"email" json:"email"`
	Affiliation string `form:"affiliation" json:"affiliation"`
	Country     string `form:"country" json:"country"`
	Website     string `form:"website" json:"website"`
}

// InstallRequest 定义接收installRequest数据的结构体。
type InstallRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
}

// SubmissionRequest 定义接收提交flag数据的结构体。
type SubmissionRequest struct {
	Cid  int    `form:"cid" json:"cid" binding:"required"`
	Flag string `form:"flag" json:"flag" binding:"required"`
}

// ScoreResponse 定义返回得分情况结构体。
type ScoreResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
}

// SolveResponse 定义返回解题情况结构体。
type SolveResponse struct {
	ID            int    `json:"id"`
	Uid           int    `json:"uid"`
	Cid           int    `json:"cid"`
	Username      string `json:"username"`
	ChallengeName string `json:"challenge_name"`
	SubmittedAt   int    `json:"submitted_at"`
}

// PublicInfoResponse 定义返回用户公开信息结构体。
type PublicInfoResponse struct {
	Username    string `json:"username"`
	Affiliation string `json:"affiliation"`
	Country     string `json:"country"`
	TeamID      int    `json:"team_id"`
}

// ChallengeResponse 定义获取题目的一个响应。
type ChallengeResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Score       int    `json:"score"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Tags        string `json:"tags"`
	Hints       string `json:"hints"`
	SolverCount int    `json:"solver_count"`
}
