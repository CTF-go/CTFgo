package apiAdmin

// newChallengeRequest 定义新增题目的一个请求
type newChallengeRequest struct {
}

// editChallengeRequest 定义修改题目的一个请求
type editChallengeRequest struct {
}

// deleteChallengeRequest 定义删除题目的一个请求
type deleteChallengeRequest struct {
}

//// Bulletin 定义一个公告
//type Bulletin struct {
//	ID        int       `json:"id"`
//	Content   string    `json:"content"`
//	CreatedAt time.Time `json:"created_at"`
//	UpdatedAt time.Time `json:"updated_at"`
//}

// newBulletinRequest 定义新增公告的一个请求
type newBulletinRequest struct {
	Title   string `form:"title" json:"title" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
}

// editBulletinRequest 定义修改公告的一个请求
type editBulletinRequest struct {
	ID      int    `form:"id" json:"id" binding:"required"`
	Title   string `form:"title" json:"title" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
}

// deleteBulletinRequest 定义删除公告的一个请求
type deleteBulletinRequest struct {
	ID int `form:"id" json:"id" binding:"required"`
}
