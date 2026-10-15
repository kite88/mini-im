package conf

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// RenderConfig 以模板（config.example.yaml）为基底渲染配置文件：
// 只替换各配置项的值，模板中的注释、空行与排版原样保留。
//
// 之所以不直接 yaml.Marshal：那样会把人工维护的中文注释全部丢掉，
// 用户拿到的 config.yaml 就成了没有说明的键值堆。
func RenderConfig(template []byte, cfg *Config) ([]byte, error) {
	values := cfg.scalarValues()

	text := strings.ReplaceAll(string(template), "\r\n", "\n")
	lines := strings.Split(text, "\n")

	section := ""
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " "))
		key, rest, ok := splitKey(line[indent:])
		if !ok {
			continue
		}
		if indent == 0 {
			section = key
			continue
		}

		val, ok := values[section+"."+key]
		if !ok {
			continue
		}
		lines[i] = replaceValue(line, indent, key, rest, val)
	}

	return []byte(strings.Join(lines, "\n")), nil
}

// Save 渲染并写入配置文件；目标已存在时先备份为 <path>.bak。
func Save(path string, template []byte, cfg *Config) (backup string, err error) {
	data, err := RenderConfig(template, cfg)
	if err != nil {
		return "", err
	}

	if _, statErr := os.Stat(path); statErr == nil {
		backup = path + ".bak"
		if err := copyFile(path, backup); err != nil {
			return "", fmt.Errorf("备份原配置到 %s 失败：%w", backup, err)
		}
	}

	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return backup, fmt.Errorf("创建配置目录 %s 失败：%w", dir, err)
		}
	}
	// 配置里含数据库密码与 JWT 密钥，文件权限收紧到仅本人可读写
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return backup, fmt.Errorf("写入配置文件 %s 失败：%w", path, err)
	}
	return backup, nil
}

// scalarValues 把「向导表单覆盖得到的」配置项映射为 "段.键" -> YAML 文本。
// 字符串统一加双引号并正确转义；不在表里的键（app.name、app.web_dir、
// postgres.timezone、snowflake.node_id、seed.accounts）保持模板原值。
func (c *Config) scalarValues() map[string]string {
	str := quoteYAMLString
	return map[string]string{
		"app.http_port": strconv.Itoa(c.App.HTTPPort),
		"app.run_mode":  str(c.App.RunMode),

		"postgres.host":     str(c.Postgres.Host),
		"postgres.port":     strconv.Itoa(c.Postgres.Port),
		"postgres.user":     str(c.Postgres.User),
		"postgres.password": str(c.Postgres.Password),
		"postgres.dbname":   str(c.Postgres.DBName),
		"postgres.sslmode":  str(c.Postgres.SSLMode),

		"redis.addr":     str(c.Redis.Addr),
		"redis.password": str(c.Redis.Password),
		"redis.db":       strconv.Itoa(c.Redis.DB),

		"jwt.secret": str(c.JWT.Secret),
		"jwt.expire": strconv.FormatInt(c.JWT.Expire, 10),

		"seed.enabled":  strconv.FormatBool(c.Seed.Enabled),
		"seed.password": str(c.Seed.Password),
	}
}

// splitKey 把 "port: 8280  # 注释" 拆成 key 与剩余部分（含冒号后的空格）。
func splitKey(s string) (key, rest string, ok bool) {
	i := strings.IndexByte(s, ':')
	if i <= 0 {
		return "", "", false
	}
	return strings.TrimSpace(s[:i]), s[i+1:], true
}

// replaceValue 在原行内替换配置值，保持原缩进与行尾注释的对齐。
func replaceValue(line string, indent int, key, rest, newVal string) string {
	prefix := strings.Repeat(" ", indent) + key + ":"

	gap := rest[:len(rest)-len(strings.TrimLeft(rest, " \t"))]
	body := rest[len(gap):]

	comment := ""
	col := commentIndex(body)
	if col >= 0 {
		comment = body[col:]
		body = body[:col]
	}

	out := prefix + gap + newVal
	if comment == "" {
		return out
	}

	commentCol := displayWidth(prefix) + commentIndexWidth(rest)
	padding := commentCol - displayWidth(out)
	if padding < 2 {
		padding = 2
	}
	return out + strings.Repeat(" ", padding) + comment
}

// displayWidth 估算终端显示宽度（CJK 等全角字符按 2 列计算），用于对齐行尾注释。
func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		if isFullWidth(r) {
			w += 2
			continue
		}
		w++
	}
	return w
}

func isFullWidth(r rune) bool {
	switch {
	case r >= 0x1100 && r <= 0x115f,
		r == 0x2329, r == 0x232a,
		r >= 0x2e80 && r <= 0xa4cf && r != 0x303f,
		r >= 0xac00 && r <= 0xd7a3,
		r >= 0xf900 && r <= 0xfaff,
		r >= 0xfe30 && r <= 0xfe6f,
		r >= 0xff00 && r <= 0xff60,
		r >= 0xffe0 && r <= 0xffe6:
		return true
	}
	return false
}

// commentIndex 返回行内注释 "#" 的字节下标（忽略引号内的 #），无注释返回 -1。
func commentIndex(s string) int {
	inSingle, inDouble := false, false
	for i, r := range s {
		switch {
		case r == '\'' && !inDouble:
			inSingle = !inSingle
		case r == '"' && !inSingle:
			inDouble = !inDouble
		case r == '#' && !inSingle && !inDouble:
			if i == 0 || s[i-1] == ' ' || s[i-1] == '\t' {
				return i
			}
		}
	}
	return -1
}

// commentIndexWidth 返回 "#" 之前的显示宽度（无注释返回 0）。
func commentIndexWidth(s string) int {
	i := commentIndex(s)
	if i < 0 {
		return 0
	}
	return displayWidth(s[:i])
}

// quoteYAMLString 借助 yaml 编码器生成安全的双引号字符串（处理引号 / 冒号 / # 等特殊字符）。
func quoteYAMLString(v string) string {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	_ = enc.Encode(&yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: v,
		Style: yaml.DoubleQuotedStyle,
	})
	_ = enc.Close()
	return strings.TrimRight(buf.String(), "\n")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
