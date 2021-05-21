/*
Package apiUser实现用户和其他普通API接口功能。
*/
package apiUser

//user定义用户结构体。
type Users struct {
	ID          int    `json:"id"`          //用户id，唯一，自增
	Token       string `json:"token"`       //用户token，唯一，API鉴权使用
	Username    string `json:"username"`    //用户名，唯一
	Password    string `json:"password"`    //用户密码md5值，md5(原密码）
	Email       string `json:"email"`       //邮箱，唯一
	Affiliation string `json:"affiliation"` //组织、战队或机构等，非必需
	Country     string `json:"country"`     //国家，非必需
	Hidden      int    `json:"hidden"`      //1：隐藏，0：显示，默认为0
	Banned      int    `json:"banned"`      //1：禁止，0：正常，默认为1，邮箱激活后为0
	Team_id     int    `json:"team_id"`     //队伍id，在团队模式下必须，个人模式非必需
	Created     string `json:"created"`     //用户注册时间，10位数时间戳
	Role        int    `json:"role"`        //0：普通用户，默认为0，1：普通管理员，2：所有者（最高权限）
}

//login_struct定义接收Login数据的结构体。
type login_struct struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	User   string `form:"username" json:"username" binding:"required"`
	Passwd string `form:"password" json:"password" binding:"required"`
}

//register_struct定义接收Login数据的结构体。
type register_struct struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	User   string `form:"username" json:"username" binding:"required"`
	Passwd string `form:"password" json:"password" binding:"required"`
	Email  string `form:"email" json:"email" binding:"required"`
}

//info_struct定义接收用户信息的结构体。
type info_struct struct {
	User        string `form:"username" json:"username"`
	Passwd      string `form:"password" json:"password"`
	Email       string `form:"email" json:"email"`
	Affiliation string `form:"affiliation" json:"affiliation"`
	Country     string `form:"country" json:"country"`
}

//scores_struct 定义返回得分情况结构体。
type scores_struct struct {
	ID    int
	User  string
	Score int
}

//captcha_struct 定义接收用户输入验证码和验证码id的结构体。
type captcha_struct struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Captcha_id       string `form:"id" json:"id" binding:"required"`
	Captcha_solution string `form:"solution" json:"solution" binding:"required"`
}

//install_struct定义接收install数据的结构体。
type install_struct struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	User   string `form:"username" json:"username" binding:"required"`
	Passwd string `form:"password" json:"password" binding:"required"`
	Email  string `form:"email" json:"email" binding:"required"`
}
