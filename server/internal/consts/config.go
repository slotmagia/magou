package consts

import (
	"github.com/gogf/gf/v2/frame/g"
)

// RequestEncryptKey
// 请求加密密钥用于敏感数据加密，16位字符，前后端需保持一致
// 安全起见，生产环境运行时请注意修改
var RequestEncryptKey = []byte(g.Cfg().MustGet(nil, "system.encryptKey", "f080a463654b2279").String())

// ConfigMaskDemoField 演示环境下需要隐藏的配置
var ConfigMaskDemoField = map[string]struct{}{
	// 邮箱
	"smtpUser": {}, "smtpPass": {},

	// 云存储
	"uploadUCloudPublicKey": {}, "uploadUCloudPrivateKey": {}, "uploadCosSecretId": {}, "uploadCosSecretKey": {},
	"uploadOssSecretId": {}, "uploadOssSecretKey": {}, "uploadQiNiuAccessKey": {}, "uploadQiNiuSecretKey": {},

	// 地图
	"geoAmapWebKey": {},

	// 短信
	"smsAliYunAccessKeyID": {}, "smsAliYunAccessKeySecret": {}, "smsTencentSecretId": {}, "smsTencentSecretKey": {},

	// 微信
	"officialAccountAppSecret": {}, "officialAccountToken": {}, "officialAccountEncodingAESKey": {}, "openPlatformAppSecret": {},
	"openPlatformToken": {}, "openPlatformEncodingAESKey": {},
}
