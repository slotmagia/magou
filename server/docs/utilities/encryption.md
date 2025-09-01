# 加密工具

## 加密工具概述

加密工具提供了项目中使用的各种加密、解密和哈希功能，位于`utility/encrypt`目录下。这些工具用于保护敏感数据、生成签名和验证请求。

## 加密工具功能

加密工具库提供以下功能：

1. **MD5 加密**：用于生成数据的 MD5 哈希值
2. **AES 加密/解密**：用于数据的对称加密和解密
3. **RSA 加密/解密**：用于数据的非对称加密和解密
4. **签名生成与验证**：用于 API 签名

## MD5 加密

MD5 用于生成不可逆的哈希值，常用于密码存储和数据完整性验证：

```go
// MD5加密
func MD5(data string) string {
    h := md5.New()
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}

// 使用示例
password := "123456"
hashedPassword := encrypt.MD5(password)
fmt.Println(hashedPassword) // e10adc3949ba59abbe56e057f20f883e
```

## AES 加密/解密

AES 是对称加密算法，用于加密和解密敏感数据：

```go
// AES加密
func EncryptByAes(plainText, key string) (string, error) {
    // 实现省略
    return cipherText, nil
}

// AES解密
func DecryptByAes(cipherText, key string) (string, error) {
    // 实现省略
    return plainText, nil
}

// 使用示例
key := "123456789abcdef0123456789abcdef0" // 32字节密钥
plainText := "sensitive data"
cipherText, err := encrypt.EncryptByAes(plainText, key)
if err != nil {
    fmt.Println("加密失败:", err)
    return
}
fmt.Println("加密结果:", cipherText)

decrypted, err := encrypt.DecryptByAes(cipherText, key)
if err != nil {
    fmt.Println("解密失败:", err)
    return
}
fmt.Println("解密结果:", decrypted)
```

## RSA 加密/解密

RSA 是非对称加密算法，用于需要更高安全性的场景：

```go
// RSA公钥加密
func RsaEncrypt(plainText string, publicKey []byte) (string, error) {
    // 实现省略
    return cipherText, nil
}

// RSA私钥解密
func RsaDecrypt(cipherText string, privateKey []byte) (string, error) {
    // 实现省略
    return plainText, nil
}

// 使用示例
publicKey := []byte("-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----")
privateKey := []byte("-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----")

plainText := "sensitive data"
cipherText, err := encrypt.RsaEncrypt(plainText, publicKey)
if err != nil {
    fmt.Println("加密失败:", err)
    return
}
fmt.Println("加密结果:", cipherText)

decrypted, err := encrypt.RsaDecrypt(cipherText, privateKey)
if err != nil {
    fmt.Println("解密失败:", err)
    return
}
fmt.Println("解密结果:", decrypted)
```

## API 签名生成与验证

API 签名用于验证请求的合法性和完整性：

```go
// 生成API签名
func GenerateSign(params map[string]interface{}, secretKey string) string {
    // 按键名升序排序
    var keys []string
    for k := range params {
        if k != "sign" { // 排除sign字段
            keys = append(keys, k)
        }
    }
    sort.Strings(keys)

    // 构建签名字符串
    var signStr string
    for _, k := range keys {
        v := params[k]
        if v != "" { // 忽略空值
            signStr += k + "=" + fmt.Sprintf("%v", v) + "&"
        }
    }
    signStr += "key=" + secretKey

    // 生成MD5签名
    return MD5(signStr)
}

// 验证API签名
func VerifySign(params map[string]interface{}, sign, secretKey string) bool {
    generatedSign := GenerateSign(params, secretKey)
    return sign == generatedSign
}

// 使用示例
params := map[string]interface{}{
    "order_id": "ORDER123456",
    "amount": "100.00",
    "timestamp": time.Now().Unix(),
}
secretKey := "your_secret_key"
sign := encrypt.GenerateSign(params, secretKey)
params["sign"] = sign

// 验证签名
isValid := encrypt.VerifySign(params, sign, secretKey)
fmt.Println("签名验证结果:", isValid)
```

## Base64 编码/解码

Base64 用于二进制数据的编码和解码：

```go
// Base64编码
func Base64Encode(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}

// Base64解码
func Base64Decode(s string) ([]byte, error) {
    return base64.StdEncoding.DecodeString(s)
}

// 使用示例
data := []byte("hello world")
encoded := encrypt.Base64Encode(data)
fmt.Println("编码结果:", encoded)

decoded, err := encrypt.Base64Decode(encoded)
if err != nil {
    fmt.Println("解码失败:", err)
    return
}
fmt.Println("解码结果:", string(decoded))
```

## 随机字符串生成

用于生成随机字符串，常用于临时密钥、会话 ID 等：

```go
// 生成随机字符串
func RandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}

// 使用示例
randomStr := encrypt.RandomString(16)
fmt.Println("随机字符串:", randomStr)
```

## 最佳实践

1. **敏感数据加密存储**：

   - 用户密码应使用不可逆哈希算法存储（如 bcrypt）
   - 敏感业务数据应使用 AES 加密存储

2. **API 安全通信**：

   - 使用 HTTPS 传输数据
   - 为 API 请求添加签名验证
   - 为敏感 API 添加时间戳防止重放攻击

3. **密钥管理**：

   - 密钥不应硬编码在源代码中
   - 使用环境变量或配置中心存储密钥
   - 定期轮换密钥

4. **日志安全**：
   - 避免在日志中记录敏感信息（如密码、完整信用卡号）
   - 对必须记录的敏感数据进行脱敏处理

## 示例应用场景

### 用户密码加密

```go
// 注册时加密存储密码
func registerUser(username, password string) error {
    hashedPassword := encrypt.MD5(password) // 实际应使用更安全的bcrypt
    // 存储username和hashedPassword到数据库
    return nil
}

// 登录时验证密码
func verifyPassword(username, password string) bool {
    // 从数据库获取存储的hashedPassword
    storedHash := getStoredHash(username)
    return encrypt.MD5(password) == storedHash
}
```

### API 请求签名验证

```go
// 客户端构造请求
func buildRequest() map[string]interface{} {
    params := map[string]interface{}{
        "method": "payment.create",
        "order_id": "ORDER123456",
        "amount": "100.00",
        "timestamp": time.Now().Unix(),
    }
    params["sign"] = encrypt.GenerateSign(params, "your_secret_key")
    return params
}

// 服务端验证请求
func verifyRequest(params map[string]interface{}) bool {
    sign, ok := params["sign"].(string)
    if !ok {
        return false
    }
    delete(params, "sign")
    return encrypt.VerifySign(params, sign, "your_secret_key")
}
```

### 敏感数据加密存储

```go
// 加密存储银行卡信息
func storeBankCard(cardNo string) (string, error) {
    // 获取加密密钥
    key := getEncryptionKey()
    // 加密卡号
    encryptedCardNo, err := encrypt.EncryptByAes(cardNo, key)
    if err != nil {
        return "", err
    }
    // 存储加密后的卡号
    return encryptedCardNo, nil
}

// 解密获取银行卡信息
func getBankCard(encryptedCardNo string) (string, error) {
    // 获取加密密钥
    key := getEncryptionKey()
    // 解密卡号
    cardNo, err := encrypt.DecryptByAes(encryptedCardNo, key)
    if err != nil {
        return "", err
    }
    return cardNo, nil
}
```
