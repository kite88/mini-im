package dao

import (
	"gorm.io/gorm"

	"mini-im/db"
	"mini-im/model"
)

// FriendDao 好友关系数据访问
type FriendDao struct{}

// NewFriendDao 创建 FriendDao
func NewFriendDao() *FriendDao { return &FriendDao{} }

func (FriendDao) table() *gorm.DB { return db.DB.Model(&model.Friend{}) }

// FindByUser 查询与指定用户相关的全部好友关系（双向）
func (d FriendDao) FindByUser(uid int64) ([]model.Friend, error) {
	list := make([]model.Friend, 0)
	err := d.table().Where("user_id = ? OR friend_id = ?", uid, uid).Find(&list).Error
	return list, err
}

// ExistsBetween 判断两个用户是否已是好友（双向匹配）
func (d FriendDao) ExistsBetween(uid, fid int64) (bool, error) {
	var count int64
	err := d.table().
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", uid, fid, fid, uid).
		Count(&count).Error
	return count > 0, err
}

// Create 新增一条好友关系（一条记录即表达双向关系）
func (d FriendDao) Create(friend *model.Friend) error {
	return d.table().Create(friend).Error
}
