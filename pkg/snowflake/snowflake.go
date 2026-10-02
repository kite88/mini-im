// Package snowflake 实现 Twitter Snowflake 算法的 64 位分布式 ID 生成器。
//
// 位布局：1 位符号位（恒 0）+ 41 位毫秒时间戳 + 10 位节点号 + 12 位序列号。
// 单节点每毫秒最多生成 4096 个 ID，41 位时间戳从起始时间起可用约 69 年，
// 生成的 ID 整体趋势递增，适合直接作为数据库主键（B 树索引友好）。
package snowflake

import (
	"hash/fnv"
	"net"
	"os"
	"sync"
	"time"
)

const (
	// epoch 起始时间：2024-01-01 00:00:00 UTC（毫秒）
	epoch int64 = 1704067200000

	// NodeBits 节点号位数；NodeMax 为节点号上限（含）
	NodeBits uint8 = 10
	NodeMax  int64 = (1 << NodeBits) - 1

	stepBits  uint8 = 12
	stepMax   int64 = (1 << stepBits) - 1
	timeShift       = NodeBits + stepBits
	nodeShift       = stepBits
)

// Generator 雪花 ID 生成器，可并发安全使用
type Generator struct {
	mu     sync.Mutex
	nodeID int64
	step   int64
	lastMs int64
}

// New 创建生成器并返回，nodeID 取值 0 ~ NodeMax；
// 传入负数或越界值时按本机地址自动推导节点号。
func New(nodeID int64) *Generator {
	if nodeID < 0 || nodeID > NodeMax {
		nodeID = localNodeID()
	}
	return &Generator{nodeID: nodeID}
}

// NodeID 当前生成器使用的节点号
func (g *Generator) NodeID() int64 { return g.nodeID }

// Next 生成下一个 ID。
//
// 同一毫秒内序列号耗尽时自旋等待下一毫秒；发生时钟回拨时沿用上次时间戳，
// 靠序列号继续递增，从而保证同一进程内 ID 永不重复、永不回退。
func (g *Generator) Next() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()
	if now < g.lastMs {
		now = g.lastMs
	}

	if now == g.lastMs {
		g.step = (g.step + 1) & stepMax
		if g.step == 0 {
			for now <= g.lastMs {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		g.step = 0
	}
	g.lastMs = now

	return (now-epoch)<<timeShift | g.nodeID<<nodeShift | g.step
}

// ---------------------------------------------------------------
// 全局生成器：绝大多数场景直接用 NextID 即可
// ---------------------------------------------------------------

var (
	globalMu sync.RWMutex
	global   *Generator
)

// Init 初始化全局生成器并返回实际生效的节点号，应在创建任何实体前调用一次。
func Init(nodeID int64) int64 {
	g := New(nodeID)

	globalMu.Lock()
	global = g
	globalMu.Unlock()

	return g.nodeID
}

// NextID 用全局生成器生成一个 ID；未显式 Init 时按本机地址惰性初始化，
// 保证「忘记初始化」不会导致空指针，也不会影响 ID 唯一性。
func NextID() int64 {
	globalMu.RLock()
	g := global
	globalMu.RUnlock()

	if g == nil {
		globalMu.Lock()
		if global == nil {
			global = New(-1)
		}
		g = global
		globalMu.Unlock()
	}

	return g.Next()
}

// localNodeID 取本机非回环 IPv4 地址的后两段（低 10 位）作为节点号，
// 同一网段内的多实例天然错开；取不到网卡地址时退化为按主机名哈希。
func localNodeID() int64 {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.IsLoopback() {
				continue
			}
			if ip4 := ipNet.IP.To4(); ip4 != nil {
				return (int64(ip4[2])<<8 | int64(ip4[3])) & NodeMax
			}
		}
	}

	host, err := os.Hostname()
	if err != nil || host == "" {
		return 1
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(host))
	return int64(h.Sum32()) & NodeMax
}
