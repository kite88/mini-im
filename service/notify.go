package service

import "encoding/json"

// WebSocket 下行事件名（报文形如 {"event": "...", "data": {...}}）
const (
	EventFriendRequest  = "friend_request"  // 收到新的好友申请
	EventFriendAccepted = "friend_accepted" // 好友申请被同意（互加成功）
	EventNotFriend      = "not_friend"      // 消息因我已删除对方（或双方不是好友）而未送达
	EventUnfriended     = "unfriended"      // 消息因对方已把我从好友中删除而未送达
	EventBlacklisted    = "blacklisted"     // 消息因被对方拉黑而未能送达
)

// Notifier 由 WebSocket 连接中心实现（*ws.Hub）。
//
// 接口定义在 service 侧、实现放在 ws 侧，是为了保持依赖方向单向：
// ws -> service，业务层不必反向引用 ws。注入动作在 main 中完成。
type Notifier interface {
	SendTo(uid int64, payload []byte) bool
}

var notifier Notifier

// SetNotifier 注入事件推送实现，进程启动时调用一次
func SetNotifier(n Notifier) { notifier = n }

// Notify 向指定用户推送一条事件报文。
//
// 对方离线时直接丢弃：好友申请本身已落库，对方下次登录拉列表依然能看到。
func Notify(uid int64, event string, data interface{}) {
	if notifier == nil || uid <= 0 {
		return
	}
	payload, err := json.Marshal(map[string]interface{}{
		"event": event,
		"data":  data,
	})
	if err != nil {
		return
	}
	notifier.SendTo(uid, payload)
}
