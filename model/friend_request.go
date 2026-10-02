package model

import (
	"gorm.io/gorm"

	"mini-im/pkg/snowflake"
)

// 好友申请状态
const (
	FriendRequestPending  int8 = 1 // 待处理
	FriendRequestAccepted int8 = 2 // 已同意（已成为好友）
	FriendRequestRejected int8 = 3 // 已拒绝
)

// FriendRequest 好友申请
//
// 同一对 (from_id, to_id) 只有一行，重复申请走 upsert 更新状态与验证消息，
// 因此「申请 → 被拒 → 再申请」不会堆积垃圾数据。
// 与用户 / 消息表一致：主键为雪花 ID，对外 ID 按字符串序列化。
type FriendRequest struct {
	ID      int64  `gorm:"primaryKey;autoIncrement:false" json:"id,string"`
	FromID  int64  `gorm:"not null;uniqueIndex:uk_im_friend_request_from_to,priority:1" json:"from_id,string"`
	ToID    int64  `gorm:"not null;uniqueIndex:uk_im_friend_request_from_to,priority:2;index:idx_im_friend_request_to_status,priority:1" json:"to_id,string"`
	Message string `gorm:"size:255;not null;default:''" json:"message"`
	Status  int8   `gorm:"not null;default:1;index:idx_im_friend_request_to_status,priority:2" json:"status"`
	// CreateTime 首次申请时间；重复申请只刷新 UpdateTime，便于列表稳定排序
	CreateTime int64 `gorm:"not null;default:0" json:"create_time"`
	UpdateTime int64 `gorm:"not null;default:0" json:"update_time"`
}

// TableName 指定表名
func (FriendRequest) TableName() string { return "im_friend_request" }

// BeforeCreate 入库前填充雪花 ID
func (r *FriendRequest) BeforeCreate(*gorm.DB) error {
	if r.ID <= 0 {
		r.ID = snowflake.NextID()
	}
	return nil
}
