// Package conf 负责加载并持有全局配置。
package conf

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"mini-im/pkg/snowflake"
)

// ErrConfigNotFound 配置文件不存在（首次启动）。
// main 据此决定进入初始化配置向导，而不是直接用内置默认值启动 —— 默认值只面向
// 本地开发，发布包里「连不上数据库还照常启动」没有任何意义。
var ErrConfigNotFound = errors.New("配置文件不存在")

// FilePath 当前使用的配置文件路径，由 Init 记录，供初始化向导展示与写回。
var FilePath = "config.yaml"

// AppConfig 应用基础配置
type AppConfig struct {
	Name     string `yaml:"name" json:"name"`
	HTTPPort int    `yaml:"http_port" json:"http_port"`
	RunMode  string `yaml:"run_mode" json:"run_mode"`
	WebDir   string `yaml:"web_dir" json:"web_dir"`
}

// Debug 是否为调试模式
func (a AppConfig) Debug() bool { return a.RunMode == "debug" }

// PostgresConfig PostgreSQL 连接配置
type PostgresConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	DBName   string `yaml:"dbname" json:"dbname"`
	SSLMode  string `yaml:"sslmode" json:"sslmode"`
	TimeZone string `yaml:"timezone" json:"timezone"`
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
	Addr     string `yaml:"addr" json:"addr"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret string `yaml:"secret" json:"secret"`
	Expire int64  `yaml:"expire" json:"expire"` // token 有效期（秒）
}

// SnowflakeConfig 雪花 ID 配置
//
// NodeID 取值 0 ~ 1023，多实例部署时必须各不相同，否则可能生成重复 ID；
// 负数表示按本机网卡地址自动推导。
type SnowflakeConfig struct {
	NodeID int64 `yaml:"node_id" json:"node_id"`
}

// SeedAccount 待初始化账号
type SeedAccount struct {
	Username string `yaml:"username" json:"username"`
	Nickname string `yaml:"nickname" json:"nickname"`
}

// SeedConfig 初始化数据配置
type SeedConfig struct {
	Enabled  bool          `yaml:"enabled" json:"enabled"`
	Password string        `yaml:"password" json:"password"`
	Accounts []SeedAccount `yaml:"accounts" json:"accounts"`
}

// Config 全局配置
type Config struct {
	App       AppConfig       `yaml:"app" json:"app"`
	Postgres  PostgresConfig  `yaml:"postgres" json:"postgres"`
	Redis     RedisConfig     `yaml:"redis" json:"redis"`
	JWT       JWTConfig       `yaml:"jwt" json:"jwt"`
	Snowflake SnowflakeConfig `yaml:"snowflake" json:"snowflake"`
	Seed      SeedConfig      `yaml:"seed" json:"seed"`
}

// C 全局唯一配置实例（包级共享，只读）
var C = defaults()

func defaults() Config {
	return Config{
		App: AppConfig{
			Name:     "im",
			HTTPPort: 8280,
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

// Init 读取配置文件；随后应用环境变量覆盖，便于容器化部署。
//
// 配置文件不存在时返回 ErrConfigNotFound（调用方据此进入初始化配置向导），
// 内容非法时返回具体错误。内置默认值只在配置项缺失时兜底，不代表「无需配置文件」。
func Init(path string) error {
	FilePath = path

	data, err := os.ReadFile(path)
	switch {
	case err == nil:
		if err = yaml.Unmarshal(data, &C); err != nil {
			return fmt.Errorf("解析配置文件 %s 失败: %w", path, err)
		}
	case os.IsNotExist(err):
		return fmt.Errorf("%w: %s", ErrConfigNotFound, path)
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

// Parse 只解析配置内容，不做环境变量覆盖，也不套用内置默认值
// （供初始化向导解析模板与校验表单使用）。
func Parse(data []byte) (*Config, error) {
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Sanitize 清空敏感字段，避免已保存的密码 / 密钥被回显到初始化向导页面。
func (c *Config) Sanitize() *Config {
	c.Postgres.Password = ""
	c.Redis.Password = ""
	c.JWT.Secret = ""
	c.Seed.Password = ""
	return c
}

// InheritUnexposed 用模板（config.example.yaml）补齐「初始化向导不展示」的配置项。
//
// 向导表单只包含 app.http_port / app.run_mode、postgres（除 timezone）、redis、jwt 与
// seed.enabled / seed.password；app.name（Redis key 前缀）、app.web_dir、postgres.timezone、
// snowflake.node_id 与 seed.accounts 刻意不进表单，它们只由配置文件维护
// —— 表单提交上来的对象里这几处是零值，直接 Validate 会误报「未配置」而无法保存。
func (c *Config) InheritUnexposed(tpl *Config) {
	if strings.TrimSpace(c.App.Name) == "" {
		c.App.Name = tpl.App.Name
	}
	if strings.TrimSpace(c.App.WebDir) == "" {
		c.App.WebDir = tpl.App.WebDir
	}
	if strings.TrimSpace(c.Postgres.TimeZone) == "" {
		c.Postgres.TimeZone = tpl.Postgres.TimeZone
	}
	// 表单不提交 node_id，零值与「显式指定 0 号节点」无法区分：
	// 向导场景下按模板值（-1，自动推导）补齐，需要固定节点号请直接编辑配置文件。
	if c.Snowflake.NodeID == 0 {
		c.Snowflake.NodeID = tpl.Snowflake.NodeID
	}
	if len(c.Seed.Accounts) == 0 {
		c.Seed.Accounts = tpl.Seed.Accounts
	}
}

// Validate 校验配置完整性，缺失即报错，供初始化向导保存前调用。
func (c *Config) Validate() error {
	var errs []string
	require := func(missing bool, key string) {
		if missing {
			errs = append(errs, key+" 未配置")
		}
	}

	require(c.App.HTTPPort <= 0 || c.App.HTTPPort > 65535, "app.http_port（1~65535）")
	switch c.App.RunMode {
	case "debug", "release", "test":
	default:
		errs = append(errs, "app.run_mode 必须是 debug / release / test")
	}
	require(strings.TrimSpace(c.App.Name) == "", "app.name")
	require(strings.TrimSpace(c.App.WebDir) == "", "app.web_dir")

	require(strings.TrimSpace(c.Postgres.Host) == "", "postgres.host")
	require(c.Postgres.Port <= 0, "postgres.port")
	require(strings.TrimSpace(c.Postgres.User) == "", "postgres.user")
	require(strings.TrimSpace(c.Postgres.DBName) == "", "postgres.dbname")
	require(strings.TrimSpace(c.Postgres.SSLMode) == "", "postgres.sslmode")
	require(strings.TrimSpace(c.Postgres.TimeZone) == "", "postgres.timezone")

	// Redis 对本服务是必需组件（会话 / 消息投递 / 在线状态都依赖它）
	require(strings.TrimSpace(c.Redis.Addr) == "", "redis.addr（host:port）")
	if addr := strings.TrimSpace(c.Redis.Addr); addr != "" {
		i := strings.LastIndex(addr, ":")
		if i <= 0 {
			errs = append(errs, "redis.addr 格式应为 host:port")
		} else if p, err := strconv.Atoi(strings.TrimSpace(addr[i+1:])); err != nil || p <= 0 {
			errs = append(errs, "redis.addr 端口无效")
		}
	}

	require(strings.TrimSpace(c.JWT.Secret) == "", "jwt.secret")
	require(c.JWT.Expire <= 0, "jwt.expire（秒）")

	if c.Snowflake.NodeID < -1 || c.Snowflake.NodeID > snowflake.NodeMax {
		errs = append(errs, "snowflake.node_id 必须是 -1~1023")
	}

	if c.Seed.Enabled {
		require(strings.TrimSpace(c.Seed.Password) == "", "seed.password")
		require(len(c.Seed.Accounts) == 0, "seed.accounts")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "；"))
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
