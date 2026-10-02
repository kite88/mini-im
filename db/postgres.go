// Package db 提供 PostgreSQL(GORM) 与 Redis 客户端。
package db

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"mini-im/conf"
	"mini-im/model"
)

// DB 全局 GORM 实例
var DB *gorm.DB

// PostgresInit 建立 PostgreSQL 连接并执行表结构自动迁移
func PostgresInit() error {
	cfg := conf.C

	level := gormlogger.Warn
	if cfg.App.Debug() {
		level = gormlogger.Info
	}

	gdb, err := gorm.Open(postgres.Open(cfg.Postgres.DSN()), &gorm.Config{
		Logger:                 gormlogger.Default.LogMode(level),
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err = gdb.AutoMigrate(&model.User{}, &model.Friend{}, &model.FriendRequest{}, &model.Message{}, &model.ReadState{}, &model.Blacklist{}, &model.ContactState{}); err != nil {
		return err
	}

	DB = gdb
	log.Printf("[postgres] 连接成功 %s:%d/%s", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)
	return nil
}

// Close 关闭数据库连接
func Close() {
	if DB == nil {
		return
	}
	if sqlDB, err := DB.DB(); err == nil {
		_ = sqlDB.Close()
	}
}
