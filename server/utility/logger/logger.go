// Package logger
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package logger

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"runtime"
	"strings"
	"time"
)

// Logger 日志工具结构
type Logger struct {
	logger *glog.Logger
	config *Config
}

// Config 日志配置
type Config struct {
	// 基本配置
	Level          string `json:"level"`          // 日志级别: debug, info, warn, error, fatal
	Format         string `json:"format"`         // 日志格式: text, json
	TimeFormat     string `json:"timeFormat"`     // 时间格式
	CallerSkip     int    `json:"callerSkip"`     // 调用栈跳过层数
	CallerLine     bool   `json:"callerLine"`     // 是否显示调用行号
	CallerFile     bool   `json:"callerFile"`     // 是否显示调用文件
	CallerFunction bool   `json:"callerFunction"` // 是否显示调用函数

	// 输出配置
	StdoutPrint bool   `json:"stdoutPrint"` // 是否输出到控制台
	StdoutColor bool   `json:"stdoutColor"` // 控制台是否彩色输出
	Path        string `json:"path"`        // 日志文件路径
	File        string `json:"file"`        // 日志文件名

	// 切割配置
	RotateExpire   time.Duration `json:"rotateExpire"`   // 切割过期时间
	RotateSize     int64         `json:"rotateSize"`     // 切割大小
	RotateCount    int           `json:"rotateCount"`    // 切割数量
	RotateBackup   bool          `json:"rotateBackup"`   // 是否备份
	RotateCompress bool          `json:"rotateCompress"` // 是否压缩

	// 性能配置
	Async         bool `json:"async"`         // 是否异步写入
	AsyncChanSize int  `json:"asyncChanSize"` // 异步通道大小

	// 上下文配置
	CtxKeys    []string `json:"ctxKeys"`    // 上下文键名
	HeaderKeys []string `json:"headerKeys"` // 请求头键名

	// 环境相关
	EnvName string `json:"envName"` // 环境名称
	AppName string `json:"appName"` // 应用名称
	Version string `json:"version"` // 版本号
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:          "info",
		Format:         "text",
		TimeFormat:     "2006-01-02 15:04:05.000",
		CallerSkip:     3,
		CallerLine:     true,
		CallerFile:     true,
		CallerFunction: false,
		StdoutPrint:    true,
		StdoutColor:    true,
		Path:           "logs",
		File:           "app.log",
		RotateExpire:   7 * 24 * time.Hour,
		RotateSize:     100 * 1024 * 1024, // 100MB
		RotateCount:    10,
		RotateBackup:   true,
		RotateCompress: true,
		Async:          true,
		AsyncChanSize:  1000,
		CtxKeys:        []string{"request_id", "trace_id", "user_id"},
		HeaderKeys:     []string{"X-Request-ID", "X-Trace-ID", "X-User-ID"},
		EnvName:        "develop",
		AppName:        "hotgo",
		Version:        "1.0.0",
	}
}

// NewLogger 创建新的日志实例
func NewLogger(config *Config) *Logger {
	if config == nil {
		config = DefaultConfig()
	}

	logger := glog.New()

	// 设置日志级别
	setLogLevel(logger, config.Level)

	// 设置输出配置
	setupOutput(logger, config)

	// 设置调用栈配置
	setupCaller(logger, config)

	// 设置切割配置
	setupRotate(logger, config)

	// 设置异步配置
	setupAsync(logger, config)

	return &Logger{
		logger: logger,
		config: config,
	}
}

// NewFromConfig 从配置文件创建日志实例
func NewFromConfig(ctx context.Context) *Logger {
	config := loadConfigFromFile(ctx)
	return NewLogger(config)
}

// setLogLevel 设置日志级别
func setLogLevel(logger *glog.Logger, level string) {
	switch strings.ToLower(level) {
	case "debug":
		logger.SetLevel(glog.LEVEL_DEBU)
	case "info":
		logger.SetLevel(glog.LEVEL_INFO)
	case "warn", "warning":
		logger.SetLevel(glog.LEVEL_WARN)
	case "error":
		logger.SetLevel(glog.LEVEL_ERRO)
	case "fatal":
		logger.SetLevel(glog.LEVEL_FATA)
	default:
		logger.SetLevel(glog.LEVEL_INFO)
	}
}

