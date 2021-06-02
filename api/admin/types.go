package apiAdmin

// challengeRequest 定义新增/修改题目的一个请求
type challengeRequest struct {
	Name        string `form:"name" json:"name" binding:"required"`
	Score       int    `form:"score" json:"score" binding:"required"`
	Flag        string `form:"flag" json:"flag"`
	Description string `form:"description" json:"description"`
	Category    string `form:"category" json:"category" binding:"required"`
	Tags        string `form:"tags" json:"tags"`
	Hints       string `form:"hints" json:"hints"`
	Visible     int    `form:"visible" json:"visible"`
}

// challengeResponse 定义获取题目的一个响应
type challengeResponse struct {
	ID          int    `form:"id" json:"id" binding:"required"`
	Name        string `form:"name" json:"name" binding:"required"`
	Score       int    `form:"score" json:"score" binding:"required"`
	Description string `form:"description" json:"description"`
	Category    string `form:"category" json:"category" binding:"required"`
	Tags        string `form:"tags" json:"tags"`
	Hints       string `form:"hints" json:"hints"`
	SolverCount int    `form:"solver_count" json:"solver_count"`
}

// noticeRequest 定义新增公告的一个请求
type noticeRequest struct {
	Title   string `form:"title" json:"title" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
}
