package service

import (
	"miniIM/server/conf"
	"miniIM/server/db"
)

func LineUser(uid string, is bool) {
	if is == true {
		db.Redis().SAdd(conf.ImLineUserKey, uid)
	} else {
		db.Redis().SRem(conf.ImLineUserKey, uid)
	}
}

func IsUserLine(uid string) bool {
	res, _ := db.Redis().SIsMember(conf.ImLineUserKey, conf.ImUserPrefix+uid).Result()
	return res
}
