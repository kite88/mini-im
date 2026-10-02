// Package dao 封装数据库访问，业务层不直接拼 SQL。
package dao

import (
	"errors"

	"gorm.io/gorm"
)

// ErrNotFound 查询记录不存在，由各 Dao 统一转换，
// 避免上层（service / api）为了判断「不存在」而依赖 gorm 的错误类型。
var ErrNotFound = errors.New("记录不存在")

// notFound 把 gorm 的「记录不存在」统一转换为 ErrNotFound，其余错误原样返回
func notFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}
