package dao

import (
	"gorm.io/gorm"

	"mini-im/db"
	"mini-im/model"
)

// BlacklistDao 黑名单数据访问
type BlacklistDao struct{}

// NewBlacklistDao 创建 BlacklistDao
func NewBlacklistDao() *BlacklistDao { return &BlacklistDao{} }

func (BlacklistDao) table() *gorm.DB { return db.DB.Model(&model.Blacklist{}) }

// FindByUser 查询「我拉黑的人」，按拉黑时间倒序
func (d BlacklistDao) FindByUser(uid int64, limit int) ([]model.Blacklist, error) {
	list := make([]model.Blacklist, 0)
	err := d.table().
		Where("user_id = ?", uid).
		Order("add_time DESC, id DESC").
		Limit(limit).
		Find(&list).Error
	return list, err
}

// FindBlockedIDs 我拉黑的全部用户 ID（不分页，用于好友 / 会话列表过滤）
func (d BlacklistDao) FindBlockedIDs(uid int64) ([]int64, error) {
	ids := make([]int64, 0)
	err := d.table().
		Where("user_id = ?", uid).
		Pluck("blocked_id", &ids).Error
	return ids, err
}

// Exists 判断 userID 是否已拉黑 blockedID（单向）
func (d BlacklistDao) Exists(userID, blockedID int64) (bool, error) {
	var count int64
	err := d.table().
		Where("user_id = ? AND blocked_id = ?", userID, blockedID).
		Count(&count).Error
	return count > 0, err
}

// Create 新增一条黑名单记录（调用方需保证同一对用户不重复）
func (d BlacklistDao) Create(item *model.Blacklist) error {
	return d.table().Create(item).Error
}

// Delete 解除拉黑，返回是否真的删除了记录（未拉黑过时为 false）
func (d BlacklistDao) Delete(userID, blockedID int64) (bool, error) {
	res := d.table().
		Where("user_id = ? AND blocked_id = ?", userID, blockedID).
		Delete(&model.Blacklist{})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
