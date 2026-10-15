package setup

import (
	_ "embed"
	"html"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"mini-im/conf"
)

//go:embed done.html
var doneHTML []byte

// DoneHandler 返回「初始化已完成」的服务状态页。
// 业务服务启动后访问 /setup 会看到该页面，而不是 404。
//
// webApp 表示前端目录里有可用页面（web_dir 指向的目录含 index.html）：
// 这种模式下页面与接口同源，首页就是本服务的根路径，因此多给一个首页入口。
func DoneHandler(cfg *conf.Config, webApp bool) gin.HandlerFunc {
	cfgPath := conf.FilePath
	if abs, err := filepath.Abs(cfgPath); err == nil {
		cfgPath = abs
	}

	return func(c *gin.Context) {
		// 首页地址取本次请求的 origin：用户可能是经网关 / 反代访问的，
		// 此时 host 与端口都和本机不同，写死 localhost 会把人带去错误的地址。
		homeURL := originOf(c)
		if homeURL == "" {
			homeURL = "http://localhost:" + strconv.Itoa(cfg.App.HTTPPort) + "/"
		} else {
			homeURL += "/"
		}

		// 页面里每处替换值都是固定的几项，逐请求替换的量可以忽略
		// （这页只在配置完成后被打开，频率极低），换来的是首页地址能跟着请求走。
		homeRow := ""
		if webApp {
			homeRow = `<div><span>首页地址</span><code class="copyable" id="homeUrl">` +
				html.EscapeString(homeURL) + `</code><a class="btn primary sm" id="btnHome" href="` +
				html.EscapeString(homeURL) + `">打开主页</a></div>`
		}

		// 从初始化向导跳过来时带上旧配置的备份路径，只在这次展示（刷新即消失）
		backupRow := ""
		if backup := c.Query("backup"); backup != "" {
			backupRow = `<div><span>原配置已备份</span><code class="copyable" title="点击复制">` +
				html.EscapeString(backup) + `</code></div>`
		}

		// 这一页有两种来路：从向导跳转过来（刚配好），或事后自己访问 /setup（只是查状态）。
		// 副标题按来路区分，否则长期访问时一直写着「初始化配置已完成」会很奇怪。
		subtitle := "服务已就绪。"
		if c.Query("fresh") == "1" {
			subtitle = "初始化配置已完成，服务已启动。"
		}

		page := strings.NewReplacer(
			"{{CONFIG_PATH}}", html.EscapeString(cfgPath),
			"{{PORT}}", strconv.Itoa(cfg.App.HTTPPort),
			"{{MODE}}", html.EscapeString(cfg.App.RunMode),
			"{{HOME_ROW}}", homeRow,
			"{{BACKUP_ROW}}", backupRow,
			"{{SUBTITLE}}", subtitle,
		).Replace(string(doneHTML))

		c.Header("Cache-Control", "no-store")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(page))
	}
}

// originOf 返回请求的 scheme://host[:port]，取不到时返回空串。
// 只信任 X-Forwarded-Proto（而不是 Host）：Host 会被客户端随意伪造，
// 而它仅用于本页展示的跳转链接，不参与任何鉴权或路由判断。
func originOf(c *gin.Context) string {
	if c == nil || c.Request == nil || c.Request.Host == "" {
		return ""
	}
	if proto := c.Request.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto + "://" + c.Request.Host
	}
	if c.Request.TLS != nil {
		return "https://" + c.Request.Host
	}
	return "http://" + c.Request.Host
}
