package model

// ReadState 会话已读游标
//
// 每个 (user_id, peer_id) 组合一行，read_time 表示「该用户已读完与 peer 的会话中
// create_time <= read_time 的消息」，因此未读数 = 对方发来的、晚于该时间的消息条数。
//
// 之所以记录游标而不是给每条消息打已读标记：消息是异步落库的（Redis 队列 + 后台任务），
// 用游标可以避免「标记已读发生在消息落库之前」从而产生的假未读。
type ReadState struct {
	ID         int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64 `gorm:"not null;uniqueIndex:uk_im_read_state_user_peer,priority:1" json:"user_id"`
	PeerID     int64 `gorm:"not null;uniqueIndex:uk_im_read_state_user_peer,priority:2" json:"peer_id"`
	ReadTime   int64 `gorm:"not null;default:0" json:"read_time"`
	UpdateTime int64 `gorm:"not null;default:0" json:"update_time"`
}

// TableName 指定表名
func (ReadState) TableName() string { return "im_read_state" }
