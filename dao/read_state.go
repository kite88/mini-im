package dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mini-im/db"
	"mini-im/model"
)

// ReadStateDao 会话已读游标数据访问
type ReadStateDao struct{}

// NewReadStateDao 创建 ReadStateDao
func NewReadStateDao() *ReadStateDao { return &ReadStateDao{} }

func (ReadStateDao) table() *gorm.DB { return db.DB.Model(&model.ReadState{}) }

// Upsert 推进「userID 已读到与 peerID 的会话」的游标，只前进不后退
func (d ReadStateDao) Upsert(userID, peerID, readTime int64) error {
	state := model.ReadState{
		UserID:     userID,
		PeerID:     peerID,
		ReadTime:   readTime,
		UpdateTime: readTime,
	}
	return d.table().Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "peer_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			// 并发上报时取较大值，避免游标被回拨
			"read_time":   gorm.Expr("GREATEST(im_read_state.read_time, ?)", readTime),
			"update_time": readTime,
		}),
	}).Create(&state).Error
}

// UnreadCounts 统计每个会话的未读消息数，key 为对方用户 ID
//
// 未读判定：对方发给我、且时间晚于我的已读游标。
// 游标缺失时（未读功能上线前就存在的会话）回退到「我在该会话中最后一次发言时间」，
// 避免把早已看过的历史消息全部算成未读。
func (d ReadStateDao) UnreadCounts(uid int64) (map[int64]int64, error) {
	type unreadRow struct {
		PeerID int64 `gorm:"column:peer_id"`
		Cnt    int64 `gorm:"column:cnt"`
	}

	rows := make([]unreadRow, 0)
	err := db.DB.Raw(`
		SELECT m.user_id AS peer_id, COUNT(*) AS cnt
		FROM im_message m
		LEFT JOIN im_read_state r ON r.user_id = ? AND r.peer_id = m.user_id
		WHERE m.friend_id = ?
		  AND m.create_time > COALESCE(r.read_time,
		      (SELECT COALESCE(MAX(s.create_time), 0)
		         FROM im_message s
		        WHERE s.user_id = ? AND s.friend_id = m.user_id))
		GROUP BY m.user_id`, uid, uid, uid).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[int64]int64, len(rows))
	for _, row := range rows {
		counts[row.PeerID] = row.Cnt
	}
	return counts, nil
}
