// Package ws 实现基于 gorilla/websocket 的消息中心。
package ws

import (
	"sync"

	"mini-im/service"
)

// Hub 维护所有在线连接，支持同一用户多端同时在线
type Hub struct {
	mu      sync.RWMutex
	clients map[int64]map[*Client]struct{}
}

// NewHub 创建消息中心
func NewHub() *Hub {
	return &Hub{clients: make(map[int64]map[*Client]struct{})}
}

// add 注册连接；某用户首个连接建立时标记为在线
func (h *Hub) add(c *Client) {
	h.mu.Lock()
	set, ok := h.clients[c.uid]
	if !ok {
		set = make(map[*Client]struct{})
		h.clients[c.uid] = set
	}
	set[c] = struct{}{}
	first := len(set) == 1
	h.mu.Unlock()

	if first {
		service.SetOnline(c.uid)
	}
}

// remove 注销连接；某用户最后一个连接断开时标记为离线
func (h *Hub) remove(c *Client) {
	h.mu.Lock()
	set, ok := h.clients[c.uid]
	if !ok {
		h.mu.Unlock()
		return
	}
	if _, exist := set[c]; !exist {
		h.mu.Unlock()
		return
	}
	delete(set, c)
	last := len(set) == 0
	if last {
		delete(h.clients, c.uid)
	}
	h.mu.Unlock()

	if last {
		service.SetOffline(c.uid)
	}
}

// SendTo 向指定用户的所有连接推送消息，返回是否至少有一条连接投递成功。
//
// 该方法同时满足 service.Notifier 接口，业务层（如好友申请）借此主动推送事件，
// 而无需 ws 与 service 相互引用。
func (h *Hub) SendTo(uid int64, payload []byte) bool {
	h.mu.RLock()
	set := h.clients[uid]
	targets := make([]*Client, 0, len(set))
	for c := range set {
		targets = append(targets, c)
	}
	h.mu.RUnlock()

	delivered := false
	for _, c := range targets {
		if c.push(payload) {
			delivered = true
		}
	}
	return delivered
}

// ConnCount 当前连接总数
func (h *Hub) ConnCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	total := 0
	for _, set := range h.clients {
		total += len(set)
	}
	return total
}
