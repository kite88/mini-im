package model

// Blacklist 用户黑名单
//
// 与好友关系（一条记录表达双向）不同，黑名单是单向的：
// UserID 拉黑了 BlockedID。业务效果只作用在一个方向 ——
// 被拉黑方发给拉黑方的消息会被服务端拦截（见 service.CanSendMessage）。
//
// 主键用数据库自增：黑名单只是关联记录，ID 不对外暴露。
type Blacklist struct {
	ID        int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64 `gorm:"not null;uniqueIndex:uk_im_blacklist_pair,priority:1" json:"user_id"`
	BlockedID int64 `gorm:"not null;uniqueIndex:uk_im_blacklist_pair,priority:2;index:idx_im_blacklist_blocked" json:"blocked_id"`
	AddTime   int64 `gorm:"not null;default:0" json:"add_time"`
}

// TableName 指定表名
func (Blacklist) TableName() string { return "im_blacklist" }
