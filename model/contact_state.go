package model

// ContactState 联系人的「我这一侧」状态（单向）
//
// 好友关系（im_friend）是一条记录表达双向，而删除好友、清空会话都是**单向**行为：
// 我删了对方，只影响我的好友列表 / 会话列表 / 我的聊天记录，对方的列表与记录不动。
// 这类只对某一侧生效的状态记在这里，每个 (user_id, peer_id) 一行。
//
// 与好友表的关系：删除好友不改 im_friend（记录保留），只把 Removed 置为 true；
// 双方重新加为好友时把 Removed 清回去，ClearTime 保留，
// 因此「重新加回来是空的」，而对方那侧的历史记录不受影响。
type ContactState struct {
	ID     int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID int64 `gorm:"not null;uniqueIndex:uk_im_contact_state_pair,priority:1" json:"user_id"`
	PeerID int64 `gorm:"not null;uniqueIndex:uk_im_contact_state_pair,priority:2;index:idx_im_contact_state_peer" json:"peer_id"`
	// Removed 我是否已把对方从好友中移除（单向，只影响我这一侧）
	Removed bool `gorm:"not null;default:false" json:"removed"`
	// ClearTime 我清空与对方会话的时间点，早于它的消息不再出现在我的会话列表 / 聊天记录里
	ClearTime int64 `gorm:"not null;default:0" json:"clear_time"`
}

// TableName 指定表名
func (ContactState) TableName() string { return "im_contact_state" }
