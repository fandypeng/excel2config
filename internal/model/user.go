package model

type Token struct {
	Token      string `json:"token" bson:"token"`
	UpdateTime int64  `json:"update_time" bson:"update_time"`
}

type UserInfo struct {
	Uid       string `json:"uid,omitempty" bson:"_id,omitempty"`
	UserName  string `json:"userName,omitempty" bson:"username"`
	Email     string `json:"email,omitempty" bson:"email"`
	Passwd    string `json:"-" bson:"passwd"`
	RegTime   int64  `json:"regTime,omitempty" bson:"reg_time"`
	Avatar    string `json:"avatar,omitempty" bson:"avatar"`
	LoginType int    `json:"login_type" bson:"login_type"`
	OpenId    string `json:"open_id" bson:"open_id"`
}

type SimpleUserInfo struct {
	Uid      string `json:"uid,omitempty" bson:"uid,omitempty"`
	UserName string `json:"userName,omitempty" bson:"username,omitempty"`
	Avatar   string `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Role     int    `json:"role,omitempty" bson:"role,omitempty"`
}
