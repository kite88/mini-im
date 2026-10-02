package dao

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mini-im/db"
	"mini-im/model"
)

// ContactStateDao 联系人「我这一侧」状态数据访问
type ContactStateDao struct{}

// NewContactStateDao 创建 ContactStateDao
func NewContactStateDao() *ContactStateDao { return &ContactStateDao{} }

func (ContactStateDao) table() *gorm.DB { return db.DB.Model(&model.ContactState{}) }

// FindByUser 查询我这一侧的全部联系人状态
func (d ContactStateDao) FindByUser(uid int64) ([]model.ContactState, error) {
	list := make([]model.ContactState, 0)
	err := d.table().Where("user_id = ?", uid).Find(&list).Error
	return list, err
}

// FindPair 查询两人之间双向的状态（最多两行），用于判断「谁删了谁」
func (d ContactStateDao) FindPair(a, b int64) ([]model.ContactState, error) {
	list := make([]model.ContactState, 0, 2)
	err := d.table().
		Where("(user_id = ? AND peer_id = ?) OR (user_id = ? AND peer_id = ?)", a, b, b, a).
		Find(&list).Error
	return list, err
}

// ClearTime 我清空与 peerID 会话的时间点；没有记录时返回 0，表示不受限
func (d ContactStateDao) ClearTime(uid, peerID int64) (int64, error) {
	var state model.ContactState
	err := d.table().Where("user_id = ? AND peer_id = ?", uid, peerID).First(&state).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return state.ClearTime, nil
}

// UpsertRemoved 标记「我把对方从好友中移除」，同时把会话清空时间推进到 now。
// 重复删除只刷新清空时间，不会产生多行。
func (d ContactStateDao) UpsertRemoved(uid, peerID, now int64) error {
	state := model.ContactState{UserID: uid, PeerID: peerID, Removed: true, ClearTime: now}
	return d.table().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "peer_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"removed": true, "clear_time": now}),
	}).Create(&state).Error
}

// Restore 双方重新成为好友：只清掉「已移除」标记，会话清空时间保留，
// 因此我这边重新加回来仍然是空会话，对方那边原本就有的记录也不会被翻出来。
func (d ContactStateDao) Restore(a, b int64) error {
	return d.table().
		Where("(user_id = ? AND peer_id = ?) OR (user_id = ? AND peer_id = ?)", a, b, b, a).
		Update("removed", false).Error
}
