package service

import (
	"strconv"

	"mini-im/conf"
	"mini-im/db"
)

// SetOnline 标记用户在线
func SetOnline(uid int64) {
	db.Redis().SAdd(db.Ctx, conf.C.OnlineSetKey(), strconv.FormatInt(uid, 10))
}

// SetOffline 标记用户离线
func SetOffline(uid int64) {
	db.Redis().SRem(db.Ctx, conf.C.OnlineSetKey(), strconv.FormatInt(uid, 10))
}