// setupOutput 设置输出配置
func setupOutput(logger *glog.Logger, config *Config) {
	// 控制台输出
	if config.StdoutPrint {
		logger.SetStdoutPrint(true)
		if !config.StdoutColor {
			logger.SetStdoutColorDisabled(true)
		}
	} else {
		logger.SetStdoutPrint(false)
	}

	// 文件输出
	if config.Path != "" {
		if !gfile.Exists(config.Path) {
			if err := gfile.Mkdir(config.Path); err != nil {
				panic(fmt.Sprintf("创建日志目录失败: %v", err))
			}
		}

		filePath := config.Path
		if config.File != "" {
			filePath = gfile.Join(config.Path, config.File)
		}
		logger.SetPath(filePath)
	}
}

// setupCaller 设置调用栈配置
func setupCaller(logger *glog.Logger, config *Config) {
	logger.SetLevelPrint(true)

	// 设置是否显示调用栈
	if config.CallerLine {
		logger.SetStack(true)
	}

	// 设置时间格式
	if config.TimeFormat != "" {
		logger.SetTimeFormat(config.TimeFormat)
	}
}

// setupRotate 设置切割配置
func setupRotate(logger *glog.Logger, config *Config) {
	// 注意：GoFrame v2的日志轮转需要通过配置文件或其他方式实现
	// 这里保留配置参数，但不直接调用不存在的方法
	// 实际的日志轮转需要通过GoFrame的配置管理器来实现

	// 可以通过设置文件路径来启用文件输出
	if config.Path != "" && config.File != "" {
		logPath := gfile.Join(config.Path, config.File)
		logger.SetFile(logPath)
	}
}

// setupAsync 设置异步配置
func setupAsync(logger *glog.Logger, config *Config) {
	if config.Async {
		logger.SetAsync(true)
		if config.AsyncChanSize > 0 {
			// GoFrame 2.0 可能需要其他方式设置异步通道大小
			// 这里保留配置，以备后续使用
		}
	}
}

