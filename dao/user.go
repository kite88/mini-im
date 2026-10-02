// Package dao 封装数据库访问，业务层不直接拼 SQL。
package dao

import (
	"strings"

	"gorm.io/gorm"

	"mini-im/db"
	"mini-im/model"
)

// UserDao 用户数据访问
type UserDao struct{}

// NewUserDao 创建 UserDao
func NewUserDao() *UserDao { return &UserDao{} }

func (UserDao) table() *gorm.DB { return db.DB.Model(&model.User{}) }

// FindByID 按主键查询
func (d UserDao) FindByID(id int64) (*model.User, error) {
	var user model.User
	if err := d.table().Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 按用户名查询
func (d UserDao) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := d.table().Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByIDs 批量按主键查询
func (d UserDao) FindByIDs(ids []int64) ([]model.User, error) {
	users := make([]model.User, 0, len(ids))
	if len(ids) == 0 {
		return users, nil
	}
	err := d.table().Where("id IN ?", ids).Find(&users).Error
	return users, err
}

// SearchByKeyword 按账号 / 昵称模糊搜索用户（不区分大小写），排除指定用户
func (d UserDao) SearchByKeyword(keyword string, excludeID int64, limit int) ([]model.User, error) {
	users := make([]model.User, 0, limit)
	like := "%" + escapeLike(keyword) + "%"
	err := d.table().
		Where("id <> ?", excludeID).
		Where("username ILIKE ? OR nickname ILIKE ?", like, like).
		Order("username ASC").
		Limit(limit).
		Find(&users).Error
	return users, err
}

// Create 新增用户
func (d UserDao) Create(user *model.User) error {
	return d.table().Create(user).Error
}

// escapeLike 转义 LIKE 通配符（PostgreSQL 默认转义符为反斜杠），
// 避免用户输入的 % / _ 被当成通配符而搜出全表。
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "%", `\%`)
	return strings.ReplaceAll(s, "_", `\_`)
}
