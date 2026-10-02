package handler

import (
	"fmt"
	"hash/fnv"
	"html"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var avatarPalette = [][2]string{
	{"6C8DFF", "4A6BF5"},
	{"FF8A65", "F4511E"},
	{"4DB6AC", "00897B"},
	{"BA68C8", "8E24AA"},
	{"FFD54F", "FB8C00"},
	{"7986CB", "3949AB"},
	{"4DD0E1", "00ACC1"},
	{"AED581", "7CB342"},
}

const avatarTemplate = `<svg xmlns="http://www.w3.org/2000/svg" width="120" height="120" viewBox="0 0 120 120">` +
	`<defs><linearGradient id="g" x1="0" y1="0" x2="1" y2="1">` +
	`<stop offset="0%%" stop-color="#%s"/><stop offset="100%%" stop-color="#%s"/>` +
	`</linearGradient></defs>` +
	`<rect width="120" height="120" rx="28" fill="url(#g)"/>` +
	`<text x="60" y="60" text-anchor="middle" dominant-baseline="central" ` +
	`font-family="Microsoft YaHei,PingFang SC,Arial,sans-serif" font-size="52" font-weight="600" fill="#ffffff">%s</text>` +
	`</svg>`

// Avatar GET /avatar?name=xxx&label=昵称
//
// 动态生成 SVG 头像：同一 name 颜色固定，文字取 label 首字符。
// 这样无需任何图片资源即可获得统一的默认头像。
func Avatar(c *gin.Context) {
	name := strings.TrimSpace(c.Query("name"))
	label := strings.TrimSpace(c.Query("label"))
	if label == "" {
		label = name
	}

	seed := name
	if seed == "" {
		seed = label
	}

	pair := avatarPalette[hashIndex(seed)]
	text := "?"
	if runes := []rune(label); len(runes) > 0 {
		text = strings.ToUpper(string(runes[0]))
	}

	c.Header("Content-Type", "image/svg+xml; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=86400")
	c.Status(http.StatusOK)
	fmt.Fprintf(c.Writer, avatarTemplate, pair[0], pair[1], html.EscapeString(text))
}

func hashIndex(seed string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(seed))
	return h.Sum32() % uint32(len(avatarPalette))
}
