/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

//user定义用户结构体。
type User struct {
	ID          int    `json:"id"`          //用户id，唯一，自增
	Token       string `json:"token"`       //用户token，唯一，API鉴权使用
	Username    string `json:"username"`    //用户名，唯一
	Password    string `json:"password"`    //用户密码md5值，md5(原密码）
	Email       string `json:"email"`       //邮箱，唯一
	Affiliation string `json:"affiliation"` //组织、战队或机构等，非必需
	Country     string `json:"country"`     //国家，非必需
	Hidden      int    `json:"hidden"`      //1：隐藏，0：显示，默认为0
	Banned      int    `json:"banned"`      //1：禁止，0：正常，默认为1，邮箱激活后为0
	TeamID      int    `json:"team_id"`     //队伍id，在团队模式下必须，个人模式非必需
	Created     string `json:"created"`     //用户注册时间，10位数时间戳
	Role        int    `json:"role"`        //0：普通用户，默认为0，1：普通管理员，2：所有者（最高权限）
}

// loginRequest 定义接收登录数据的结构体。
type loginRequest struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// registerRequest 定义接收注册数据的结构体。
type registerRequest struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	Email     string `form:"email" json:"email" binding:"required"`
	CaptchaID string `form:"captchaid" json:"captchaid" binding:"required"`
	Solution  string `form:"solution" json:"solution" binding:"required"`
}

// infoRequest 定义接收用户信息的结构体。
type infoRequest struct {
	Username    string `form:"username" json:"username"`
	Password    string `form:"password" json:"password"`
	Email       string `form:"email" json:"email"`
	Affiliation string `form:"affiliation" json:"affiliation"`
	Country     string `form:"country" json:"country"`
}

// scores 定义返回得分情况结构体。
type Score struct {
	ID       int
	Username string
	Score    int
}

// installRequest 定义接收installRequest数据的结构体。
type installRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
}

// isExistedRequest 定义接收用户名或邮箱数据的结构体。
type isExistedRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
}

type submissionRequest struct {
	ChallengeID int    `form:"challenge_id" json:"challenge_id" binding:"required"`
	Flag        string `form:"flag" json:"flag" binding:"required"`
}
