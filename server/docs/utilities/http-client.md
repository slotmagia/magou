# HTTP 客户端工具

## HTTP 客户端概述

HTTP 客户端工具提供了用于发起 HTTP 请求的工具包，位于`utility/httpclient`目录下。该工具封装了标准库的`http`包，并提供了更便捷的 API，用于与外部服务进行通信。

## 主要功能

1. **发送 HTTP 请求**：支持 GET、POST、PUT、DELETE 等常见 HTTP 方法
2. **请求参数设置**：支持 URL 参数、表单数据、JSON 数据等
3. **响应处理**：自动解析 JSON 响应、处理错误
4. **超时控制**：设置连接超时和请求超时
5. **重试机制**：支持请求失败自动重试
6. **TLS 配置**：支持 HTTPS 请求和证书配置

## 基本使用

### 创建 HTTP 客户端

```go
// 创建默认HTTP客户端
client := httpclient.New()

// 创建带配置的HTTP客户端
client := httpclient.New(
    httpclient.WithTimeout(5*time.Second),
    httpclient.WithRetry(3, 1*time.Second),
)
```

### 发送 GET 请求

```go
// 简单GET请求
resp, err := client.Get("https://api.example.com/users")
if err != nil {
    // 处理错误
    return
}

// 带URL参数的GET请求
resp, err := client.Get("https://api.example.com/users",
    httpclient.WithQuery(map[string]string{
        "page": "1",
        "limit": "10",
    }),
)
```

### 发送 POST 请求

```go
// 发送JSON数据
resp, err := client.Post("https://api.example.com/users",
    httpclient.WithJSON(map[string]interface{}{
        "name": "张三",
        "age": 30,
    }),
)

// 发送表单数据
resp, err := client.Post("https://api.example.com/login",
    httpclient.WithForm(map[string]string{
        "username": "admin",
        "password": "123456",
    }),
)
```

### 处理响应

```go
// 解析JSON响应到结构体
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

var user User
resp, err := client.Get("https://api.example.com/users/1")
if err != nil {
    return
}
defer resp.Body.Close()

if err := resp.DecodeJSON(&user); err != nil {
    return
}
fmt.Printf("用户名: %s, 年龄: %d\n", user.Name, user.Age)

// 检查响应状态码
if resp.StatusCode != http.StatusOK {
    fmt.Println("请求失败，状态码:", resp.StatusCode)
    return
}

// 读取响应body
body, err := resp.ReadBody()
if err != nil {
    return
}
fmt.Println("响应内容:", string(body))
```

## 高级功能

### 请求头设置

```go
resp, err := client.Get("https://api.example.com/users",
    httpclient.WithHeaders(map[string]string{
        "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "Content-Type":  "application/json",
        "Accept":        "application/json",
    }),
)
```

### 超时控制

```go
// 设置客户端全局超时
client := httpclient.New(
    httpclient.WithTimeout(5*time.Second),
)

// 单次请求的超时控制
resp, err := client.Get("https://api.example.com/users",
    httpclient.WithRequestTimeout(2*time.Second),
)
```

### 重试机制

```go
// 设置重试次数和间隔
client := httpclient.New(
    httpclient.WithRetry(3, 1*time.Second),
)

// 自定义重试条件
client := httpclient.New(
    httpclient.WithRetryCondition(func(resp *http.Response, err error) bool {
        return err != nil || resp.StatusCode >= 500
    }),
)
```

### TLS 配置

```go
// 跳过证书验证（不推荐用于生产环境）
client := httpclient.New(
    httpclient.WithInsecureSkipVerify(true),
)

// 使用自定义证书
cert, _ := tls.LoadX509KeyPair("client.crt", "client.key")
client := httpclient.New(
    httpclient.WithTLSClientConfig(&tls.Config{
        Certificates: []tls.Certificate{cert},
    }),
)
```

### 代理设置

```go
// 设置HTTP代理
client := httpclient.New(
    httpclient.WithProxy("http://proxy.example.com:8080"),
)
```

### 请求上下文

```go
// 创建可取消的上下文
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// 使用上下文发送请求
resp, err := client.GetWithContext(ctx, "https://api.example.com/users")

// 10秒后自动取消请求
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
resp, err := client.GetWithContext(ctx, "https://api.example.com/users")
```

## 完整示例

### 支付 API 调用示例

