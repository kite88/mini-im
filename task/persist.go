// Package task 承载后台常驻任务。
package task

import (
	"log"
	"time"

	"mini-im/service"
)

const popTimeout = 5 * time.Second

// Start 启动所有后台任务
func Start() {
	go persistMessages()
	log.Println("[task] 消息落库任务已启动")
}

// persistMessages 从 Redis 队列阻塞读取消息并持久化到 PostgreSQL
//
// 使用 BRPop 而非定时轮询：无消息时阻塞等待，避免空转。
func persistMessages() {
	for {
		msg, ok := service.PopMessageQueue(popTimeout)
		if !ok {
			continue
		}

		if err := service.PersistMessage(msg); err != nil {
			log.Printf("[task] 消息落库失败: %v", err)
			service.PushErrorQueue(msg)
		}
	}
}
