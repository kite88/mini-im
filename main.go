// mini-im —— 基于 Go + Gin + GORM(PostgreSQL) + Redis 的即时通讯 / 客服系统
package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	// 内嵌 IANA 时区数据库：CGO_ENABLED=0 交叉编译出的归档，在没装 Go、
	// 系统也没有 IANA tzdata 的机器上（典型是 Windows）同样能解析配置中的
	// Asia/Shanghai，否则 pgx 建连时直接报 unknown time zone
	_ "time/tzdata"

	"github.com/gin-gonic/gin"

	"mini-im/api/middleware"
	"mini-im/conf"
	"mini-im/db"
	"mini-im/pkg/port"
	"mini-im/pkg/snowflake"
	"mini-im/route"
	"mini-im/service"
	"mini-im/setup"
	"mini-im/task"
	"mini-im/ws"
)

// configTemplate 配置模板，用作初始化向导写入配置文件时的底稿（保留注释与排版）
//
//go:embed config.example.yaml
var configTemplate []byte

func main() {
	configPath := flag.String("c", "config.yaml", "配置文件路径")
	flag.Parse()

	// 配置文件缺失，或显式设置 SETUP_FORCE 时先跑初始化向导。
	// 注意 SETUP_FORCE 必须独立于加载结果判断：配置本身是好的，只是想重新配置。
	force := isTruthy(os.Getenv("SETUP_FORCE"))
	if err := conf.Init(*configPath); err != nil || force {
		if err := loadOrSetup(*configPath, err, force); err != nil {
			log.Fatalf("%v", err)
		}
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
			// 端口被系统保留 / 被占用时，原始错误只是一句「权限不足」，
			// 这里补上换端口的具体做法
			log.Fatalf("HTTP 服务异常退出: %v",
				port.ListenError(strconv.Itoa(cfg.App.HTTPPort), "IM_HTTP_PORT", err))
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

// loadOrSetup 处理「配置缺失」与「强制重配」两种情形：
// 先启动初始化配置向导，待配置文件写入完成后重新加载，再由调用方按正常流程启动。
//
// 只在配置文件不存在（loadErr 为 ErrConfigNotFound），或调用方判定 SETUP_FORCE 生效时
// 进入；其余加载错误（格式非法等）直接返回给调用方报错退出，
// 免得把「配置写坏了」误当成「需要重新配置」。
func loadOrSetup(path string, loadErr error, force bool) error {
	notFound := errors.Is(loadErr, conf.ErrConfigNotFound)
	if !notFound && loadErr != nil && !force {
		return loadErr
	}
	if !notFound && loadErr != nil {
		log.Printf("现有配置存在问题（SETUP_FORCE 已开启，将在向导中重新配置）：%v", loadErr)
	}

	// 预填表单：首次启动用模板里的示例值，强制重配时用当前生效的配置
	// （Sanitize 会清掉密码与密钥，不回显到页面上）
	prefill := &conf.Config{}
	if notFound {
		if parsed, perr := conf.Parse(configTemplate); perr == nil {
			prefill = parsed
		}
	} else {
		current := conf.C
		prefill = current.Sanitize()
	}

	if err := setup.Run(setup.Options{
		Template:   configTemplate,
		Prefill:    prefill.Sanitize(),
		Port:       setupPort(),
		ConfigPath: path,
		Force:      force,
	}); err != nil {
		return err
	}

	if err := conf.Init(path); err != nil {
		return fmt.Errorf("初始化配置已写入，但加载失败：%w", err)
	}
	return nil
}

// setupPort 初始化向导自身的监听端口：SETUP_PORT 环境变量优先，默认 8280
// （此时还没有可用配置，无法从配置里读端口）。
func setupPort() string {
	if v := strings.TrimSpace(os.Getenv("SETUP_PORT")); v != "" {
		return v
	}
	return "8280"
}

func isTruthy(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true
	}
	return false
}
