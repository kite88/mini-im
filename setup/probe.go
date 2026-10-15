package setup

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/redis/go-redis/v9"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"mini-im/conf"
)

const probeTimeout = 5 * time.Second

// ProbeResult 连接测试结果
type ProbeResult struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	Detail    string `json:"detail,omitempty"`
	LatencyMs int64  `json:"latencyMs"`
}

// probePostgres 用提交的表单参数真实连一次 PostgreSQL
// （与 db.PostgresInit 使用同一套 DSN 拼装，见 conf.PostgresConfig.DSN）。
func probePostgres(p conf.PostgresConfig) ProbeResult {
	if strings.TrimSpace(p.Host) == "" || p.Port <= 0 || strings.TrimSpace(p.DBName) == "" {
		return ProbeResult{Message: "请先填写 host / port / dbname"}
	}
	// 表单不含 timezone，探测时按默认补齐：留空会让 pgx 在建连阶段就报 unknown time zone
	if strings.TrimSpace(p.TimeZone) == "" {
		p.TimeZone = "Asia/Shanghai"
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), probeTimeout)
	defer cancel()

	db, err := gorm.Open(postgres.Open(p.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return ProbeResult{Message: "连接失败", Detail: fixEncoding(err.Error()), LatencyMs: elapsed(start)}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return ProbeResult{Message: "连接失败", Detail: fixEncoding(err.Error()), LatencyMs: elapsed(start)}
	}
	defer sqlDB.Close()

	if err := sqlDB.PingContext(ctx); err != nil {
		return ProbeResult{Message: "连接失败，请检查地址 / 账号 / 密码 / 库名", Detail: fixEncoding(err.Error()), LatencyMs: elapsed(start)}
	}

	detail := ""
	var version string
	if err := sqlDB.QueryRowContext(ctx, "select version()").Scan(&version); err == nil {
		detail = firstSentence(version)
	}
	return ProbeResult{OK: true, Message: "连接成功", Detail: detail, LatencyMs: elapsed(start)}
}

// probeRedis 探测 Redis。
// 与 mini-kanban 不同：Redis 对本服务是必需组件（会话 / 消息投递 / 在线状态都依赖它），
// 因此探测失败只是警告，但文案要点明「服务启动会失败」。
func probeRedis(r conf.RedisConfig) ProbeResult {
	if strings.TrimSpace(r.Addr) == "" {
		return ProbeResult{Message: "请先填写 Redis addr（host:port）"}
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), probeTimeout)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})
	defer client.Close()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return ProbeResult{Message: "连接失败，请检查地址 / 密码 / 库号", Detail: err.Error(), LatencyMs: elapsed(start)}
	}
	return ProbeResult{OK: true, Message: "连接成功（" + pong + "）", LatencyMs: elapsed(start)}
}

func elapsed(start time.Time) int64 { return time.Since(start).Milliseconds() }

// fixEncoding 修掉 PostgreSQL 在「连接阶段」下发的非 UTF-8 报错文本。
//
// client_encoding 要等连接建立（选库完成）之后才生效，而认证失败、库不存在
// 这类连接期错误是服务端直接按自己的消息目录编码下发的：中文 Windows 上的
// PostgreSQL（lc_messages = Chinese (Simplified)_China.936）给的是 GBK 字节，
// 在 Go 里就是非法 UTF-8，塞进 JSON 响应会变成一串问号。
//
// 已经是合法 UTF-8 的文本（绝大多数情况，含认证成功后的所有查询报错）原样返回；
// 否则按 GBK 解回来；实在解不干净就丢弃非法字节——保证页面上不出现乱码。
func fixEncoding(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	if out, _, err := transform.String(simplifiedchinese.GBK.NewDecoder(), s); err == nil &&
		!strings.ContainsRune(out, utf8.RuneError) {
		return out
	}
	return strings.ToValidUTF8(s, "")
}

func firstSentence(s string) string {
	s = strings.SplitN(strings.TrimSpace(s), ",", 2)[0]
	if len([]rune(s)) > 80 {
		return fmt.Sprintf("%.80s…", s)
	}
	return s
}
