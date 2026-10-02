package dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mini-im/db"
	"mini-im/model"
)

// FriendRequestDao 好友申请数据访问
type FriendRequestDao struct{}

// NewFriendRequestDao 创建 FriendRequestDao
func NewFriendRequestDao() *FriendRequestDao { return &FriendRequestDao{} }

func (FriendRequestDao) table() *gorm.DB { return db.DB.Model(&model.FriendRequest{}) }

// Upsert 写入一条好友申请：
// 首次插入使用雪花 ID；已存在同方向记录时只刷新验证消息与时间，并重新置为待处理。
func (d FriendRequestDao) Upsert(req *model.FriendRequest) error {
	return d.table().Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "from_id"}, {Name: "to_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"message":     req.Message,
			"status":      model.FriendRequestPending,
			"update_time": req.UpdateTime,
		}),
	}).Create(req).Error
}

// FindByID 按主键查询，不存在返回 ErrNotFound
func (d FriendRequestDao) FindByID(id int64) (*model.FriendRequest, error) {
	var req model.FriendRequest
	if err := d.table().Where("id = ?", id).First(&req).Error; err != nil {
		return nil, notFound(err)
	}
	return &req, nil
}

// FindPending 查询 fromID -> toID 方向待处理的申请，不存在返回 ErrNotFound
func (d FriendRequestDao) FindPending(fromID, toID int64) (*model.FriendRequest, error) {
	var req model.FriendRequest
	err := d.table().
		Where("from_id = ? AND to_id = ? AND status = ?", fromID, toID, model.FriendRequestPending).
		First(&req).Error
	if err != nil {
		return nil, notFound(err)
	}
	return &req, nil
}

// FindPendingByTo 查询「待我处理」的申请，按申请时间倒序
func (d FriendRequestDao) FindPendingByTo(toID int64, limit int) ([]model.FriendRequest, error) {
	list := make([]model.FriendRequest, 0)
	err := d.table().
		Where("to_id = ? AND status = ?", toID, model.FriendRequestPending).
		Order("create_time DESC, id DESC").
		Limit(limit).
		Find(&list).Error
	return list, err
}

// FindPendingWith 查询 uid 与一批用户之间未处理完的申请（双向），
// 用于搜索结果里标注「我已申请 / 待我处理」。
func (d FriendRequestDao) FindPendingWith(uid int64, ids []int64) ([]model.FriendRequest, error) {
	list := make([]model.FriendRequest, 0)
	if len(ids) == 0 {
		return list, nil
	}
	err := d.table().
		Where("status = ? AND ((from_id = ? AND to_id IN ?) OR (to_id = ? AND from_id IN ?))",
			model.FriendRequestPending, uid, ids, uid, ids).
		Find(&list).Error
	return list, err
}

// DeletePendingBetween 删除两人之间「待处理」的申请（双向）。
// 拉黑时调用：残留的申请若被同意，会让刚拉黑的人又变成好友。
func (d FriendRequestDao) DeletePendingBetween(a, b int64) error {
	return d.table().
		Where("status = ? AND ((from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?))",
			model.FriendRequestPending, a, b, b, a).
		Delete(&model.FriendRequest{}).Error
}

// UpdatePending 把一条「待处理」的申请改为目标状态，返回是否真正命中。
// 带上 status 条件是为了并发处理同一申请时只有一个请求能生效。
func (d FriendRequestDao) UpdatePending(id int64, status int8, now int64) (bool, error) {
	res := d.table().
		Where("id = ? AND status = ?", id, model.FriendRequestPending).
		Updates(map[string]interface{}{
			"status":      status,
			"update_time": now,
		})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
