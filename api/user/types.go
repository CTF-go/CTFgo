/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

//user定义用户结构体。
type User struct {
	ID          int    `json:"id"`          //用户id，唯一，自增
	Token       string `json:"token"`       //用户token，唯一
	Username    string `json:"username"`    //用户名，唯一
	Password    string `json:"password"`    //用户密码md5值，md5(原密码）
	Email       string `json:"email"`       //邮箱，唯一
	Affiliation string `json:"affiliation"` //组织、战队或机构等，非必需，默认为0
	Country     string `json:"country"`     //国家，非必需，默认为0
	Website     string `json:"website"`     //个人链接，默认为0
	Hidden      int    `json:"hidden"`      //1：隐藏，0：显示，默认为0
	Banned      int    `json:"banned"`      //1：禁止，0：正常，默认为1，邮箱激活后为0
	TeamID      int    `json:"team_id"`     //队伍id，在团队模式下必须，个人模式非必需，默认为0
	Created     int    `json:"created"`     //用户注册时间，10位数时间戳
	Role        int    `json:"role"`        //0：普通用户，默认为0，1：普通管理员，2：所有者（最高权限）
}

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

type getSubmissionsByUidRequest struct {
	Uid int `form:"uid" json:"uid" binding:"required"`
}

type getSubmissionsByCidRequest struct {
	Cid int `form:"cid" json:"cid" binding:"required"`
}

type getSolvesByUidRequest struct {
	Uid int `form:"uid" json:"uid" binding:"required"`
}

type getSolvesByCidRequest struct {
	Cid int `form:"cid" json:"cid" binding:"required"`
}

// scoreResponse 定义返回得分情况结构体。
type scoreResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
}

// publicInfoResponse 定义返回用户公开信息结构体。
type publicInfoResponse struct {
	Username    string `json:"username"`
	Affiliation string `json:"affiliation"`
	Country     string `json:"country"`
	TeamID      int    `json:"team_id"`
}
