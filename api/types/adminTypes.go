package types

// ChallengeRequest 定义新增/修改题目的一个请求
type ChallengeRequest struct {
	Name        string   `json:"name" binding:"required"`
	Score       int      `json:"score" binding:"required"`
	Flag        string   `json:"flag"` // 暂时一个题只能一个flag
	Description string   `json:"description"`
	Attachment  []string `json:"attachment"`
	Category    string   `json:"category" binding:"required"`
	Tags        string   `json:"tags"`
	Hints       []string `json:"hints"`
	Visible     int      `json:"visible"`
}

// NoticeRequest 定义新增公告的一个请求
type NoticeRequest struct {
	Title   string `form:"title" json:"title" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
}
