package model

type FriendModel struct {
	ID       int64 `json:"id"`
	UserId   int64 `json:"user_id"`
	FriendId int64 `json:"friend_id"`
	AddTime  int64 `json:"add_time"`
}

func (FriendModel) TableName() string {
	return "im_friend"
}