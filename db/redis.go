package db

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"mini-im/conf"
)

// RDB 全局 Redis 客户端
var RDB *redis.Client

// Ctx 默认上下文（Redis / 后台任务共用）
var Ctx = context.Background()

// RedisInit 建立 Redis 连接并做一次心跳校验
func RedisInit() error {
	cfg := conf.C.Redis

	RDB = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolSize:     50,
	})

	ctx, cancel := context.WithTimeout(Ctx, 5*time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		return err
	}
	log.Printf("[redis] 连接成功 %s (db=%d)", cfg.Addr, cfg.DB)
	return nil
}

// Redis 返回 Redis 客户端
func Redis() *redis.Client { return RDB }
