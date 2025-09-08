package logger

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// sqlTracingEnabled 控制是否启用 SQL 追踪输出
var sqlTracingEnabled = true

// EnableSQLTracing 启用 SQL 追踪输出
func EnableSQLTracing() {
	sqlTracingEnabled = true
}

// DisableSQLTracing 禁用 SQL 追踪输出
func DisableSQLTracing() {
	sqlTracingEnabled = false
}

// TraceSQL 是 LogSQLWithCaller 的别名，便于直观调用
func TraceSQL(ctx context.Context, sql string, args []interface{}, duration time.Duration) {
	LogSQLWithCaller(ctx, sql, args, duration)
}

// LogSQLWithCaller 输出带调用方(文件:行:函数)的 SQL 日志，支持 IDE 控制台点击跳转
// 注意：为保证点击跳转，尽量输出绝对路径或相对工程根目录的路径。
func LogSQLWithCaller(ctx context.Context, sql string, args []interface{}, duration time.Duration) {
	if !sqlTracingEnabled {
		return
	}

	caller := findCallerForSQL()
	// 着色阈值：<10ms 绿色，10-100ms 黄色，>100ms 红色
	color := ""
	reset := ""
	ms := duration.Milliseconds()
	if IsDebugEnabled() {
		// 仅在控制台彩色输出启用时加颜色（简单判断，避免文件日志混乱）
		if GetDefaultLogger().config.StdoutPrint && GetDefaultLogger().config.StdoutColor {
			reset = "\x1b[0m"
			switch {
			case ms < 10:
				color = "\x1b[32m" // 绿色
			case ms < 100:
				color = "\x1b[33m" // 黄色
			default:
				color = "\x1b[31m" // 红色
			}
		}
	}

	// 第一行：时间/耗时/调用方
	Info(ctx, fmt.Sprintf("[%dms] %s%s%s [%s]",
		ms, color, duration.String(), reset, caller,
	))

	// 第二行：SQL 语句
	Debug(ctx, fmt.Sprintf("SQL: %s", sql))

	// 第三行：参数
	if len(args) > 0 {
		Debug(ctx, fmt.Sprintf("ARGS: %v", args))
	}
}

// findCallerForSQL 通过遍历调用栈，找到业务侧发起 SQL 的调用点
func findCallerForSQL() string {
	// 跳过前若干层为本日志封装/GoFrame/数据库层
	const maxDepth = 32
	for i := 2; i < maxDepth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		name := ""
		if fn != nil {
			name = fn.Name()
		}

		// 过滤掉非业务包的帧
		if shouldSkipFrame(file, name) {
			continue
		}

		// 使用工程内相对路径以便 IDE 识别并可点击
		clean := filepath.ToSlash(file)
		return fmt.Sprintf("%s:%d:%s", clean, line, shortFunc(name))
	}

	// 兜底：返回最外层一帧（可能不准确）
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown:0:unknown"
	}
	return fmt.Sprintf("%s:%d:%s", filepath.ToSlash(file), line, "unknown")
}

func shouldSkipFrame(file, funcName string) bool {
	path := filepath.ToSlash(file)
	// 跳过 Go 标准库、GoFrame、自身日志封装
	if strings.Contains(path, "/go/pkg/mod/") ||
		strings.Contains(path, "/runtime/") ||
		strings.Contains(path, "/github.com/gogf/gf/") ||
		strings.Contains(path, "/utility/logger/") {
		return true
	}
	// 可按需继续扩展要跳过的内部封装
	return false
}

func shortFunc(full string) string {
	if full == "" {
		return ""
	}
	if idx := strings.LastIndex(full, "/"); idx >= 0 {
		full = full[idx+1:]
	}
	// 去包名，仅保留函数
	if idx := strings.LastIndex(full, "."); idx >= 0 {
		return full[idx+1:]
	}
	return full
}
