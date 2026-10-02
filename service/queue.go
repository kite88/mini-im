package service

import (
	"encoding/json"
	"time"

	"mini-im/conf"
	"mini-im/dao"
	"mini-im/db"
	"mini-im/model"
)

// PushMessageQueue 投递消息到「待落库」队列，由后台任务异步写入 PostgreSQL
func PushMessageQueue(msg model.Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return db.Redis().LPush(db.Ctx, conf.C.MsgQueueKey(), payload).Err()
}

// PopMessageQueue 阻塞读取一条待落库消息
func PopMessageQueue(timeout time.Duration) (model.Message, bool) {
	var msg model.Message

	res, err := db.Redis().BRPop(db.Ctx, timeout, conf.C.MsgQueueKey()).Result()
	if err != nil || len(res) < 2 {
		return msg, false
	}
	if err = json.Unmarshal([]byte(res[1]), &msg); err != nil {
		return msg, false
	}
	return msg, true
}

// PushErrorQueue 落库失败的消息转入错误队列，便于人工补偿
func PushErrorQueue(msg model.Message) {
	payload, err := json.Marshal(msg)
	if err != nil {
		return
	}
	db.Redis().LPush(db.Ctx, conf.C.MsgErrorQueueKey(), payload)
}

// PersistMessage 将消息持久化到数据库
//
// 主键由 model.Message 的 BeforeCreate 钩子填充 UUID v7，这里只清理两个不落库的展示字段。
func PersistMessage(msg model.Message) error {
	msg.Type = 0
	msg.User = nil
	return dao.NewMessageDao().Create(&msg)
}

// PushOffline 缓存离线消息
//
// 调用方（WebSocket 分发）已确认对方无活跃连接，这里不再重复判断在线状态，
// 以免出现「在线集合与真实连接不一致」导致消息丢失。
func PushOffline(uid int64, payload []byte) {
	if uid <= 0 || len(payload) == 0 {
		return
	}
	db.Redis().LPush(db.Ctx, conf.C.OfflineKey(uid), payload)
}

// DrainOffline 取出并清空离线消息（LPUSH + RPop 保证原始顺序）
func DrainOffline(uid int64) []string {
	key := conf.C.OfflineKey(uid)

	length, err := db.Redis().LLen(db.Ctx, key).Result()
	if err != nil || length == 0 {
		return nil
	}

	list := make([]string, 0, length)
	for i := int64(0); i < length; i++ {
		value, err := db.Redis().RPop(db.Ctx, key).Result()
		if err != nil {
			break
		}
		list = append(list, value)
	}
	return list
}
