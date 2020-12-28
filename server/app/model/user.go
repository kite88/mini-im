package model

type UserModel struct {
	ID         int64  `json:"id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	CreateTime int64  `json:"-"`
	UpdateTime int64  `json:"-"`
	Username   string `json:"username"`
	Password   string `json:"-"`
}

func (UserModel) TableName() string {
	return "im_user"
}
