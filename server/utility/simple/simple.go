// Package simple
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package simple

import (
	"client-app/internal/consts"
	"client-app/utility/encrypt"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"strings"
	"time"
)

// JWT Token payload 结构
type JWTPayload struct {
	UserId     int64  `json:"userId"`     // 用户ID
	TenantId   int64  `json:"tenantId"`   // 租户ID
	TenantCode string `json:"tenantCode"` // 租户编码
	Username   string `json:"username"`   // 用户名
	RoleId     int64  `json:"roleId"`     // 主要角色ID
	RoleKey    string `json:"roleKey"`    // 角色标识
	DeptId     int64  `json:"deptId"`     // 部门ID
	App        string `json:"app"`        // 应用标识
	Iat        int64  `json:"iat"`        // 签发时间
	Exp        int64  `json:"exp"`        // 过期时间
}

// RouterPrefix 获取应用路由前缀
func RouterPrefix(ctx context.Context, app string) string {
	return g.Cfg().MustGet(ctx, "router."+app+".prefix", "/"+app+"").String()
}

// FilterMaskDemo 过滤演示环境下的配置隐藏字段
func FilterMaskDemo(ctx context.Context, src g.Map) g.Map {
	if src == nil {
		return nil
	}

	if !IsDemo(ctx) {
		return src
	}

	for k := range src {
		if _, ok := consts.ConfigMaskDemoField[k]; ok {
			src[k] = consts.DemoTips
		}
	}
	return src
}

// DefaultErrorTplContent 获取默认的错误模板内容
func DefaultErrorTplContent(ctx context.Context) string {
	return gfile.GetContents(g.Cfg().MustGet(ctx, "viewer.paths").String() + "/error/default.html")
}

// DecryptText 解密文本
func DecryptText(text string) (string, error) {
	str, err := gbase64.Decode([]byte(text))
	if err != nil {
		return "", err
	}

	str, err = encrypt.AesECBDecrypt(str, consts.RequestEncryptKey)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

// CheckPassword 检查密码
func CheckPassword(input, salt, hash string) (err error) {
	// 解密密码

	password, err := DecryptText(input)
	if err != nil {
		fmt.Println("-----------", err)
		return err
	}

	fmt.Println("------------", gmd5.MustEncryptString(password+salt))
	fmt.Println("------------", password)
	fmt.Println("------------", salt)
	if hash != gmd5.MustEncryptString(password+salt) {

		err = gerror.New("用户密码不正确")
		return
	}
	return
}

// GenerateJWTToken 生成JWT Token
func GenerateJWTToken(payload *JWTPayload, secretKey string) (string, error) {
	// JWT Header
	header := map[string]interface{}{
		"typ": "JWT",
		"alg": "HS256",
	}

	// 设置签发时间和过期时间
	payload.Iat = time.Now().Unix()
	if payload.Exp == 0 {
		payload.Exp = time.Now().Add(24 * time.Hour).Unix() // 默认24小时过期
	}

	headerJson := gjson.New(header).MustToJson()
	payloadJson := gjson.New(payload).MustToJson()

	// Base64 编码
	headerB64 := gbase64.EncodeToString([]byte(headerJson))
	payloadB64 := gbase64.EncodeToString([]byte(payloadJson))

	// 去掉填充字符
	headerB64 = strings.TrimRight(headerB64, "=")
	payloadB64 = strings.TrimRight(payloadB64, "=")

	// 生成签名
	message := fmt.Sprintf("%s.%s", headerB64, payloadB64)
	signature := gmd5.MustEncryptString(message + secretKey)

	// 返回完整token
	return fmt.Sprintf("%s.%s.%s", headerB64, payloadB64, signature), nil
}

// ParseJWTToken 解析JWT Token
func ParseJWTToken(token, secretKey string) (*JWTPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, gerror.New("无效的JWT格式")
	}

	headerB64, payloadB64, signature := parts[0], parts[1], parts[2]

	// 验证签名
	message := fmt.Sprintf("%s.%s", headerB64, payloadB64)
	expectedSignature := gmd5.MustEncryptString(message + secretKey)
	if signature != expectedSignature {
		return nil, gerror.New("JWT签名验证失败")
	}

	// 解码payload
	// 添加填充字符
	for len(payloadB64)%4 != 0 {
		payloadB64 += "="
	}

	payloadBytes, err := gbase64.Decode([]byte(payloadB64))
	if err != nil {
		return nil, gerror.Newf("解码JWT payload失败: %v", err)
	}

	var payload JWTPayload
	if err := gjson.New(payloadBytes).Scan(&payload); err != nil {
		return nil, gerror.Newf("解析JWT payload失败: %v", err)
	}

	// 检查过期时间
	if payload.Exp > 0 && time.Now().Unix() > payload.Exp {
		return nil, gerror.New("JWT token已过期")
	}

	return &payload, nil
}

// ExtractTokenFromHeader 从请求头中提取token
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", gerror.New("Authorization header缺失")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", gerror.New("无效的Authorization格式，应为: Bearer <token>")
	}

	return parts[1], nil
}

// GetJWTSecretKey 获取JWT密钥
func GetJWTSecretKey(ctx context.Context) string {
	return g.Cfg().MustGet(ctx, "server.ApiKey", "default_secret_key").String()
}

// SafeGo 安全的调用协程，遇到错误时输出错误日志而不是抛出panic
func SafeGo(ctx context.Context, goroutineFunc func(ctx context.Context)) {
	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					glog.Error(ctx, v)
				} else {
					glog.Errorf(ctx, "%+v", exception)
				}
			}
		}()
		goroutineFunc(ctx)
	}()
}

func Logf(level int, ctx context.Context, format string, v ...interface{}) {
	switch level {
	case glog.LEVEL_DEBU:
		g.Log().Debugf(ctx, format, v...)
	case glog.LEVEL_INFO:
		g.Log().Infof(ctx, format, v...)
	case glog.LEVEL_NOTI:
		g.Log().Noticef(ctx, format, v...)
	case glog.LEVEL_WARN:
		g.Log().Warningf(ctx, format, v...)
	case glog.LEVEL_ERRO:
		g.Log().Errorf(ctx, format, v...)
	case glog.LEVEL_CRIT:
		g.Log().Criticalf(ctx, format, v...)
	case glog.LEVEL_PANI:
		g.Log().Panicf(ctx, format, v...)
	case glog.LEVEL_FATA:
		g.Log().Fatalf(ctx, format, v...)
	default:
		g.Log().Errorf(ctx, format, v...)
	}
}
