package service

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"mini-im/dao"
	"mini-im/model"
)

const (
	// friendSearchLimit 单次搜索返回的最大用户数
	friendSearchLimit = 20
	// friendRequestLimit 单次拉取的待处理申请上限
	friendRequestLimit = 50
	// friendRequestMessageMax 验证消息最大长度（字符数）
	friendRequestMessageMax = 100
)

// 搜索结果中「我」与目标用户的关系，供前端决定按钮形态
const (
	RelationNone      = "none"      // 可添加
	RelationFriend    = "friend"    // 已是好友
	RelationRequested = "requested" // 我已申请，等对方处理
	RelationIncoming  = "incoming"  // 对方申请了我，待我处理
	RelationBlocked   = "blocked"   // 我已把对方拉黑，需先移出黑名单
)

var (
	// ErrUserNotFound 目标用户不存在
	ErrUserNotFound = errors.New("用户不存在")
	// ErrFriendSelf 不能添加自己
	ErrFriendSelf = errors.New("不能添加自己为好友")
	// ErrAlreadyFriend 已经是好友
	ErrAlreadyFriend = errors.New("你们已经是好友了")
	// ErrFriendNotFound 双方不是好友（或好友关系已被删除）
	ErrFriendNotFound = errors.New("你们还不是好友")
	// ErrFriendRequestNotFound 申请不存在或无权处理
	ErrFriendRequestNotFound = errors.New("好友申请不存在")
	// ErrFriendRequestHandled 申请已被处理
	ErrFriendRequestHandled = errors.New("该好友申请已处理")
	// ErrSearchFailed 搜索失败
	ErrSearchFailed = errors.New("搜索用户失败")
)

