package dao

import (
	"github.com/jinzhu/gorm"
	"miniIM/server/app/model"
	"miniIM/server/db"
)

type MessageDao struct{}

func (MessageDao) D() *gorm.DB {
	var m model.MessageModel
	return db.Db.Table(m.TableName())
}