```go
// 支付服务客户端
type PaymentClient struct {
    httpClient *httpclient.Client
    baseURL    string
    appID      string
    appSecret  string
}

// 创建支付客户端
func NewPaymentClient(baseURL, appID, appSecret string) *PaymentClient {
    return &PaymentClient{
        httpClient: httpclient.New(
            httpclient.WithTimeout(5*time.Second),
            httpclient.WithRetry(3, 1*time.Second),
        ),
        baseURL:    baseURL,
        appID:      appID,
        appSecret:  appSecret,
    }
}

// 生成签名
func (c *PaymentClient) generateSign(params map[string]interface{}) string {
    // 实现签名逻辑
    return "signature"
}

// 创建支付订单
func (c *PaymentClient) CreateOrder(orderID string, amount float64, description string) (*OrderResult, error) {
    // 构建请求参数
    params := map[string]interface{}{
        "app_id":      c.appID,
        "order_id":    orderID,
        "amount":      amount,
        "description": description,
        "timestamp":   time.Now().Unix(),
    }

    // 生成签名
    params["sign"] = c.generateSign(params)

    // 发送请求
    resp, err := c.httpClient.Post(c.baseURL+"/api/orders",
        httpclient.WithJSON(params),
        httpclient.WithHeaders(map[string]string{
            "Content-Type": "application/json",
        }),
    )
    if err != nil {
        return nil, fmt.Errorf("请求支付API失败: %w", err)
    }
    defer resp.Body.Close()

    // 检查响应状态
    if resp.StatusCode != http.StatusOK {
        body, _ := resp.ReadBody()
        return nil, fmt.Errorf("支付API返回错误: %s, 状态码: %d", string(body), resp.StatusCode)
    }

    // 解析响应
    var result struct {
        Code    int         `json:"code"`
        Message string      `json:"message"`
        Data    OrderResult `json:"data"`
    }

    if err := resp.DecodeJSON(&result); err != nil {
        return nil, fmt.Errorf("解析支付API响应失败: %w", err)
    }

    // 检查业务状态码
    if result.Code != 0 {
        return nil, fmt.Errorf("支付API业务错误: %s, 错误码: %d", result.Message, result.Code)
    }

    return &result.Data, nil
}

// 订单结果
type OrderResult struct {
    OrderID     string `json:"order_id"`
    PaymentURL  string `json:"payment_url"`
    QRCodeURL   string `json:"qrcode_url"`
    ExpireTime  int64  `json:"expire_time"`
}

// 使用示例
func main() {
    client := NewPaymentClient(
        "https://api.payment.example.com",
        "APP_ID_123456",
        "APP_SECRET_abcdef",
    )

    result, err := client.CreateOrder("ORDER_123456", 100.50, "商品购买")
    if err != nil {
        fmt.Println("创建订单失败:", err)
        return
    }

    fmt.Println("订单创建成功:")
    fmt.Println("订单ID:", result.OrderID)
    fmt.Println("支付链接:", result.PaymentURL)
    fmt.Println("二维码链接:", result.QRCodeURL)
    fmt.Println("过期时间:", time.Unix(result.ExpireTime, 0).Format("2006-01-02 15:04:05"))
}
```

## 最佳实践

1. **连接池复用**

   HTTP 客户端默认使用连接池，避免频繁创建和关闭连接。在应用中应该创建一个全局的客户端实例并复用它。

   ```go
   var defaultClient = httpclient.New(
       httpclient.WithTimeout(5*time.Second),
       httpclient.WithRetry(3, 1*time.Second),
   )
   ```

2. **错误处理**

   始终检查 HTTP 请求的错误并正确处理：

   ```go
   resp, err := client.Get(url)
   if err != nil {
       log.Printf("请求失败: %v", err)
       return
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
       body, _ := resp.ReadBody()
       log.Printf("服务器返回错误: %s, 状态码: %d", string(body), resp.StatusCode)
       return
   }
   ```

3. **超时设置**

   总是设置合理的超时时间，防止请求长时间阻塞：

   ```go
   client := httpclient.New(
       httpclient.WithTimeout(5*time.Second), // 总超时
       httpclient.WithDialTimeout(2*time.Second), // 连接超时
   )
   ```

4. **请求日志**

   记录请求和响应日志，便于问题排查：

   ```go
   client := httpclient.New(
       httpclient.WithLogger(log.New(os.Stdout, "HTTP ", log.LstdFlags)),
       httpclient.WithLogLevel(httpclient.LogLevelDebug),
   )
   ```

5. **服务监控**

   监控 HTTP 请求的成功率、延迟等指标：

   ```go
   client := httpclient.New(
       httpclient.WithMetrics(metrics.DefaultRegistry),
   )
   ```

## 常见问题解答

1. **如何处理大文件上传？**

   ```go
   file, _ := os.Open("large_file.zip")
   defer file.Close()

   resp, err := client.Post("https://api.example.com/upload",
       httpclient.WithMultipartField("file", "large_file.zip", file),
   )
   ```

2. **如何下载大文件？**

   ```go
   resp, err := client.Get("https://api.example.com/download/large_file.zip")
   if err != nil {
       return
   }
   defer resp.Body.Close()

   // 创建目标文件
   out, _ := os.Create("downloaded_file.zip")
   defer out.Close()

   // 复制响应内容到文件
   _, err = io.Copy(out, resp.Body)
   ```

3. **如何实现 API 限流？**

   ```go
   // 创建限流器，每秒最多10个请求
   limiter := rate.NewLimiter(10, 1)

   client := httpclient.New(
       httpclient.WithRequestInterceptor(func(req *http.Request) error {
           return limiter.Wait(req.Context())
       }),
   )
   ```
