package types

// LoginRequest 定义接收登录数据的结构体。
type LoginRequest struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Remember bool   `json:"remember"`
}

// RegisterRequest 定义接收注册数据的结构体。
type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	CaptchaID string `json:"captchaid" binding:"required"`
	Solution  string `json:"solution" binding:"required"`
}

// InfoRequest 定义接收用户修改信息的结构体。
type InfoRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Affiliation string `json:"affiliation"`
	Country     string `json:"country"`
	Website     string `json:"website"`
}

// InstallRequest 定义接收installRequest数据的结构体。
type InstallRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

// SubmissionRequest 定义接收提交flag数据的结构体。
type SubmissionRequest struct {
	Cid  int    `json:"cid" binding:"required"`
	Flag string `json:"flag" binding:"required"`
}

// StudentInfo 定义校内学生信息结构体。
type StudentInfo struct {
	Username  string `json:"username"`
	StudentID string `json:"student_id"`
	QQ        string `json:"qq"`
}

// SubmitStudentInfoRequest 定义接收提交校内学生用户信息的结构体。
type SubmitStudentInfoRequest struct {
	Student1 StudentInfo `json:"student1" binding:"required"`
	Student2 StudentInfo `json:"student2"`
	Student3 StudentInfo `json:"student3"`
	Student4 StudentInfo `json:"student4"`
}

// OthersInfo 定义校外用户信息结构体。
type OthersInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	QQ       string `json:"qq"`
}

// SubmitOthersInfoRequest 定义接收提交校外学生用户信息的结构体。
type SubmitOthersInfoRequest struct {
	Others1 OthersInfo `json:"others1" binding:"required"`
	Others2 OthersInfo `json:"others2"`
	Others3 OthersInfo `json:"others3"`
	Others4 OthersInfo `json:"others4"`
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
	Score         int    `json:"score"`
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
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Score       int      `json:"score"`
	Description string   `json:"description"`
	Attachment  []string `json:"attachment"`
	Category    string   `json:"category"`
	Tags        string   `json:"tags"`
	Hints       []string `json:"hints"`
	SolverCount int      `json:"solver_count"`
	IsSolved    bool     `json:"is_solved"` // true：已解决，false：未解决
}

// ScoreRankResponse 定义获取当前用户分数和排名的一个响应。
type ScoreRankResponse struct {
	Score int `json:"score"`
	Rank  int `json:"rank"`
}

// StudentsOrOthersInfoResponse 定义获取校内用户或校外用户信息的响应。
type StudentsOrOthersInfoResponse struct {
	Username  string `json:"username"`
	IDOrEmail string `json:"id_email"`
	QQ        string `json:"qq"`
}
