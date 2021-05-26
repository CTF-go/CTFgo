package apiAdmin

// newChallengeRequest 定义新增题目的一个请求
type newChallengeRequest struct {
	Name        string `form:"name" json:"name" binding:"required"`
	Score       int    `form:"score" json:"score" binding:"required"`
	Flag        string `form:"flag" json:"flag"`
	Description string `form:"description" json:"description"`
	Category    string `form:"category" json:"category" binding:"required"`
	Tags        string `form:"tags" json:"tags"`
	Hints       string `form:"hints" json:"hints"`
	Visible     int    `form:"visible" json:"visible"`
}

// editChallengeRequest 定义修改题目的一个请求
type editChallengeRequest struct {
	ID          int    `form:"id" json:"id" binding:"required"`
	Name        string `form:"name" json:"name" binding:"required"`
	Score       int    `form:"score" json:"score" binding:"required"`
	Flag        string `form:"flag" json:"flag"`
	Description string `form:"description" json:"description"`
	Category    string `form:"category" json:"category" binding:"required"`
	Tags        string `form:"tags" json:"tags"`
	Hints       string `form:"hints" json:"hints"`
	Visible     int    `form:"visible" json:"visible"`
}

// deleteChallengeRequest 定义删除题目的一个请求
type deleteChallengeRequest struct {
	ID int `form:"id" json:"id" binding:"required"`
}

// getChallengeByCategoryRequest 定义获取指定类别题目的一个请求
type getChallengeByCategoryRequest struct {
	Category string `form:"category" json:"category" binding:"required"`
}

// newNoticeRequest 定义新增公告的一个请求
type newNoticeRequest struct {
	Title   string `form:"title" json:"title" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
}

// editNoticeRequest 定义修改公告的一个请求
type editNoticeRequest struct {
	ID      int    `form:"id" json:"id" binding:"required"`
	Title   string `form:"title" json:"title" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
}

// deleteNoticeRequest 定义删除公告的一个请求
type deleteNoticeRequest struct {
	ID int `form:"id" json:"id" binding:"required"`
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
}
