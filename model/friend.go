package model

// Friend 好友关系
//
// 一条记录即表达双向关系，查询时按 (user_id / friend_id) 双向匹配。
// 写入来源有两个：seed 的演示数据（两两互加）与「加好友」被同意后的建立关系。
//
// 记录只在建立好友时写入，删除好友**不删除**这行：
// 谁删了谁、谁的会话被清空这类单侧状态记在 im_contact_state，
// 因此双方重新加为好友时不需要重新插行，历史也能按人区分。
type Friend struct {
	ID       int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID   int64 `gorm:"not null;index:idx_im_friend_user" json:"user_id"`
	FriendID int64 `gorm:"not null;index:idx_im_friend_friend" json:"friend_id"`
	AddTime  int64 `gorm:"not null;default:0" json:"add_time"`
}

// TableName 指定表名
func (Friend) TableName() string { return "im_friend" }