// loadConfigFromFile 从配置文件加载配置
func loadConfigFromFile(ctx context.Context) *Config {
	config := DefaultConfig()

	// 从系统配置读取
	if mode := g.Cfg().MustGet(ctx, "system.mode").String(); mode != "" {
		config.EnvName = mode
	}

	if appName := g.Cfg().MustGet(ctx, "system.appName").String(); appName != "" {
		config.AppName = appName
	}

	if version := g.Cfg().MustGet(ctx, "system.appVersion").String(); version != "" {
		config.Version = version
	}

	// 根据环境设置不同的日志级别
	switch config.EnvName {
	case "develop":
		config.Level = "debug"
		config.StdoutPrint = true
		config.StdoutColor = true
		config.CallerLine = true
		config.CallerFile = true
		config.CallerFunction = true
	case "testing":
		config.Level = "debug"
		config.StdoutPrint = true
		config.StdoutColor = false
		config.CallerLine = true
		config.CallerFile = true
		config.CallerFunction = false
	case "staging":
		config.Level = "info"
		config.StdoutPrint = false
		config.StdoutColor = false
		config.CallerLine = false
		config.CallerFile = false
		config.CallerFunction = false
	case "product":
		config.Level = "warn"
		config.StdoutPrint = false
		config.StdoutColor = false
		config.CallerLine = false
		config.CallerFile = false
		config.CallerFunction = false
	}

	// 从日志配置读取
	if level := g.Cfg().MustGet(ctx, "logger.level").String(); level != "" {
		config.Level = level
	}
	if format := g.Cfg().MustGet(ctx, "logger.format").String(); format != "" {
		config.Format = format
	}
	if path := g.Cfg().MustGet(ctx, "logger.path").String(); path != "" {
		config.Path = path
	}
	if file := g.Cfg().MustGet(ctx, "logger.file").String(); file != "" {
		config.File = file
	}
	if stdoutPrint := g.Cfg().MustGet(ctx, "logger.stdoutPrint"); !stdoutPrint.IsEmpty() {
		config.StdoutPrint = stdoutPrint.Bool()
	}
	if async := g.Cfg().MustGet(ctx, "logger.async"); !async.IsEmpty() {
		config.Async = async.Bool()
	}
	if rotateSize := g.Cfg().MustGet(ctx, "logger.rotateSize").Int64(); rotateSize > 0 {
		config.RotateSize = rotateSize
	}
	if rotateExpire := g.Cfg().MustGet(ctx, "logger.rotateExpire").String(); rotateExpire != "" {
		if duration, err := time.ParseDuration(rotateExpire); err == nil {
			config.RotateExpire = duration
		}
	}
	if rotateCount := g.Cfg().MustGet(ctx, "logger.rotateCount").Int(); rotateCount > 0 {
		config.RotateCount = rotateCount
	}
	if rotateBackup := g.Cfg().MustGet(ctx, "logger.rotateBackup"); !rotateBackup.IsEmpty() {
		config.RotateBackup = rotateBackup.Bool()
	}
	if rotateCompress := g.Cfg().MustGet(ctx, "logger.rotateCompress"); !rotateCompress.IsEmpty() {
		config.RotateCompress = rotateCompress.Bool()
	}
	if callerSkip := g.Cfg().MustGet(ctx, "logger.callerSkip").Int(); callerSkip > 0 {
		config.CallerSkip = callerSkip
	}
	if callerLine := g.Cfg().MustGet(ctx, "logger.callerLine"); !callerLine.IsEmpty() {
		config.CallerLine = callerLine.Bool()
	}
	if callerFile := g.Cfg().MustGet(ctx, "logger.callerFile"); !callerFile.IsEmpty() {
		config.CallerFile = callerFile.Bool()
	}
	if callerFunction := g.Cfg().MustGet(ctx, "logger.callerFunction"); !callerFunction.IsEmpty() {
		config.CallerFunction = callerFunction.Bool()
	}
	if ctxKeys := g.Cfg().MustGet(ctx, "logger.ctxKeys"); !ctxKeys.IsEmpty() {
		config.CtxKeys = ctxKeys.Strings()
	}
	if headerKeys := g.Cfg().MustGet(ctx, "logger.headerKeys"); !headerKeys.IsEmpty() {
		config.HeaderKeys = headerKeys.Strings()
	}

	// 设置日志文件路径
	if config.Path != "" {
		envPath := gfile.Join(config.Path, config.EnvName)
		config.Path = envPath
	}

	return config
}

// Debug 记录调试日志
func (l *Logger) Debug(ctx context.Context, v ...interface{}) {
	l.logger.Debug(ctx, l.formatMessage(ctx, v...)...)
}

// Debugf 记录格式化调试日志
func (l *Logger) Debugf(ctx context.Context, format string, v ...interface{}) {
	l.logger.Debugf(ctx, l.formatFormat(ctx, format), v...)
}

// Info 记录信息日志
func (l *Logger) Info(ctx context.Context, v ...interface{}) {
	l.logger.Info(ctx, l.formatMessage(ctx, v...)...)
}

// Infof 记录格式化信息日志
func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.logger.Infof(ctx, l.formatFormat(ctx, format), v...)
}

// Warn 记录警告日志
func (l *Logger) Warn(ctx context.Context, v ...interface{}) {
	l.logger.Warning(ctx, l.formatMessage(ctx, v...)...)
}

// Warnf 记录格式化警告日志
func (l *Logger) Warnf(ctx context.Context, format string, v ...interface{}) {
	l.logger.Warningf(ctx, l.formatFormat(ctx, format), v...)
}

// Error 记录错误日志
func (l *Logger) Error(ctx context.Context, v ...interface{}) {
	l.logger.Error(ctx, l.formatMessage(ctx, v...)...)
}

// Errorf 记录格式化错误日志
func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.logger.Errorf(ctx, l.formatFormat(ctx, format), v...)
}

// Fatal 记录致命错误日志
func (l *Logger) Fatal(ctx context.Context, v ...interface{}) {
	l.logger.Fatal(ctx, l.formatMessage(ctx, v...)...)
}

