// mini-im —— 基于 Go + Gin + GORM(PostgreSQL) + Redis 的即时通讯 / 客服系统
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"mini-im/api/middleware"
	"mini-im/conf"
	"mini-im/db"
	"mini-im/pkg/snowflake"
	"mini-im/route"
	"mini-im/service"
	"mini-im/task"
	"mini-im/ws"
)

func main() {
	configPath := flag.String("c", "config.yaml", "配置文件路径")
	flag.Parse()

	if err := conf.Init(*configPath); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg := conf.C
	log.Printf("mini-im 启动中... 模式=%s 端口=%d", cfg.App.RunMode, cfg.App.HTTPPort)
	gin.SetMode(ginMode(cfg.App.RunMode))

	// 雪花 ID：用户表 / 消息表主键在应用侧生成，必须在写入任何数据（含 seed）之前就绪
	nodeID := snowflake.Init(cfg.Snowflake.NodeID)
	log.Printf("[snowflake] 节点号=%d", nodeID)

	// 基础组件
	if err := db.RedisInit(); err != nil {
		log.Fatalf("连接 Redis 失败: %v", err)
	}
	if err := db.PostgresInit(); err != nil {
		log.Fatalf("连接 PostgreSQL 失败: %v", err)
	}
	if err := db.Seed(); err != nil {
		log.Fatalf("初始化演示数据失败: %v", err)
	}

	// 后台任务
	hub := ws.NewHub()
	// 把连接中心注入业务层：好友申请、同意等事件由 service 主动推送给在线用户
	service.SetNotifier(hub)
	task.Start()

	// HTTP 服务
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery(), middleware.CORS())
	route.Register(engine, hub)

	srv := &http.Server{
		Addr:        ":" + strconv.Itoa(cfg.App.HTTPPort),
		Handler:     engine,
		ReadTimeout: 30 * time.Second,
		// WebSocket 是长连接，不能设置写超时，否则会被强制中断
		WriteTimeout: 0,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP 服务异常退出: %v", err)
		}
	}()
	log.Printf("HTTP 服务已就绪: http://127.0.0.1:%d/web/index.html", cfg.App.HTTPPort)

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("收到退出信号，正在优雅关闭...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP 服务关闭异常: %v", err)
	}
	db.Close()
	log.Println("已退出")
}

// ginMode 将配置中的运行模式映射为 gin 常量，非法值回退到 debug
func ginMode(mode string) string {
	switch mode {
	case gin.ReleaseMode, gin.TestMode, gin.DebugMode:
		return mode
	default:
		return gin.DebugMode
	}
}
