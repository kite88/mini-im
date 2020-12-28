package dao

import (
	"github.com/jinzhu/gorm"
	"miniIM/server/app/model"
	"miniIM/server/db"
)

type FriendDao struct {}

func (FriendDao) D() *gorm.DB {
	var m model.FriendModel
	return db.Db.Table(m.TableName())
}