// Fatalf 记录格式化致命错误日志
func (l *Logger) Fatalf(ctx context.Context, format string, v ...interface{}) {
	l.logger.Fatalf(ctx, l.formatFormat(ctx, format), v...)
}

// formatMessage 格式化消息
func (l *Logger) formatMessage(ctx context.Context, v ...interface{}) []interface{} {
	if l.config.Format == "json" {
		return []interface{}{l.buildStructuredLog(ctx, v...)}
	}
	return l.addContext(ctx, v...)
}

// formatFormat 格式化格式化字符串
func (l *Logger) formatFormat(ctx context.Context, format string) string {
	if l.config.Format == "json" {
		return "%s" // JSON格式将在buildStructuredLog中处理
	}
	return l.addContextToFormat(ctx, format)
}

// addContext 添加上下文信息
func (l *Logger) addContext(ctx context.Context, v ...interface{}) []interface{} {
	result := make([]interface{}, 0, len(v)+len(l.config.CtxKeys))

	// 添加上下文信息
	contextInfo := l.getContextInfo(ctx)
	if contextInfo != "" {
		result = append(result, contextInfo)
	}

	// 添加原始消息
	result = append(result, v...)

	return result
}

// addContextToFormat 添加上下文到格式化字符串
func (l *Logger) addContextToFormat(ctx context.Context, format string) string {
	contextInfo := l.getContextInfo(ctx)
	if contextInfo != "" {
		return contextInfo + " " + format
	}
	return format
}

// getContextInfo 获取上下文信息
func (l *Logger) getContextInfo(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	info := make([]string, 0, len(l.config.CtxKeys)+1)

	// 首先添加GoFrame框架的链路追踪ID
	if traceID := gctx.CtxId(ctx); traceID != "" {
		info = append(info, fmt.Sprintf("trace_id=%s", traceID))
	}

	// 添加配置的上下文键
	for _, key := range l.config.CtxKeys {
		if key == "trace_id" {
			continue // 跳过trace_id，已经在上面处理了
		}
		if value := ctx.Value(key); value != nil {
			info = append(info, fmt.Sprintf("%s=%v", key, value))
		}
	}

	if len(info) > 0 {
		return "[" + strings.Join(info, " ") + "]"
	}

	return ""
}

// buildStructuredLog 构建结构化日志
func (l *Logger) buildStructuredLog(ctx context.Context, v ...interface{}) map[string]interface{} {
	log := map[string]interface{}{
		"timestamp": gtime.Now().Format(l.config.TimeFormat),
		"level":     l.getCurrentLevel(),
		"env":       l.config.EnvName,
		"app":       l.config.AppName,
		"version":   l.config.Version,
		"message":   fmt.Sprint(v...),
	}

	// 添加上下文信息
	if ctx != nil {
		// 首先添加GoFrame框架的链路追踪ID
		if traceID := gctx.CtxId(ctx); traceID != "" {
			log["trace_id"] = traceID
		}

		// 添加配置的上下文键
		for _, key := range l.config.CtxKeys {
			if value := ctx.Value(key); value != nil {
				log[key] = value
			}
		}
	}

	// 添加调用栈信息
	if l.config.CallerFile || l.config.CallerFunction {
		if caller := l.getCaller(); caller != "" {
			log["caller"] = caller
		}
	}

	return log
}

// getCurrentLevel 获取当前日志级别
func (l *Logger) getCurrentLevel() string {
	level := l.logger.GetLevel()
	switch level {
	case glog.LEVEL_DEBU:
		return "debug"
	case glog.LEVEL_INFO:
		return "info"
	case glog.LEVEL_WARN:
		return "warn"
	case glog.LEVEL_ERRO:
		return "error"
	case glog.LEVEL_FATA:
		return "fatal"
	default:
		return "info"
	}
}

