// Package setup 实现首次启动时的初始化配置向导：
// 配置文件不存在时，先启动一个只用于填写配置的临时 HTTP 服务，
// 配置写入 config.yaml 后立即关闭，主流程随后按正常方式启动业务服务。
package setup

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"mini-im/conf"
	"mini-im/pkg/port"
)

//go:embed page.html
var pageHTML []byte

// Options 向导启动参数
type Options struct {
	Template   []byte       // 配置模板（config.example.yaml），用于渲染写入文件
	Prefill    *conf.Config // 表单预填值（敏感字段已清空）
	Port       string       // 向导自身监听端口，默认 8280
	ConfigPath string       // 配置写入路径
	Force      bool         // 强制重新配置（目标文件已存在时先备份）
}

type saveRequest struct {
	Config conf.Config `json:"config"`
	Force  bool        `json:"force"`
}

type saveResponse struct {
	OK         bool     `json:"ok"`
	Error      string   `json:"error,omitempty"`
	CanForce   bool     `json:"canForce,omitempty"`
	ConfigPath string   `json:"configPath,omitempty"`
	BackupPath string   `json:"backupPath,omitempty"`
	ServerPort string   `json:"serverPort,omitempty"`
	Warnings   []string `json:"warnings,omitempty"`
}

type stateResponse struct {
	NeedsSetup bool         `json:"needsSetup"`
	Forced     bool         `json:"forced"`
	ConfigPath string       `json:"configPath"`
	SetupPort  string       `json:"setupPort"`
	Existing   bool         `json:"existing"`
	Prefill    *conf.Config `json:"prefill"`
}

// Run 启动配置向导并阻塞，直到配置成功写入（返回 nil）或服务异常（返回错误）。
func Run(opts Options) error {
	opts = normalize(opts)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), quietLogger())

	done := make(chan struct{})

	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/setup") })
	r.GET("/setup", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", pageHTML)
	})
	r.GET("/favicon.ico", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "setup"})
	})

	r.GET("/api/setup/state", func(c *gin.Context) {
		_, statErr := os.Stat(opts.ConfigPath)
		c.JSON(http.StatusOK, stateResponse{
			NeedsSetup: true,
			Forced:     opts.Force,
			ConfigPath: opts.ConfigPath,
			SetupPort:  opts.Port,
			Existing:   statErr == nil,
			Prefill:    opts.Prefill,
		})
	})

	r.POST("/api/setup/test/postgres", func(c *gin.Context) {
		var pg conf.PostgresConfig
		if err := c.ShouldBindJSON(&pg); err != nil {
			c.JSON(http.StatusBadRequest, ProbeResult{Message: "请求参数格式错误", Detail: err.Error()})
			return
		}
		c.JSON(http.StatusOK, probePostgres(pg))
	})

	r.POST("/api/setup/test/redis", func(c *gin.Context) {
		var rd conf.RedisConfig
		if err := c.ShouldBindJSON(&rd); err != nil {
			c.JSON(http.StatusBadRequest, ProbeResult{Message: "请求参数格式错误", Detail: err.Error()})
			return
		}
		c.JSON(http.StatusOK, probeRedis(rd))
	})

	r.POST("/api/setup/save", func(c *gin.Context) {
		var req saveRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, saveResponse{Error: "请求参数格式错误：" + err.Error()})
			return
		}

		cfg := req.Config
		// 表单不含 app.name / app.web_dir / postgres.timezone / snowflake.node_id /
		// seed.accounts（这几项只由配置文件维护），先按模板补齐再校验，
		// 否则会因「未配置」直接拦下保存——而表单里根本没有可填的地方。
		if tpl, err := conf.Parse(opts.Template); err == nil {
			cfg.InheritUnexposed(tpl)
		}
		if err := cfg.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, saveResponse{Error: err.Error()})
			return
		}

		warnings := make([]string, 0, 2)
		if probe := probePostgres(cfg.Postgres); !probe.OK {
			msg := "数据库连接失败：" + probe.Message
			if probe.Detail != "" {
				msg += "（" + probe.Detail + "）"
			}
			if !req.Force {
				c.JSON(http.StatusBadRequest, saveResponse{Error: msg, CanForce: true})
				return
			}
			warnings = append(warnings, msg+"，已按你的要求忽略该问题并保存")
		}
		if probe := probeRedis(cfg.Redis); !probe.OK {
			warnings = append(warnings,
				"Redis 连接失败："+probe.Message+"。本服务的会话与消息投递依赖 Redis，服务启动会失败，请检查后重启")
		}

		backup, err := conf.Save(opts.ConfigPath, opts.Template, &cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, saveResponse{Error: err.Error()})
			return
		}

		c.Header("Connection", "close")
		c.JSON(http.StatusOK, saveResponse{
			OK:         true,
			ConfigPath: opts.ConfigPath,
			BackupPath: backup,
			ServerPort: strconv.Itoa(cfg.App.HTTPPort),
			Warnings:   warnings,
		})

		close(done)
	})

	ln, err := net.Listen("tcp", ":"+opts.Port)
	if err != nil {
		// 端口被系统保留时错误原文只有「权限不足」，补一句换端口的做法
		return fmt.Errorf("配置向导端口 %s 无法监听：%v", opts.Port, port.ListenError(opts.Port, "SETUP_PORT", err))
	}

	srv := &http.Server{Handler: r}
	errCh := make(chan error, 1)
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	log.Printf("未检测到可用配置，已进入初始化配置向导")
	log.Printf("==> 请在浏览器打开：http://localhost:%s/setup", opts.Port)
	for _, ip := range localIPs() {
		log.Printf("==> 远程访问：http://%s:%s/setup", ip, opts.Port)
	}
	log.Printf("==> 配置将写入：%s", opts.ConfigPath)
	if opts.Force {
		log.Printf("==> 已开启强制重配（SETUP_FORCE），原配置会备份为 %s.bak", opts.ConfigPath)
	}

	select {
	case <-done:
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
		log.Printf("初始化配置已保存，正在启动业务服务……")
		return nil
	case err := <-errCh:
		return fmt.Errorf("配置向导运行失败：%w", err)
	}
}

func normalize(opts Options) Options {
	if strings.TrimSpace(opts.Port) == "" {
		opts.Port = "8280"
	}
	if strings.TrimSpace(opts.ConfigPath) == "" {
		opts.ConfigPath = "config.yaml"
	}
	if abs, err := filepath.Abs(opts.ConfigPath); err == nil {
		opts.ConfigPath = abs
	}
	if opts.Prefill == nil {
		opts.Prefill = &conf.Config{}
	}
	return opts
}

// quietLogger 只记录向导接口的简短访问日志，避免刷屏。
func quietLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			log.Printf("[setup] %s %s -> %d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
		}
	}
}

// localIPs 返回本机可用于远程访问的 IPv4 地址，方便在服务器上首次部署时打开向导。
func localIPs() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var ips []string
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.To4() == nil {
				continue
			}
			ips = append(ips, ipNet.IP.String())
		}
	}
	return ips
}
