package middleware

import (
	"client-app/internal/library/contexts"
	"client-app/utility/logger"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"strings"
	"time"
)

// LogRequest 请求日志中间件
func (s *sMiddleware) LogRequest(r *ghttp.Request) {
	// 跳过静态资源和健康检查
	if s.shouldSkipLogging(r) {
		r.Middleware.Next()
		return
	}

	// 生成请求ID并设置响应头
	requestID := s.generateRequestID(r)

	// 在现有上下文基础上添加日志相关信息
	ctx := s.enhanceContextForLogging(r, requestID)
	r.SetCtx(ctx)

	// 记录请求开始
	startTime := gtime.Now()
	s.logRequestStart(ctx, r)

	// 执行请求
	r.Middleware.Next()

	// 记录请求结束
	s.logRequestEnd(ctx, r, startTime)
}

// shouldSkipLogging 是否跳过日志记录
func (s *sMiddleware) shouldSkipLogging(r *ghttp.Request) bool {
	// 跳过静态文件请求
	if r.IsFileRequest() {
		return true
	}

	// 跳过特定的HTTP方法
	ignoredMethods := []string{"HEAD", "PRI", "OPTIONS"}
	if gstr.InArray(ignoredMethods, strings.ToUpper(r.Method)) {
		return true
	}

	// 跳过健康检查等路径
	skipPaths := []string{
		"/health",
		"/ping",
		"/favicon.ico",
		"/robots.txt",
	}
	for _, path := range skipPaths {
		if strings.HasSuffix(r.URL.Path, path) {
			return true
		}
	}

	return false
}

// generateRequestID 生成请求ID
func (s *sMiddleware) generateRequestID(r *ghttp.Request) string {
	// 优先使用请求头中的ID
	requestID := r.Header.Get("X-Request-ID")

	// 如果没有，则生成新的
	if requestID == "" {
		requestID = guid.S()
	}

	// 设置响应头
	r.Response.Header().Set("X-Request-ID", requestID)

	return requestID
}

// enhanceContextForLogging 在现有上下文基础上增强日志信息
func (s *sMiddleware) enhanceContextForLogging(r *ghttp.Request, requestID string) context.Context {
	ctx := r.Context()

	// 获取GoFrame框架的链路追踪ID作为traceID（与现有系统保持一致）
	traceID := gctx.CtxId(ctx)

	// 在现有上下文基础上添加日志相关信息
	logContext := context.WithValue(ctx, "request_id", requestID)
	logContext = context.WithValue(logContext, "trace_id", traceID)
	logContext = context.WithValue(logContext, "method", r.Method)
	logContext = context.WithValue(logContext, "path", r.URL.Path)
	logContext = context.WithValue(logContext, "client_ip", r.GetClientIp())

	// 添加查询参数（如果存在）
	if r.URL.RawQuery != "" {
		logContext = context.WithValue(logContext, "query", r.URL.RawQuery)
	}

	// 添加常用请求头信息
	if userAgent := r.Header.Get("User-Agent"); userAgent != "" {
		logContext = context.WithValue(logContext, "user_agent", userAgent)
	}

	if contentType := r.Header.Get("Content-Type"); contentType != "" {
		logContext = context.WithValue(logContext, "content_type", contentType)
	}

	// 添加代理相关请求头
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		logContext = context.WithValue(logContext, "x_forwarded_for", forwardedFor)
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		logContext = context.WithValue(logContext, "x_real_ip", realIP)
	}

	return logContext
}

// logRequestStart 记录请求开始
func (s *sMiddleware) logRequestStart(ctx context.Context, r *ghttp.Request) {
	// 只在调试模式下记录请求开始
	if !logger.IsDebugEnabled() {
		return
	}

	// 记录请求体（敏感信息需要过滤）
	var requestBody string
	if r.Method != "GET" {
		body := r.GetBodyString()
		if body != "" {
			requestBody = s.filterSensitiveData(body)
		}
	}

	// 构建请求信息
	message := fmt.Sprintf("HTTP请求开始 %s %s", r.Method, r.URL.Path)
	if requestBody != "" {
		message += fmt.Sprintf(" Body: %s", requestBody)
	}

	logger.Debug(ctx, message)
}

// logRequestEnd 记录请求结束
func (s *sMiddleware) logRequestEnd(ctx context.Context, r *ghttp.Request, startTime *gtime.Time) {
	// 计算请求耗时
	duration := gtime.Now().Sub(startTime)

	// 获取响应信息
	status := r.Response.Status
	responseSize := r.Response.BufferLength()

	// 创建带响应信息的上下文
	endCtx := context.WithValue(ctx, "status", status)
	endCtx = context.WithValue(endCtx, "duration_ms", duration.Milliseconds())
	endCtx = context.WithValue(endCtx, "duration_s", duration.Seconds())
	endCtx = context.WithValue(endCtx, "response_size", responseSize)

	// 根据状态码和耗时选择日志级别
	level := s.getLogLevel(status, duration)

	// 构建日志消息
	message := fmt.Sprintf("HTTP请求完成 %s %s %d %dms %db",
		r.Method, r.URL.Path, status, duration.Milliseconds(), responseSize)

	// 记录日志
	switch level {
	case "debug":
		logger.Debug(endCtx, message)
	case "info":
		logger.Info(endCtx, message)
	case "warn":
		logger.Warn(endCtx, message)
	case "error":
		logger.Error(endCtx, message)
	}

	// 记录慢请求
	s.logSlowRequest(endCtx, r, duration)

	// 记录错误响应
	s.logErrorResponse(endCtx, r, status)
}

