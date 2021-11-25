package types

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

// Notice 定义一个公告
type Notice struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int    `json:"created_at"`
}

// Challenge 定义一个题目
type Challenge struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Score       int      `json:"score"`
	Flag        string   `json:"flag"`
	Description string   `json:"description"`
	Attachment  []string `json:"attachment"`
	Category    string   `json:"category"`
	Tags        string   `json:"tags"`
	Hints       []string `json:"hints"`
	Visible     int      `json:"visible"` // 0表示隐藏，1表示可见
}

// Submission 表示一次flag提交记录
type Submission struct {
	ID          int    `json:"id"`
	UserID      int    `json:"uid"`
	ChallengeID int    `json:"cid"`
	Flag        string `json:"flag"`
	IP          string `json:"ip"`
	Time        int    `json:"submitted_at"`
}

// Solve 表示一次正确的flag提交记录
type Solve struct {
	ID          int `json:"id"`
	UserID      int `json:"uid"`
	ChallengeID int `json:"cid"`
	Time        int `json:"solved_at"`
}