// getCaller 获取调用栈信息
func (l *Logger) getCaller() string {
	_, file, line, ok := runtime.Caller(l.config.CallerSkip)
	if !ok {
		return ""
	}

	if l.config.CallerFile {
		file = gfile.Basename(file)
	}

	if l.config.CallerFunction {
		pc, _, _, ok := runtime.Caller(l.config.CallerSkip)
		if ok {
			funcName := runtime.FuncForPC(pc).Name()
			if idx := strings.LastIndex(funcName, "."); idx != -1 {
				funcName = funcName[idx+1:]
			}
			return fmt.Sprintf("%s:%d@%s", file, line, funcName)
		}
	}

	return fmt.Sprintf("%s:%d", file, line)
}

// IsDebugEnabled 是否启用调试级别
func (l *Logger) IsDebugEnabled() bool {
	return l.logger.GetLevel() <= glog.LEVEL_DEBU
}

// IsInfoEnabled 是否启用信息级别
func (l *Logger) IsInfoEnabled() bool {
	return l.logger.GetLevel() <= glog.LEVEL_INFO
}

// IsWarnEnabled 是否启用警告级别
func (l *Logger) IsWarnEnabled() bool {
	return l.logger.GetLevel() <= glog.LEVEL_WARN
}

// IsErrorEnabled 是否启用错误级别
func (l *Logger) IsErrorEnabled() bool {
	return l.logger.GetLevel() <= glog.LEVEL_ERRO
}

// Close 关闭日志器
func (l *Logger) Close() error {
	// GoFrame的glog会自动处理关闭
	return nil
}

// GetConfig 获取配置
func (l *Logger) GetConfig() *Config {
	return l.config
}

// WithContext 创建带上下文的日志器
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		logger: l.logger,
		config: l.config,
	}
}

// 全局日志实例
var defaultLogger *Logger

// initDefaultLogger 初始化默认日志器
func initDefaultLogger() {
	if defaultLogger == nil {
		ctx := gctx.New()
		defaultLogger = NewFromConfig(ctx)
	}
}

// Debug 全局调试日志
func Debug(ctx context.Context, v ...interface{}) {
	initDefaultLogger()
	defaultLogger.Debug(ctx, v...)
}

// Debugf 全局格式化调试日志
func Debugf(ctx context.Context, format string, v ...interface{}) {
	initDefaultLogger()
	defaultLogger.Debugf(ctx, format, v...)
}

// Info 全局信息日志
func Info(ctx context.Context, v ...interface{}) {
	initDefaultLogger()
	defaultLogger.Info(ctx, v...)
}

// Infof 全局格式化信息日志
func Infof(ctx context.Context, format string, v ...interface{}) {
	initDefaultLogger()
	defaultLogger.Infof(ctx, format, v...)
}

// Warn 全局警告日志
func Warn(ctx context.Context, v ...interface{}) {
	initDefaultLogger()
	defaultLogger.Warn(ctx, v...)
}

// Warnf 全局格式化警告日志
func Warnf(ctx context.Context, format string, v ...interface{}) {
	initDefaultLogger()
	defaultLogger.Warnf(ctx, format, v...)
}

// Error 全局错误日志
func Error(ctx context.Context, v ...interface{}) {
	initDefaultLogger()
	defaultLogger.Error(ctx, v...)
}

// Errorf 全局格式化错误日志
func Errorf(ctx context.Context, format string, v ...interface{}) {
	initDefaultLogger()
	defaultLogger.Errorf(ctx, format, v...)
}

// Fatal 全局致命错误日志
func Fatal(ctx context.Context, v ...interface{}) {
	initDefaultLogger()
	defaultLogger.Fatal(ctx, v...)
}

// Fatalf 全局格式化致命错误日志
func Fatalf(ctx context.Context, format string, v ...interface{}) {
	initDefaultLogger()
	defaultLogger.Fatalf(ctx, format, v...)
}

// IsDebugEnabled 全局是否启用调试级别
func IsDebugEnabled() bool {
	initDefaultLogger()
	return defaultLogger.IsDebugEnabled()
}

// IsInfoEnabled 全局是否启用信息级别
func IsInfoEnabled() bool {
	initDefaultLogger()
	return defaultLogger.IsInfoEnabled()
}

// SetDefaultLogger 设置默认日志器
func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
}

// GetDefaultLogger 获取默认日志器
func GetDefaultLogger() *Logger {
	initDefaultLogger()
	return defaultLogger
}