// getLogLevel 根据状态码和耗时确定日志级别
func (s *sMiddleware) getLogLevel(status int, duration time.Duration) string {
	// 5xx错误
	if status >= 500 {
		return "error"
	}

	// 4xx错误
	if status >= 400 {
		return "warn"
	}

	// 慢请求
	if duration > 5*time.Second {
		return "warn"
	}

	// 正常请求
	if duration > 1*time.Second {
		return "info"
	}

	return "debug"
}

// logSlowRequest 记录慢请求
func (s *sMiddleware) logSlowRequest(ctx context.Context, r *ghttp.Request, duration time.Duration) {
	slowThreshold := 2 * time.Second
	if duration > slowThreshold {
		message := fmt.Sprintf("慢请求告警 %s %s 耗时: %s",
			r.Method, r.URL.Path, duration.String())
		logger.Warn(ctx, message)
	}
}

// logErrorResponse 记录错误响应
func (s *sMiddleware) logErrorResponse(ctx context.Context, r *ghttp.Request, status int) {
	if status >= 400 {
		// 获取错误响应体
		var errorBody string
		if contexts.Get(ctx) != nil && contexts.Get(ctx).Response != nil {
			if data := contexts.Get(ctx).Response.Data; data != nil {
				// 使用类型断言处理interface{}类型
				switch v := data.(type) {
				case string:
					errorBody = v
				case []byte:
					errorBody = string(v)
				default:
					errorBody = gconv.String(data)
				}
			}
		}

		message := fmt.Sprintf("错误响应 %s %s 状态码: %d",
			r.Method, r.URL.Path, status)

		if errorBody != "" {
			message += fmt.Sprintf(" 响应: %s", s.filterSensitiveData(errorBody))
		}

		if status >= 500 {
			logger.Error(ctx, message)
		} else {
			logger.Warn(ctx, message)
		}
	}
}

// filterSensitiveData 过滤敏感数据
func (s *sMiddleware) filterSensitiveData(data string) string {
	// 敏感字段列表
	sensitiveFields := []string{
		"password", "pwd", "passwd", "secret", "token", "key",
		"authorization", "auth", "credential", "private", "confidential",
		"card", "credit", "debit", "account", "payment", "billing",
	}

	result := data
	for _, field := range sensitiveFields {
		// 简单的正则替换，实际项目中可能需要更复杂的处理
		patterns := []string{
			fmt.Sprintf(`"%s":\s*"[^"]*"`, field),
			fmt.Sprintf(`"%s":\s*[^,}\s]*`, field),
			fmt.Sprintf(`%s=[\w\d]*`, field),
		}

		for _, pattern := range patterns {
			result = strings.ReplaceAll(result, pattern, fmt.Sprintf(`"%s":"***"`, field))
		}
	}

	// 限制长度
	if len(result) > 1000 {
		result = result[:1000] + "..."
	}

	return result
}

// LogError 记录错误日志（用于业务逻辑层）
func (s *sMiddleware) LogError(ctx context.Context, err error, message string, args ...interface{}) {
	if err == nil {
		return
	}

	// 构建错误信息
	errorMsg := fmt.Sprintf(message, args...)
	if errorMsg != "" {
		errorMsg += ": "
	}
	errorMsg += err.Error()

	logger.Error(ctx, errorMsg)
}

// LogBusiness 记录业务日志
func (s *sMiddleware) LogBusiness(ctx context.Context, action string, data interface{}, result string) {
	message := fmt.Sprintf("业务操作 %s 结果: %s", action, result)
	if data != nil {
		message += fmt.Sprintf(" 数据: %v", data)
	}

	logger.Info(ctx, message)
}

// LogSecurity 记录安全日志
func (s *sMiddleware) LogSecurity(ctx context.Context, event string, level string, details ...interface{}) {
	message := fmt.Sprintf("安全事件 %s", event)
	if len(details) > 0 {
		message += fmt.Sprintf(" 详情: %v", details)
	}

	switch level {
	case "critical":
		logger.Fatal(ctx, message)
	case "high":
		logger.Error(ctx, message)
	case "medium":
		logger.Warn(ctx, message)
	default:
		logger.Info(ctx, message)
	}
}

// LogMetrics 记录性能指标
func (s *sMiddleware) LogMetrics(ctx context.Context, metric string, value interface{}, tags ...string) {
	message := fmt.Sprintf("性能指标 %s: %v", metric, value)
	if len(tags) > 0 {
		message += fmt.Sprintf(" 标签: %s", strings.Join(tags, ", "))
	}

	logger.Info(ctx, message)
}

// LogAudit 记录审计日志
func (s *sMiddleware) LogAudit(ctx context.Context, action string, resource string, result string, details ...interface{}) {
	message := fmt.Sprintf("审计日志 操作: %s 资源: %s 结果: %s", action, resource, result)
	if len(details) > 0 {
		message += fmt.Sprintf(" 详情: %v", details)
	}

	logger.Info(ctx, message)
}
