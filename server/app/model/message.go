package model

type MessageModel struct {
	ID          int64     `json:"id"`
	ContentType int64     `json:"content_type"`
	Content     string    `json:"content"`
	CreateTime  int64     `json:"create_time"`
	UserId      int64     `json:"user_id"`
	FriendId    int64     `json:"friend_id"`
	Type        int8      `json:"type" gorm:"-"`
	User        UserModel `json:"user" gorm:"-"`
}

func (MessageModel) TableName() string {
	return "im_message"
}
