// Package conf 负责加载并持有全局配置。
package conf

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"

	"mini-im/pkg/snowflake"
)

// AppConfig 应用基础配置
type AppConfig struct {
	Name     string `yaml:"name"`
	HTTPPort int    `yaml:"http_port"`
	RunMode  string `yaml:"run_mode"`
	WebDir   string `yaml:"web_dir"`
}

// Debug 是否为调试模式
func (a AppConfig) Debug() bool { return a.RunMode == "debug" }

// PostgresConfig PostgreSQL 连接配置
type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	TimeZone string `yaml:"timezone"`
}

// DSN 生成 gorm postgres 驱动所需的数据源字符串
func (p PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode, p.TimeZone,
	)
}

// RedisConfig Redis 连接配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int64  `yaml:"expire"`
}

// SnowflakeConfig 雪花 ID 配置
//
// NodeID 取值 0 ~ 1023，多实例部署时必须各不相同，否则可能生成重复 ID；
// 负数表示按本机网卡地址自动推导。
type SnowflakeConfig struct {
	NodeID int64 `yaml:"node_id"`
}

// SeedAccount 待初始化账号
type SeedAccount struct {
	Username string `yaml:"username"`
	Nickname string `yaml:"nickname"`
}

// SeedConfig 初始化数据配置
type SeedConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Password string        `yaml:"password"`
	Accounts []SeedAccount `yaml:"accounts"`
}

// Config 全局配置
type Config struct {
	App       AppConfig       `yaml:"app"`
	Postgres  PostgresConfig  `yaml:"postgres"`
	Redis     RedisConfig     `yaml:"redis"`
	JWT       JWTConfig       `yaml:"jwt"`
	Snowflake SnowflakeConfig `yaml:"snowflake"`
	Seed      SeedConfig      `yaml:"seed"`
}

// C 全局唯一配置实例（包级共享，只读）
var C = defaults()

func defaults() Config {
	return Config{
		App: AppConfig{
			Name:     "im",
			HTTPPort: 2580,
			RunMode:  "debug",
			WebDir:   "web",
		},
		Postgres: PostgresConfig{
			Host:     "127.0.0.1",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			DBName:   "im",
			SSLMode:  "disable",
			TimeZone: "Asia/Shanghai",
		},
		Redis: RedisConfig{
			Addr: "127.0.0.1:6379",
			DB:   1,
		},
		JWT: JWTConfig{
			Secret: "hjJlbE8LW7NAO15v",
			Expire: 60 * 60 * 24 * 7,
		},
		Snowflake: SnowflakeConfig{
			NodeID: -1, // 跟随本机网卡地址自动推导
		},
	}
}

// Init 读取配置文件，文件不存在时使用内置默认值；
// 随后应用环境变量覆盖，便于容器化部署。
func Init(path string) error {
	data, err := os.ReadFile(path)
	switch {
	case err == nil:
		if err = yaml.Unmarshal(data, &C); err != nil {
			return fmt.Errorf("解析配置文件 %s 失败: %w", path, err)
		}
	case os.IsNotExist(err):
		fmt.Printf("[conf] 未找到配置文件 %s，使用内置默认配置\n", path)
	default:
		return fmt.Errorf("读取配置文件 %s 失败: %w", path, err)
	}

	if C.App.Name == "" {
		C.App.Name = "im"
	}
	if C.JWT.Expire <= 0 {
		C.JWT.Expire = 60 * 60 * 24 * 7
	}
	applyEnv()

	// 节点号越界视为「自动推导」：越界的节点号会被位移丢弃，反而引入重复 ID 风险
	if C.Snowflake.NodeID < 0 || C.Snowflake.NodeID > snowflake.NodeMax {
		C.Snowflake.NodeID = -1
	}
	return nil
}

// applyEnv 用环境变量覆盖关键配置
func applyEnv() {
	setStr("IM_APP_NAME", &C.App.Name)
	setInt("IM_HTTP_PORT", &C.App.HTTPPort)
	setStr("IM_RUN_MODE", &C.App.RunMode)
	setStr("IM_WEB_DIR", &C.App.WebDir)

	setStr("IM_PG_HOST", &C.Postgres.Host)
	setInt("IM_PG_PORT", &C.Postgres.Port)
	setStr("IM_PG_USER", &C.Postgres.User)
	setStr("IM_PG_PASSWORD", &C.Postgres.Password)
	setStr("IM_PG_DBNAME", &C.Postgres.DBName)
	setStr("IM_PG_SSLMODE", &C.Postgres.SSLMode)

	setStr("IM_REDIS_ADDR", &C.Redis.Addr)
	setStr("IM_REDIS_PASSWORD", &C.Redis.Password)
	setInt("IM_REDIS_DB", &C.Redis.DB)

	setStr("IM_JWT_SECRET", &C.JWT.Secret)
	setInt64("IM_JWT_EXPIRE", &C.JWT.Expire)

	setInt64("IM_SNOWFLAKE_NODE_ID", &C.Snowflake.NodeID)
}

func setStr(key string, dst *string) {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		*dst = v
	}
}

func setInt(key string, dst *int) {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			*dst = n
		}
	}
}

func setInt64(key string, dst *int64) {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			*dst = n
		}
	}
}

// ---------------------------------------------------------------
// Redis key 统一在这里生成，避免散落各处的字符串拼接
// ---------------------------------------------------------------

func key(parts ...string) string {
	k := C.App.Name
	for _, p := range parts {
		k += ":" + p
	}
	return k
}

// AuthTokenKey 单点登录 token，value 为当前有效 token
func (Config) AuthTokenKey(uid int64) string {
	return key("auth_token", strconv.FormatInt(uid, 10))
}

// OnlineSetKey 在线用户集合
func (Config) OnlineSetKey() string { return key("online_users") }

// MsgQueueKey 待落库消息队列
func (Config) MsgQueueKey() string { return key("msg_queue") }

// MsgErrorQueueKey 落库失败消息队列
func (Config) MsgErrorQueueKey() string { return key("msg_queue_error") }

// OfflineKey 某用户的离线消息队列
func (Config) OfflineKey(uid int64) string {
	return key("offline", strconv.FormatInt(uid, 10))
}
