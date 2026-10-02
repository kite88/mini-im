package dao

import (
	"gorm.io/gorm"

	"mini-im/db"
	"mini-im/model"
)

// MessageDao 消息数据访问
type MessageDao struct{}

// NewMessageDao 创建 MessageDao
func NewMessageDao() *MessageDao { return &MessageDao{} }

func (MessageDao) table() *gorm.DB { return db.DB.Model(&model.Message{}) }

// Create 持久化一条消息
func (d MessageDao) Create(msg *model.Message) error {
	return d.table().Create(msg).Error
}

// FindBetween 查询两个用户之间的消息，按时间倒序。
// since > 0 时只返回该时间点之后的消息（「我清空过会话」时的过滤条件）。
func (d MessageDao) FindBetween(uid, fid int64, limit int, since int64) ([]model.Message, error) {
	query := d.table().
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", uid, fid, fid, uid)
	if since > 0 {
		query = query.Where("create_time > ?", since)
	}

	list := make([]model.Message, 0)
	err := query.
		Order("create_time DESC, id DESC").
		Limit(limit).
		Find(&list).Error
	return list, err
}

// FindRecentByUser 查询与指定用户相关的最近消息，按时间倒序
func (d MessageDao) FindRecentByUser(uid int64, limit int) ([]model.Message, error) {
	list := make([]model.Message, 0)
	err := d.table().
		Where("user_id = ? OR friend_id = ?", uid, uid).
		Order("create_time DESC, id DESC").
		Limit(limit).
		Find(&list).Error
	return list, err
}
