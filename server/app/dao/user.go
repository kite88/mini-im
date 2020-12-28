package dao

import (
	"github.com/jinzhu/gorm"
	"miniIM/server/app/model"
	"miniIM/server/db"
)

type UserDao struct {}

func (UserDao) D() *gorm.DB {
	var m model.UserModel
	return db.Db.Table(m.TableName())
}
