package model

import (
	"strings"

	"gorm.io/gorm"

	"mini-im/pkg/uuidv7"
)

// 消息内容类型（字符串标识）
//
// 用可读标识代替 1 / 2 这类魔法数字：落库、HTTP 出参、WebSocket 报文三处
// 用的是同一套值，前端直接比较字符串，不必再维护「数字 -> 含义」的映射。
const (
	ContentTypeText  = "text"  // 文本
	ContentTypeImage = "image" // 图片 / 表情（内容为 HTML 片段）
)

// NormalizeContentType 归一化客户端传来的内容类型：
// 只接受已知标识，未知 / 空值一律回落为文本，避免脏数据落库。
func NormalizeContentType(value string) string {
	switch strings.TrimSpace(value) {
	case ContentTypeImage:
		return ContentTypeImage
	default:
		return ContentTypeText
	}
}

// 会话方向（仅接口返回时使用，不落库）
const (
	MsgTypeSelf int8 = 1 // 自己发出
	MsgTypePeer int8 = 2 // 对方发来
)

// Message 聊天消息实体
//
// 主键为 UUID v7（应用侧生成，字符串进出）：消息量最大、只增不改，用 UUID 可以
// 彻底摆脱节点号协调与 64 位容量焦虑，同时时间戳在最高位、同毫秒用序列号递增，
// 因此仍然按时间有序，B 树索引表现为顺序追加（与雪花 ID 的索引友好性一致）。
// 用户 / 好友申请仍保留雪花 ID，见 im_user / im_friend_request。
//
// 消息只增不改（清理会话是「我这一侧」的状态，见 im_contact_state），
// 因此消息行本身不带软删除标记。
type Message struct {
	ID          string `gorm:"type:uuid;primaryKey" json:"id"`
	ContentType string `gorm:"size:16;not null;default:'text'" json:"content_type"`
	Content     string `gorm:"type:text;not null;default:''" json:"content"`
	UserID      int64  `gorm:"not null;index:idx_im_message_user" json:"user_id,string"`
	FriendID    int64  `gorm:"not null;index:idx_im_message_friend" json:"friend_id,string"`
	CreateTime  int64  `gorm:"not null;index:idx_im_message_create_time" json:"create_time"`

	// 以下字段不落库，仅用于接口返回
	Type int8  `gorm:"-" json:"type"`
	User *User `gorm:"-" json:"user,omitempty"`
}

// TableName 指定表名
func (Message) TableName() string { return "im_message" }

// BeforeCreate 入库前填充 UUID v7，保证任何写入口都不会漏掉主键
func (m *Message) BeforeCreate(*gorm.DB) error {
	if m.ID == "" {
		m.ID = uuidv7.New()
	}
	return nil
}
