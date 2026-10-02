package service

import (
	"time"

	"mini-im/dao"
	"mini-im/model"
)

const (
	// chatHistoryLimit 单次会话拉取的最大消息数
	chatHistoryLimit = 200
	// conversationRecentLimit 计算会话列表时扫描的最近消息数
	conversationRecentLimit = 500
)

// FriendList 好友列表（不含自己）
//
// 两类人不出现在列表里，且都只影响我这一侧的展示：
//   - 我主动拉黑的人（只保留在黑名单列表，移出后恢复）
//   - 我主动删除的人（对方删我不影响我的列表，只影响能否发消息）
func FriendList(uid int64) ([]model.User, error) {
	relations, err := dao.NewFriendDao().FindByUser(uid)
	if err != nil {
		return nil, err
	}

	blocked, err := blockedIDSet(uid)
	if err != nil {
		return nil, err
	}

	removed, _, err := contactStateOf(uid)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(relations)*2)
	seen := make(map[int64]struct{}, len(relations)*2)
	for _, r := range relations {
		for _, id := range [...]int64{r.UserID, r.FriendID} {
			if id == 0 || id == uid || blocked[id] || removed[id] {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
	}

	return dao.NewUserDao().FindByIDs(ids)
}

// MessageItem 会话列表项
//
// Uid 是对方的雪花 ID，按字符串输出，避免前端 JavaScript 精度丢失。
type MessageItem struct {
	Time        int64      `json:"time"`
	Message     string     `json:"message"`
	ContentType string     `json:"content_type"`
	Uid         int64      `json:"uid,string"`
	MsgCount    int64      `json:"msg_count"`
	User        model.User `json:"user"`
}

// MessageList 会话列表：每个聊天对象只保留最新一条消息，按时间倒序
//
// 三类会话不出现在列表里，且都只影响我这一侧：
//   - 我主动拉黑的人（消息记录保留，移出黑名单后重新出现）
//   - 我主动删除的人（会话记录被我清空，对方的会话不受影响）
//   - 我清空过会话、且清空之后没有新消息的人
func MessageList(uid int64) ([]MessageItem, error) {
	msgs, err := dao.NewMessageDao().FindRecentByUser(uid, conversationRecentLimit)
	if err != nil {
		return nil, err
	}

	blocked, err := blockedIDSet(uid)
	if err != nil {
		return nil, err
	}

	removed, clearTime, err := contactStateOf(uid)
	if err != nil {
		return nil, err
	}

	items := make([]MessageItem, 0, len(msgs))
	partnerIDs := make([]int64, 0, len(msgs))
	seen := make(map[int64]struct{}, len(msgs))

	for _, m := range msgs {
		partner := m.FriendID
		if m.FriendID == uid {
			partner = m.UserID
		}
		if partner == 0 || partner == uid || blocked[partner] || removed[partner] {
			continue
		}
		// 已经清空过这个会话：清空时间之前的消息不再出现（之后的会正常展示）
		if since, ok := clearTime[partner]; ok && m.CreateTime <= since {
			continue
		}
		if _, ok := seen[partner]; ok {
			continue
		}
		seen[partner] = struct{}{}
		partnerIDs = append(partnerIDs, partner)
		items = append(items, MessageItem{
			Time:        m.CreateTime,
			Message:     m.Content,
			ContentType: m.ContentType,
			Uid:         partner,
		})
	}

	users, err := dao.NewUserDao().FindByIDs(partnerIDs)
	if err != nil {
		return nil, err
	}
	userMap := make(map[int64]model.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}
	for i := range items {
		if u, ok := userMap[items[i].Uid]; ok {
			items[i].User = u
		}
	}

	// 未读数由服务端已读游标计算，前端刷新页面 / 重新登录后靠它恢复角标
	unread, err := dao.NewReadStateDao().UnreadCounts(uid)
	if err != nil {
		return nil, err
	}
	for i := range items {
		if count, ok := unread[items[i].Uid]; ok {
			items[i].MsgCount = count
		}
	}

	return items, nil
}

// MarkConversationRead 推进「uid 已读完与 fid 的会话」的读游标
func MarkConversationRead(uid, fid int64) error {
	if uid <= 0 || fid <= 0 || uid == fid {
		return nil
	}
	return dao.NewReadStateDao().Upsert(uid, fid, time.Now().Unix())
}

// ChatList 与指定好友的聊天记录，按时间倒序返回（前端会 reverse 后正序渲染）
//
// 我清空过会话时只返回清空时间之后的消息；对方清空不影响我这边。
func ChatList(uid, fid int64) ([]model.Message, error) {
	since, err := dao.NewContactStateDao().ClearTime(uid, fid)
	if err != nil {
		return nil, err
	}

	msgs, err := dao.NewMessageDao().FindBetween(uid, fid, chatHistoryLimit, since)
	if err != nil {
		return nil, err
	}

	users, err := dao.NewUserDao().FindByIDs([]int64{uid, fid})
	if err != nil {
		return nil, err
	}
	userMap := make(map[int64]model.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	for i := range msgs {
		if msgs[i].UserID == uid {
			msgs[i].Type = model.MsgTypeSelf
		} else {
			msgs[i].Type = model.MsgTypePeer
		}
		if u, ok := userMap[msgs[i].UserID]; ok {
			user := u
			msgs[i].User = &user
		}
	}

	return msgs, nil
}

// GetUserBrief 按 ID 查询用户（WebSocket 补全发送者资料用）
func GetUserBrief(uid int64) (*model.User, error) {
	return dao.NewUserDao().FindByID(uid)
}
