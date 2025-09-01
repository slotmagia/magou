// Package consts
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package consts

import (
	"github.com/gogf/gf/v2/text/gstr"
)

// 错误解释
const (
	ErrorORM         = "sql执行异常"
	ErrorNotData     = "数据不存在"
	ErrorRotaPointer = "指针转换异常"
)

// 鉴权相关错误码
const (
	// 认证错误
	ErrUnauthorized = "UNAUTHORIZED"   // 未授权访问
	ErrTokenMissing = "TOKEN_MISSING"  // Token缺失
	ErrTokenInvalid = "TOKEN_INVALID"  // Token无效
	ErrTokenExpired = "TOKEN_EXPIRED"  // Token已过期
	ErrUserNotFound = "USER_NOT_FOUND" // 用户不存在
	ErrUserDisabled = "USER_DISABLED"  // 用户已禁用
	ErrUserLocked   = "USER_LOCKED"    // 用户已锁定

	// 权限错误
	ErrForbidden        = "FORBIDDEN"         // 禁止访问
	ErrPermissionDenied = "PERMISSION_DENIED" // 权限不足
	ErrRoleMissing      = "ROLE_MISSING"      // 角色缺失
	ErrDataScopeLimit   = "DATA_SCOPE_LIMIT"  // 数据权限限制
)

// 鉴权错误信息映射
var AuthErrorMessages = map[string]string{
	ErrUnauthorized:     "未授权访问，请先登录",
	ErrTokenMissing:     "访问令牌不能为空",
	ErrTokenInvalid:     "访问令牌无效",
	ErrTokenExpired:     "访问令牌已过期，请重新登录",
	ErrUserNotFound:     "用户不存在或已被删除",
	ErrUserDisabled:     "用户已被禁用，无法访问系统",
	ErrUserLocked:       "用户已被锁定，请联系管理员",
	ErrForbidden:        "禁止访问该资源",
	ErrPermissionDenied: "权限不足，无法执行该操作",
	ErrRoleMissing:      "用户角色缺失，请联系管理员分配角色",
	ErrDataScopeLimit:   "数据权限受限，无法访问该数据",
}

// GetAuthErrorMessage 获取鉴权错误信息
func GetAuthErrorMessage(errCode string) string {
	if msg, ok := AuthErrorMessages[errCode]; ok {
		return msg
	}
	return "认证或权限验证失败"
}

// 需要隐藏真实错误的Wrap，开启访问日志后仍然会将真实错误记录
var concealErrorSlice = []string{ErrorORM, ErrorRotaPointer}

// ErrorMessage 错误显示信息，非debug模式有效
func ErrorMessage(err error) (message string) {
	if err == nil {
		return "操作失败！~"
	}

	message = err.Error()
	for _, e := range concealErrorSlice {
		if gstr.Contains(message, e) {
			return "操作失败，请稍后重试！~"
		}
	}
	return
}
