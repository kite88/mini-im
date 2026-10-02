package service

import (
	"errors"
	"time"

	"mini-im/dao"
	"mini-im/model"
)

const (
	// blacklistLimit 单次拉取的黑名单上限
	blacklistLimit = 200
)

var (
	// ErrBlacklistSelf 不能拉黑自己
	ErrBlacklistSelf = errors.New("不能将自己拉入黑名单")
	// ErrBlacklistNotFound 目标用户不在黑名单中
	ErrBlacklistNotFound = errors.New("该用户不在黑名单中")
	// ErrBlacklistBlockedPeer 我已拉黑对方，需先移出黑名单
	ErrBlacklistBlockedPeer = errors.New("请先将对方移出黑名单")
	// ErrBlacklistRejected 对方拉黑了我，好友申请一律拒绝（不暴露对方拉黑的事实）
	ErrBlacklistRejected = errors.New("对方拒绝接收好友申请")
)

// BlacklistItem 黑名单项：被拉黑的用户资料 + 拉黑时间
type BlacklistItem struct {
	ID       int64  `json:"id,string"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	AddTime  int64  `json:"add_time"`
}

// Blacklist 我拉黑的用户列表（按拉黑时间倒序）
func Blacklist(uid int64) ([]BlacklistItem, error) {
	list, err := dao.NewBlacklistDao().FindByUser(uid, blacklistLimit)
	if err != nil {
		return nil, err
	}

	items := make([]BlacklistItem, 0, len(list))
	if len(list) == 0 {
		return items, nil
	}

	ids := make([]int64, 0, len(list))
	for _, item := range list {
		ids = append(ids, item.BlockedID)
	}
	users, err := dao.NewUserDao().FindByIDs(ids)
	if err != nil {
		return nil, err
	}
	userMap := make(map[int64]model.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	for _, item := range list {
		u, ok := userMap[item.BlockedID]
		if !ok {
			continue // 用户已被删除，跳过
		}
		items = append(items, BlacklistItem{
			ID:       u.ID,
			Username: u.Username,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			AddTime:  item.AddTime,
		})
	}
	return items, nil
}

// AddBlacklist 把目标用户拉入黑名单（重复拉黑幂等）
func AddBlacklist(uid, targetID int64) error {
	if uid <= 0 || targetID <= 0 {
		return ErrUserNotFound
	}
	if uid == targetID {
		return ErrBlacklistSelf
	}
	if _, err := dao.NewUserDao().FindByID(targetID); err != nil {
		return ErrUserNotFound
	}

	blacklistDao := dao.NewBlacklistDao()
	exists, err := blacklistDao.Exists(uid, targetID)
	if err != nil {
		return err
	}
	if !exists {
		if err = blacklistDao.Create(&model.Blacklist{
			UserID:    uid,
			BlockedID: targetID,
			AddTime:   time.Now().Unix(),
		}); err != nil {
			return err
		}
	}

	// 清掉两人之间待处理的申请：否则残留的申请被同意后又能互加好友
	return dao.NewFriendRequestDao().DeletePendingBetween(uid, targetID)
}

// RemoveBlacklist 把目标用户移出黑名单
func RemoveBlacklist(uid, targetID int64) error {
	if uid <= 0 || targetID <= 0 {
		return ErrBlacklistNotFound
	}

	removed, err := dao.NewBlacklistDao().Delete(uid, targetID)
	if err != nil {
		return err
	}
	if !removed {
		return ErrBlacklistNotFound
	}
	return nil
}

// CanSendMessage 判断 fromID 能否把消息发给 toID。
//
// 黑名单只约束「被拉黑方 → 拉黑方」一个方向：toID 拉黑了 fromID 时消息不再送达，
// 反之（自己拉黑的人）不受限制。查询失败时放行，避免数据库抖动误伤正常聊天。
func CanSendMessage(fromID, toID int64) bool {
	if fromID <= 0 || toID <= 0 {
		return false
	}
	blocked, err := dao.NewBlacklistDao().Exists(toID, fromID)
	if err != nil {
		return true
	}
	return !blocked
}

// blockedIDSet 我拉黑的用户 ID 集合。
// 好友列表 / 会话列表据此剔除被拉黑的人（他们只保留在黑名单列表里），
// 搜索结果据此标注关系。这里不设上限：过滤必须完整，漏掉的人会重新冒出来。
func blockedIDSet(uid int64) (map[int64]bool, error) {
	ids, err := dao.NewBlacklistDao().FindBlockedIDs(uid)
	if err != nil {
		return nil, err
	}

	set := make(map[int64]bool, len(ids))
	for _, id := range ids {
		set[id] = true
	}
	return set, nil
}