// UserItem 搜索结果项：用户公开资料 + 与当前用户的关系
type UserItem struct {
	ID       int64  `json:"id,string"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Relation string `json:"relation"`
}

// RequestItem 好友申请项（申请人为谁 + 验证消息）
type RequestItem struct {
	ID         int64      `json:"id,string"`
	FromID     int64      `json:"from_id,string"`
	Message    string     `json:"message"`
	CreateTime int64      `json:"create_time"`
	User       model.User `json:"user"`
}

// SearchUsers 按账号 / 昵称搜索用户（排除自己），并标注与当前用户的关系
func SearchUsers(uid int64, keyword string) ([]UserItem, error) {
	keyword = strings.TrimSpace(keyword)
	items := make([]UserItem, 0, friendSearchLimit)
	if keyword == "" {
		return items, nil
	}

	users, err := dao.NewUserDao().SearchByKeyword(keyword, uid, friendSearchLimit)
	if err != nil {
		return nil, ErrSearchFailed
	}
	if len(users) == 0 {
		return items, nil
	}

	ids := make([]int64, 0, len(users))
	for _, u := range users {
		ids = append(ids, u.ID)
	}

	friends, err := friendIDSet(uid)
	if err != nil {
		return nil, ErrSearchFailed
	}

	pending, err := dao.NewFriendRequestDao().FindPendingWith(uid, ids)
	if err != nil {
		return nil, ErrSearchFailed
	}
	pendingMap := make(map[int64]model.FriendRequest, len(pending))
	for _, r := range pending {
		if r.FromID == uid {
			pendingMap[r.ToID] = r
			continue
		}
		pendingMap[r.FromID] = r
	}

	// 已拉黑的用户优先标注：此时不允许再发申请（后端也会拦）
	blocked, err := blockedIDSet(uid)
	if err != nil {
		return nil, ErrSearchFailed
	}

	for _, u := range users {
		relation := RelationNone
		switch r, ok := pendingMap[u.ID]; {
		case blocked[u.ID]:
			relation = RelationBlocked
		case friends[u.ID]:
			relation = RelationFriend
		case ok && r.FromID == uid:
			relation = RelationRequested
		case ok:
			relation = RelationIncoming
		}
		items = append(items, UserItem{
			ID:       u.ID,
			Username: u.Username,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Relation: relation,
		})
	}
	return items, nil
}

// SendFriendRequest 发起好友申请。
//
// 返回 true 表示无需等待对方处理就已互为好友：对方此前申请过我，
// 此时再次申请等同于双方互相确认。
func SendFriendRequest(uid, toID int64, message string) (bool, error) {
	if uid <= 0 || toID <= 0 {
		return false, ErrUserNotFound
	}
	if uid == toID {
		return false, ErrFriendSelf
	}

	userDao := dao.NewUserDao()
	if _, err := userDao.FindByID(toID); err != nil {
		return false, ErrUserNotFound
	}

	friendDao := dao.NewFriendDao()
	exists, err := friendDao.ExistsBetween(uid, toID)
	if err != nil {
		return false, err
	}
	if exists {
		// 好友记录只在首次加好友时写入，删除好友不会删行：
		// 只有「我已删除对方 / 对方已删除我」时才算真正重新申请
		removed, err := removedBetween(uid, toID)
		if err != nil {
			return false, err
		}
		if !removed {
			return false, ErrAlreadyFriend
		}
	}

	// 黑名单：我拉黑了对方要先移出；对方拉黑了我则一律拒绝
	blacklistDao := dao.NewBlacklistDao()
	mine, err := blacklistDao.Exists(uid, toID)
	if err != nil {
		return false, err
	}
	if mine {
		return false, ErrBlacklistBlockedPeer
	}
	theirs, err := blacklistDao.Exists(toID, uid)
	if err != nil {
		return false, err
	}
	if theirs {
		return false, ErrBlacklistRejected
	}

	now := time.Now().Unix()
	message = trimFriendMessage(message)
	reqDao := dao.NewFriendRequestDao()

	// 对方已申请过我：双向确认，直接成为好友
	incoming, err := reqDao.FindPending(toID, uid)
	switch {
	case err == nil:
		if err = createFriendRelation(uid, toID, now); err != nil {
			return false, err
		}
		if _, err = reqDao.UpdatePending(incoming.ID, model.FriendRequestAccepted, now); err != nil {
			return false, err
		}
		notifyUser(toID, uid, EventFriendAccepted)
		return true, nil
	case errors.Is(err, dao.ErrNotFound):
		// 正常路径：写入申请
	default:
		return false, err
	}

	if err = reqDao.Upsert(&model.FriendRequest{
		FromID:     uid,
		ToID:       toID,
		Message:    message,
		Status:     model.FriendRequestPending,
		CreateTime: now,
		UpdateTime: now,
	}); err != nil {
		return false, err
	}

	notifyUser(toID, uid, EventFriendRequest, message)
	return false, nil
}

// FriendRequestList 待「我」处理的申请列表（含申请人资料）
func FriendRequestList(uid int64) ([]RequestItem, error) {
	list, err := dao.NewFriendRequestDao().FindPendingByTo(uid, friendRequestLimit)
	if err != nil {
		return nil, err
	}

	items := make([]RequestItem, 0, len(list))
	if len(list) == 0 {
		return items, nil
	}

	ids := make([]int64, 0, len(list))
	for _, r := range list {
		ids = append(ids, r.FromID)
	}
	users, err := dao.NewUserDao().FindByIDs(ids)
	if err != nil {
		return nil, err
	}
	userMap := make(map[int64]model.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	for _, r := range list {
		items = append(items, RequestItem{
			ID:         r.ID,
			FromID:     r.FromID,
			Message:    r.Message,
			CreateTime: r.CreateTime,
			User:       userMap[r.FromID],
		})
	}
	return items, nil
}

// HandleFriendRequest 同意 / 拒绝一条好友申请
func HandleFriendRequest(uid, requestID int64, accept bool) error {
	if uid <= 0 || requestID <= 0 {
		return ErrFriendRequestNotFound
	}

	reqDao := dao.NewFriendRequestDao()
	req, err := reqDao.FindByID(requestID)
	if err != nil {
		if errors.Is(err, dao.ErrNotFound) {
			return ErrFriendRequestNotFound
		}
		return err
	}
	// 只能处理「别人发给我」的申请，其他人的申请一律视为不存在
	if req.ToID != uid {
		return ErrFriendRequestNotFound
	}
	if req.Status != model.FriendRequestPending {
		return ErrFriendRequestHandled
	}

	now := time.Now().Unix()
	status := model.FriendRequestRejected
	if accept {
		status = model.FriendRequestAccepted
	}

	updated, err := reqDao.UpdatePending(req.ID, status, now)
	if err != nil {
		return err
	}
	if !updated {
		return ErrFriendRequestHandled
	}

	if !accept {
		return nil
	}

	if err = createFriendRelation(req.FromID, req.ToID, now); err != nil {
		return err
	}
	notifyUser(req.FromID, uid, EventFriendAccepted)
	return nil
}

// createFriendRelation 建立双向好友关系（一条记录即表达双向），已存在时幂等返回。
//
// 双方重新加为好友时清掉各自的「已删除」标记，但保留会话清空时间：
// 删除过的一方仍然是空会话，对方那侧的历史记录也不会被翻出来。
func createFriendRelation(uid, fid, now int64) error {
	if uid > fid {
		uid, fid = fid, uid
	}

	friendDao := dao.NewFriendDao()
	exists, err := friendDao.ExistsBetween(uid, fid)
	if err != nil {
		return err
	}
	if !exists {
		if err = friendDao.Create(&model.Friend{UserID: uid, FriendID: fid, AddTime: now}); err != nil {
			return err
		}
	}

	return dao.NewContactStateDao().Restore(uid, fid)
}

// contactStateOf 我这一侧的联系人状态：
// removed = 我已从好友中移除的人；clearTime = 我清空过会话的人及清空时间点。
// 好友列表 / 会话列表 / 聊天记录都据此过滤，只影响我自己的展示。
func contactStateOf(uid int64) (removed map[int64]bool, clearTime map[int64]int64, err error) {
	states, err := dao.NewContactStateDao().FindByUser(uid)
	if err != nil {
		return nil, nil, err
	}

	removed = make(map[int64]bool, len(states))
	clearTime = make(map[int64]int64, len(states))
	for _, s := range states {
		if s.Removed {
			removed[s.PeerID] = true
		}
		if s.ClearTime > 0 {
			clearTime[s.PeerID] = s.ClearTime
		}
	}
	return removed, clearTime, nil
}

// removedBetween 两人之间是否有任意一方删除了好友（用于判断能否重新申请）
func removedBetween(a, b int64) (bool, error) {
	states, err := dao.NewContactStateDao().FindPair(a, b)
	if err != nil {
		return false, err
	}
	for _, s := range states {
		if s.Removed {
			return true, nil
		}
	}
	return false, nil
}

// RemovedByPeer 对方是否已把我从好友中删除（用于给发送方更准确的提示）
func RemovedByPeer(uid, peerID int64) bool {
	states, err := dao.NewContactStateDao().FindPair(uid, peerID)
	if err != nil {
		return false
	}
	for _, s := range states {
		if s.Removed && s.UserID == peerID {
			return true
		}
	}
	return false
}

// IsFriend 双方是否还能互发消息：任何一方把对方从好友中移除都不行。
//
// 查询失败时放行，避免数据库抖动导致正常聊天发不出去。
func IsFriend(uid, fid int64) bool {
	if uid <= 0 || fid <= 0 || uid == fid {
		return false
	}

	removed, err := removedBetween(uid, fid)
	if err != nil {
		return true
	}
	if removed {
		return false
	}

	exists, err := dao.NewFriendDao().ExistsBetween(uid, fid)
	if err != nil {
		return true
	}
	return exists
}

// DeleteFriend 删除好友：只改「我这一侧」的状态（与拉黑一样是单向行为）。
//
// 我的好友列表与会话列表不再有对方，我的聊天记录从此刻起不可见（写入清空时间）；
// 对方的列表与历史记录都不受影响，对方只是发不出消息，并在发送时收到提醒。
// 双方重新加为好友后，我这边是全新会话，旧消息仍保留在对方那边。
// im_friend 记录不删除，删除状态记在 im_contact_state。
func DeleteFriend(uid, fid int64) error {
	if uid <= 0 || fid <= 0 {
		return ErrFriendNotFound
	}

	exists, err := dao.NewFriendDao().ExistsBetween(uid, fid)
	if err != nil {
		return err
	}
	if !exists {
		return ErrFriendNotFound
	}

	return dao.NewContactStateDao().UpsertRemoved(uid, fid, time.Now().Unix())
}

// friendIDSet 当前用户的有效好友 ID 集合（已被我删除的不算）
func friendIDSet(uid int64) (map[int64]bool, error) {
	relations, err := dao.NewFriendDao().FindByUser(uid)
	if err != nil {
		return nil, err
	}

	removed, _, err := contactStateOf(uid)
	if err != nil {
		return nil, err
	}

	set := make(map[int64]bool, len(relations)*2)
	for _, r := range relations {
		for _, id := range [...]int64{r.UserID, r.FriendID} {
			if id == 0 || id == uid || removed[id] {
				continue
			}
			set[id] = true
		}
	}
	return set, nil
}

// trimFriendMessage 规整验证消息：去空白 + 截断到上限
func trimFriendMessage(message string) string {
	message = strings.TrimSpace(message)
	if message == "" {
		return ""
	}
	runes := []rune(message)
	if len(runes) > friendRequestMessageMax {
		return string(runes[:friendRequestMessageMax])
	}
	return message
}

// notifyUser 把 fromID 的简要资料推送给 toID（离线则丢弃）
func notifyUser(toID, fromID int64, event string, extra ...string) {
	from, err := dao.NewUserDao().FindByID(fromID)
	if err != nil {
		return
	}

	data := map[string]interface{}{
		"id":       strconv.FormatInt(from.ID, 10),
		"username": from.Username,
		"nickname": from.Nickname,
		"avatar":   from.Avatar,
	}
	if len(extra) > 0 {
		data["message"] = extra[0]
	}
	Notify(toID, event, data)
}
