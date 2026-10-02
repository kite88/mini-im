// Package uuidv7 实现 RFC 9562 定义的 UUID v7（时间有序 UUID）。
//
// 128 位大端布局：
//
//	48 位毫秒时间戳 | 4 位版本号(0111) | 12 位序列号 | 2 位变体(10) | 62 位随机数
//
// 时间戳占据最高位，因此生成的 UUID 天然按时间递增，作为主键时新数据总是追加写入
// B 树右端，与雪花 ID 的索引友好性一致，但不需要协调节点号、128 位也不会用尽。
// 同一毫秒内用 12 位序列号递增（而非纯随机）来区分，保证同毫秒产生的多条记录
// 依旧有稳定、递增的顺序，避免纯随机 UUID 在索引页里来回插入。
package uuidv7

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

const (
	// seqMax 12 位序列号上限，即单进程每毫秒最多生成 4096 个 UUID
	seqMax = 1<<12 - 1

	// tsMask 48 位时间戳掩码，超出部分截断（可用到公元 10889 年）
	tsMask = 1<<48 - 1
)

var (
	mu     sync.Mutex
	lastMs int64
	seq    int64
)

// New 生成一个 UUID v7，返回标准 36 位小写字符串（8-4-4-4-12）。
//
// 同一毫秒内序列号耗尽时自旋等待下一毫秒；发生时钟回拨时沿用上次时间戳，
// 靠序列号继续递增，从而保证同一进程内永不重复、永不回退（与 pkg/snowflake 的策略一致）。
func New() string {
	var b [16]byte

	ms, s := next()
	ts := uint64(ms) & tsMask
	b[0] = byte(ts >> 40)
	b[1] = byte(ts >> 32)
	b[2] = byte(ts >> 24)
	b[3] = byte(ts >> 16)
	b[4] = byte(ts >> 8)
	b[5] = byte(ts)

	b[6] = 0x70 | byte(s>>8) // 版本号 0111 + 序列号高 4 位
	b[7] = byte(s)           // 序列号低 8 位

	var rnd [8]byte
	// 唯一性由「时间戳 + 序列号」保证，随机位只用于跨进程区分，
	// 因此 crypto/rand 极端失败时也没必要中断生成。
	_, _ = rand.Read(rnd[:])
	b[8] = 0x80 | rnd[0]&0x3f // 变体 10 + 随机数高 6 位
	copy(b[9:], rnd[1:])      // 随机数低 56 位

	return format(b)
}

// next 推进内部状态，返回本次使用的时间戳（毫秒）与序列号
func next() (int64, int64) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now().UnixMilli()
	if now < lastMs {
		now = lastMs
	}

	if now == lastMs {
		seq++
		if seq > seqMax {
			// 同一毫秒内序列号耗尽：自旋等到下一毫秒，保证 uuid 严格递增
			for now <= lastMs {
				now = time.Now().UnixMilli()
			}
			seq = 0
		}
	} else {
		seq = 0
	}
	lastMs = now

	return now, seq
}

// format 把 16 字节按 8-4-4-4-12 输出为小写十六进制，与 RFC 9562 的字符串表示一致
func format(b [16]byte) string {
	var out [36]byte
	hex.Encode(out[0:8], b[0:4])
	out[8] = '-'
	hex.Encode(out[9:13], b[4:6])
	out[13] = '-'
	hex.Encode(out[14:18], b[6:8])
	out[18] = '-'
	hex.Encode(out[19:23], b[8:10])
	out[23] = '-'
	hex.Encode(out[24:36], b[10:16])
	return string(out[:])
}
