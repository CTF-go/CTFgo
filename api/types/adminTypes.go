package types

// ChallengeRequest 定义新增/修改题目的一个请求
type ChallengeRequest struct {
	Name        string `form:"name" json:"name" binding:"required"`
	Score       int    `form:"score" json:"score" binding:"required"`
	Flag        string `form:"flag" json:"flag"`
	Description string `form:"description" json:"description"`
	Category    string `form:"category" json:"category" binding:"required"`
	Tags        string `form:"tags" json:"tags"`
	Hints       string `form:"hints" json:"hints"`
	Visible     int    `form:"visible" json:"visible"`
}

// NoticeRequest 定义新增公告的一个请求
type NoticeRequest struct {
	Title   string `form:"title" json:"title" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
}
