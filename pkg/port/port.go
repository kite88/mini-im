// Package port 把「监听端口失败」这类操作系统层面的错误翻译成人能看懂的下一步操作。
//
// 起因：Windows 上 Hyper-V / WSL / Docker 的 NAT 服务（winnat）会随机排除若干端口区间，
// 落在区间里的端口即使没有任何进程监听，bind 也一律失败，报的却是
// "bind: An attempt was made to access a socket in a way forbidden by its access permissions"
// —— 看着像权限问题，实际是端口被系统保留，不看这段提示根本想不到要换端口。
package port

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"syscall"
)

// 直接比数值而不是 syscall.EACCES / syscall.EADDRINUSE：
// 这两个常量在 Windows 与 Unix 上的取值并不一致（Windows 的 WSAEACCES 是 10013，
// 而 syscall.EACCES 是 POSIX 的 13），用常量比较会在某个平台上永远匹配不上。
const (
	errnoEACCES        syscall.Errno = 13    // POSIX：权限不足（Unix 上多为 1024 以下端口）
	errnoWSAEACCES     syscall.Errno = 10013 // Windows：bind 被拒绝（端口被系统保留）
	errnoEADDRINUSE    syscall.Errno = 98    // Linux：地址已在使用
	errnoWSAEADDRINUSE syscall.Errno = 10048 // Windows：地址已在使用
)

// ListenError 包装监听失败的错误，补上一句「下一步该干什么」。
// 与端口被拒 / 被占用无关的错误原样返回，不改变原有行为。
//
// portText 是出问题的端口号（用于文案），envVar 是可以换端口的环境变量名（没有就传空串）。
func ListenError(portText, envVar string, err error) error {
	var errno syscall.Errno
	if !errors.As(err, &errno) {
		return err
	}

	switch errno {
	case errnoEACCES, errnoWSAEACCES:
		return fmt.Errorf("%w\n%s", err, reservedHint(portText, envVar))
	case errnoEADDRINUSE, errnoWSAEADDRINUSE:
		return fmt.Errorf("%w\n%s", err, busyHint(portText, envVar))
	default:
		return err
	}
}

// reservedHint 端口被系统保留 / 权限不足时的提示
func reservedHint(portText, envVar string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "提示：端口 %s ", portText)
	if runtime.GOOS == "windows" {
		b.WriteString("可能被 Windows 保留了：Hyper-V / WSL / Docker 的 NAT 服务（winnat）会随机排除一些端口区间，")
		b.WriteString("落在区间里的端口即使没有进程监听也 bind 不了。\n")
		b.WriteString("  · 查看被排除的区间：netsh int ipv4 show excludedportrange protocol=tcp\n")
		writeEnvHint(&b, envVar)
		b.WriteString("  · 或让系统重新分配：管理员 PowerShell 执行 net stop winnat && net start winnat")
		return b.String()
	}
	b.WriteString("监听被拒绝：1024 以下的端口通常需要 root 权限。\n")
	fmt.Fprintf(&b, "  · 换端口启动：%s=8080 后重新运行", envVar)
	return b.String()
}

// busyHint 端口已被其他进程占用时的提示
func busyHint(portText, envVar string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "提示：端口 %s 已被其他进程占用。\n", portText)
	if runtime.GOOS == "windows" {
		fmt.Fprintf(&b, "  · 查是谁在用：netstat -ano | findstr :%s\n", portText)
		writeEnvHint(&b, envVar)
		return b.String()
	}
	fmt.Fprintf(&b, "  · 查是谁在用：lsof -i :%s\n", portText)
	fmt.Fprintf(&b, "  · 换端口启动：%s=8080 后重新运行", envVar)
	return b.String()
}

func writeEnvHint(b *strings.Builder, envVar string) {
	if envVar == "" {
		return
	}
	fmt.Fprintf(b, "  · 换端口启动：set %s=8080（PowerShell： $env:%s='8080'）\n", envVar, envVar)
}
