package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"miniIM/server/conf"
	"time"
)

var Db *gorm.DB

func DbInit() {
	fmt.Println("启动数据库连接", time.Now().Format("2006-01-02 15:04:05"))
	dbHost := conf.DbHost
	dbPort := conf.DbPort
	dbName := conf.DbName
	dbUser := conf.DbUser
	dbPassword := conf.DbPassword
	var err error
	Db, err = gorm.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln("连接数据库失败", err.Error())
	} else {
		fmt.Println("数据库连接成功", dbVersion(), time.Now().Format("2006-01-02 15:04:05"))
	}
	Db.LogMode(true)
	Db.SingularTable(true)
}

func dbVersion() string {
	var rows string
	result := Db.Raw("select version() as version").Row()
	result.Scan(&rows)
	return "mysql " + rows
}
