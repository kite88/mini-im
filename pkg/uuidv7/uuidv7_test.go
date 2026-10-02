package uuidv7

import (
	"encoding/hex"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestNewFormat 校验字符串形态：36 位、连字符位置、版本号 7、变体 10xx
func TestNewFormat(t *testing.T) {
	s := New()
	if len(s) != 36 {
		t.Fatalf("长度应为 36，实际 %d：%s", len(s), s)
	}
	for _, i := range []int{8, 13, 18, 23} {
		if s[i] != '-' {
			t.Fatalf("第 %d 位应为连字符：%s", i, s)
		}
	}
	if s[14] != '7' {
		t.Fatalf("版本号应为 7：%s", s)
	}
	if !strings.ContainsRune("89ab", rune(s[19])) {
		t.Fatalf("变体位应为 8/9/a/b：%s", s)
	}
	if _, err := hex.DecodeString(strings.ReplaceAll(s, "-", "")); err != nil {
		t.Fatalf("非法的十六进制：%v", err)
	}
}

// TestNewTimestamp 前 48 位应能还原出接近当前时间的时间戳
func TestNewTimestamp(t *testing.T) {
	s := New()

	raw, err := hex.DecodeString(s[0:8] + s[9:13]) // 前 6 字节即 48 位毫秒时间戳
	if err != nil {
		t.Fatal(err)
	}
	var ts int64
	for _, c := range raw {
		ts = ts<<8 | int64(c)
	}

	if diff := time.Now().UnixMilli() - ts; diff < 0 || diff > 5000 {
		t.Fatalf("时间戳偏差过大：%d ms（ts=%d）", diff, ts)
	}
}

// TestNewUniqueMonotonic 批量生成：全局唯一，且字符串序严格递增
func TestNewUniqueMonotonic(t *testing.T) {
	const n = 20000

	seen := make(map[string]struct{}, n)
	prev := ""
	for i := 0; i < n; i++ {
		s := New()
		if _, ok := seen[s]; ok {
			t.Fatalf("第 %d 个出现重复：%s", i, s)
		}
		seen[s] = struct{}{}
		if s <= prev {
			t.Fatalf("未保持递增：%s 出现在 %s 之后", s, prev)
		}
		prev = s
	}
}

// TestNewConcurrent 并发调用下同样不允许重复
func TestNewConcurrent(t *testing.T) {
	const workers, per = 8, 2000

	var mu sync.Mutex
	seen := make(map[string]struct{}, workers*per)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < per; j++ {
				s := New()
				mu.Lock()
				if _, ok := seen[s]; ok {
					mu.Unlock()
					t.Errorf("并发下出现重复：%s", s)
					return
				}
				seen[s] = struct{}{}
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if len(seen) != workers*per {
		t.Fatalf("生成数量不足：%d/%d", len(seen), workers*per)
	}
}
