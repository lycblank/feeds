package users

type UserLogin struct {
	Uid int64            `json:"uid"`
	LoginName string     `json:"login_name"`
	LoginPassword string `json:"login_password"`
	Salt string          `json:"salt"`
	LastLoginTime int64  `json:"last_login_time"`
	CreateTime int64     `json:"create_time"`
}






