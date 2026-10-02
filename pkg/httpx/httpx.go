// Package httpx 放置 HTTP 层的通用小工具。
package httpx

import "strings"

const bearerPrefix = "bearer "

// Bearer 从 Authorization 头中提取 token，兼容 "Bearer xxx" 与裸 token 两种写法
func Bearer(header string) string {
	value := strings.TrimSpace(header)
	if len(value) >= len(bearerPrefix) && strings.EqualFold(value[:len(bearerPrefix)], bearerPrefix) {
		return strings.TrimSpace(value[len(bearerPrefix):])
	}
	return value
}
