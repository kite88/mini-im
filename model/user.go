// Package model 定义数据库实体。
package model

import (
	"gorm.io/gorm"

	"mini-im/pkg/snowflake"
)

// User 用户实体（客服 / 客户）
//
// 时间字段统一使用 Unix 秒（int64），与前端约定保持一致。
//
// 主键使用应用侧生成的雪花 ID，而非数据库自增（autoIncrement:false），
// 这样 ID 在插入前即可确定；同时因雪花 ID 是 18~19 位整数，超出
// JavaScript 的安全整数范围（2^53），JSON 统一按字符串输出避免精度丢失。
type User struct {
	ID         int64  `gorm:"primaryKey;autoIncrement:false" json:"id,string"`
	Username   string `gorm:"size:64;not null;default:'';uniqueIndex:uk_im_user_username" json:"username"`
	Password   string `gorm:"size:128;not null;default:''" json:"-"`
	Nickname   string `gorm:"size:64;not null;default:''" json:"nickname"`
	Avatar     string `gorm:"size:255;not null;default:''" json:"avatar"`
	CreateTime int64  `gorm:"not null;default:0" json:"-"`
	UpdateTime int64  `gorm:"not null;default:0" json:"-"`
}

// TableName 指定表名
func (User) TableName() string { return "im_user" }

// BeforeCreate 入库前填充雪花 ID，保证任何写入口都不会漏掉主键
func (u *User) BeforeCreate(*gorm.DB) error {
	if u.ID <= 0 {
		u.ID = snowflake.NextID()
	}
	return nil
}
