# 配置工具

## 配置工具概述

配置工具提供了统一的配置管理机制，位于`utility/config`目录下。通过该工具，可以轻松加载、访问和管理项目中的各种配置项，支持从多种来源（文件、环境变量、远程配置中心等）获取配置，并支持多环境配置。

## 主要功能

1. **多格式配置文件支持**：支持 YAML、JSON、TOML 等格式的配置文件
2. **多环境配置**：开发、测试、生产等不同环境的配置管理
3. **配置热重载**：支持配置动态更新，无需重启应用
4. **配置覆盖机制**：环境变量可以覆盖文件配置
5. **类型安全访问**：支持强类型配置访问
6. **配置加密**：支持敏感配置的加密存储和解密使用

## 配置文件结构

配置文件位于`configs`目录下，通常按以下结构组织：

```
configs/
  ├── config.yaml           # 基础配置
  ├── config.dev.yaml       # 开发环境特定配置
  ├── config.test.yaml      # 测试环境特定配置
  ├── config.prod.yaml      # 生产环境特定配置
  └── config.local.yaml     # 本地开发配置（不提交到git）
```

## 基本配置示例

一个典型的配置文件（YAML 格式）内容如下：

```yaml
# 应用基本配置
app:
  name: "支付通道服务"
  version: "1.0.0"
  mode: "dev" # dev, test, prod

# 服务器配置
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 60
  write_timeout: 60

# 数据库配置
database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "payment"
  max_idle: 10
  max_open: 100
  log_level: "info" # debug, info, warn, error

# Redis配置
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10

# 日志配置
logger:
  level: "info" # debug, info, warn, error
  format: "json" # text, json
  output: "stdout" # stdout, file
  file_path: "logs/app.log"

# 支付网关配置
payment:
  alipay:
    app_id: "2021000000000000"
    private_key: "MIIEvgIBADANBgkq..."
    alipay_public_key: "MIIBIjANBgkq..."
    notify_url: "https://api.example.com/payment/alipay/notify"
    return_url: "https://www.example.com/payment/return"

  wechat:
    app_id: "wx1234567890"
    mch_id: "1234567890"
    mch_key: "abcdefghijklmn1234567890ABCDEFG"
    notify_url: "https://api.example.com/payment/wechat/notify"
```

## 配置加载与初始化

### 初始化配置

```go
package main

import (
    "fmt"
    "github.com/yourusername/yourproject/utility/config"
)

func main() {
    // 初始化配置
    if err := config.Init("configs/config.yaml"); err != nil {
        panic(fmt.Sprintf("配置初始化失败: %v", err))
    }

    // 应用启动...
}
```

### 带环境的配置初始化

```go
// 初始化配置，自动根据环境加载对应配置文件
env := "dev" // 也可以从环境变量获取: os.Getenv("APP_ENV")
if err := config.InitWithEnv("configs/config.yaml", env); err != nil {
    panic(fmt.Sprintf("配置初始化失败: %v", err))
}
```

## 访问配置

### 获取字符串配置

```go
// 获取应用名称
appName := config.GetString("app.name")
fmt.Println("应用名称:", appName)

// 获取带默认值的配置
serverHost := config.GetString("server.host", "localhost")
```

### 获取数值类型配置

```go
// 获取整数配置
port := config.GetInt("server.port")
fmt.Println("服务端口:", port)

// 获取浮点数配置
timeout := config.GetFloat64("server.timeout")

// 获取布尔值配置
debug := config.GetBool("app.debug")
```

### 获取复杂类型配置

```go
// 获取字符串切片
allowedOrigins := config.GetStringSlice("cors.allowed_origins")

// 获取字符串映射
headers := config.GetStringMap("http.headers")

// 获取时间配置
duration := config.GetDuration("cache.expiration")
```

### 获取嵌套配置

```go
// 获取嵌套配置
dbConfig := config.GetStringMap("database")
fmt.Println("数据库主机:", dbConfig["host"])
fmt.Println("数据库端口:", dbConfig["port"])

// 或者直接使用点表示法
dbHost := config.GetString("database.host")
dbPort := config.GetInt("database.port")
```

### 绑定到结构体

```go
// 定义配置结构体
type DatabaseConfig struct {
    Driver   string `mapstructure:"driver"`
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    Database string `mapstructure:"database"`
    MaxIdle  int    `mapstructure:"max_idle"`
    MaxOpen  int    `mapstructure:"max_open"`
    LogLevel string `mapstructure:"log_level"`
}

// 绑定配置到结构体
var dbConfig DatabaseConfig
if err := config.UnmarshalKey("database", &dbConfig); err != nil {
    fmt.Println("绑定配置失败:", err)
    return
}

fmt.Printf("数据库配置: %+v\n", dbConfig)
```

## 配置热重载

```go
// 注册配置变更回调
config.OnConfigChange(func() {
    fmt.Println("配置已更新")
    // 重新加载受影响的组件
    reloadLogger()
    reloadDatabase()
})

// 启用配置热重载
config.WatchConfig()
```

## 环境变量覆盖

环境变量可以覆盖配置文件中的设置，例如：

```bash
# 设置环境变量
export APP_SERVER_PORT=9090
export APP_DATABASE_PASSWORD=strongpassword
```

在代码中，可以设置环境变量前缀：

```go
config.SetEnvPrefix("APP")
config.AutomaticEnv()

// 现在 APP_SERVER_PORT 环境变量会覆盖 server.port 配置
port := config.GetInt("server.port") // 返回 9090
```

## 加密配置

对于敏感信息（如密码、密钥等），可以使用加密存储：

```go
// 加密配置（使用AES加密）
encryptedPassword := config.Encrypt("database.password", "my-secret-password")
// 将加密后的值保存到配置文件

// 使用时自动解密
password := config.GetDecrypted("database.password")
```

## 配置分组

对于大型应用，可以对配置进行分组：

```go
// 获取数据库配置子集
dbConfig := config.Sub("database")

// 从子集获取配置
dbHost := dbConfig.GetString("host")
dbPort := dbConfig.GetInt("port")
```

## 配置中心集成

支持从配置中心（如 etcd、consul 等）加载配置：

```go
// 从etcd加载配置
config.AddRemoteProvider("etcd", "http://etcd-server:2379", "/configs/myapp")
config.ReadRemoteConfig()

// 启用远程配置自动更新
config.WatchRemoteConfig()
```

## 多种配置源优先级

配置工具支持多种配置源，按以下优先级使用（从高到低）：

1. 命令行参数
2. 环境变量
3. 远程配置中心
4. 特定环境配置文件（如 config.prod.yaml）
5. 基础配置文件（config.yaml）
6. 默认值

## 完整示例

### 配置初始化与使用

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/yourusername/yourproject/utility/config"
)

// 应用配置结构
type AppConfig struct {
    App struct {
        Name    string `mapstructure:"name"`
        Version string `mapstructure:"version"`
        Mode    string `mapstructure:"mode"`
        Debug   bool   `mapstructure:"debug"`
    } `mapstructure:"app"`

    Server struct {
        Host         string `mapstructure:"host"`
        Port         int    `mapstructure:"port"`
        ReadTimeout  int    `mapstructure:"read_timeout"`
        WriteTimeout int    `mapstructure:"write_timeout"`
    } `mapstructure:"server"`

    Database struct {
        Driver   string `mapstructure:"driver"`
        Host     string `mapstructure:"host"`
        Port     int    `mapstructure:"port"`
        Username string `mapstructure:"username"`
        Password string `mapstructure:"password"`
        Database string `mapstructure:"database"`
    } `mapstructure:"database"`

    // 其他配置...
}

var Conf AppConfig

func init() {
    // 设置环境变量前缀
    config.SetEnvPrefix("APP")
    config.AutomaticEnv()

    // 确定环境
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "dev" // 默认为开发环境
    }

    // 加载配置文件
    configPath := fmt.Sprintf("configs/config.%s.yaml", env)
    if err := config.Init(configPath); err != nil {
        // 尝试使用基础配置
        if err := config.Init("configs/config.yaml"); err != nil {
            log.Fatalf("无法加载配置文件: %v", err)
        }
    }

    // 绑定配置到结构体
    if err := config.Unmarshal(&Conf); err != nil {
        log.Fatalf("无法解析配置: %v", err)
    }

    // 启用配置热重载
    config.OnConfigChange(func() {
        fmt.Println("配置已更新，重新加载...")
        if err := config.Unmarshal(&Conf); err != nil {
            log.Printf("更新配置失败: %v", err)
            return
        }
        // 重新初始化受影响的组件
    })

    config.WatchConfig()
}

func main() {
    fmt.Printf("应用名称: %s\n", Conf.App.Name)
    fmt.Printf("应用版本: %s\n", Conf.App.Version)
    fmt.Printf("运行模式: %s\n", Conf.App.Mode)
    fmt.Printf("服务地址: %s:%d\n", Conf.Server.Host, Conf.Server.Port)

    // 应用逻辑...
}
```

## 配置安全管理

### 敏感信息处理

1. **避免硬编码**：不要在代码中硬编码敏感信息
2. **使用环境变量**：生产环境的敏感信息应通过环境变量注入
3. **加密存储**：配置文件中的敏感信息应该加密存储
4. **权限隔离**：确保配置文件的访问权限受到严格控制

### 配置加密示例

```go
// 配置加密器
type ConfigEncryptor struct {
    key []byte
}

// 创建新的加密器
func NewConfigEncryptor(key string) *ConfigEncryptor {
    return &ConfigEncryptor{
        key: []byte(key),
    }
}

// 加密配置值
func (e *ConfigEncryptor) Encrypt(value string) (string, error) {
    // 实现AES加密
    // ...
    return encryptedValue, nil
}

// 解密配置值
func (e *ConfigEncryptor) Decrypt(value string) (string, error) {
    // 实现AES解密
    // ...
    return decryptedValue, nil
}

// 使用示例
encryptor := NewConfigEncryptor(os.Getenv("CONFIG_ENCRYPTION_KEY"))

// 加密数据库密码
dbPassword := "super-secret-password"
encryptedPassword, err := encryptor.Encrypt(dbPassword)
if err != nil {
    log.Fatal("加密失败:", err)
}
fmt.Println("加密后的密码:", encryptedPassword)

// 保存加密后的密码到配置文件
// ...

// 使用时解密
decryptedPassword, err := encryptor.Decrypt(config.GetString("database.password"))
if err != nil {
    log.Fatal("解密失败:", err)
}
fmt.Println("解密后的密码:", decryptedPassword)
```

## 最佳实践

1. **使用版本控制管理配置模板**

   将配置模板（不含敏感信息）提交到版本控制系统，敏感信息使用环境变量或单独的安全存储方式管理。

2. **使用环境特定配置文件**

   为不同环境（dev、test、prod）创建特定的配置文件，以满足不同环境的需求。

3. **使用默认值**

   为配置项提供合理的默认值，使应用在缺少某些配置时仍能正常运行。

4. **集中管理配置**

   对于分布式系统，考虑使用配置中心统一管理配置。

5. **配置验证**

   在应用启动时验证关键配置的有效性，确保配置错误不会导致运行时故障。

6. **文档化配置**

   为配置项提供详细的文档，包括用途、可选值、默认值等信息。

## 常见问题解答

1. **如何处理配置文件不存在的情况？**

   ```go
   // 检查配置文件是否存在
   if _, err := os.Stat(configPath); os.IsNotExist(err) {
       // 创建默认配置
       createDefaultConfig(configPath)
   }
   ```

2. **如何在容器环境中管理配置？**

   在容器环境中，通常使用环境变量或配置卷挂载来提供配置：

   ```go
   // 优先使用环境变量
   config.AutomaticEnv()

   // 然后尝试从挂载的配置文件加载
   configPath := "/etc/app/config.yaml"
   if _, err := os.Stat(configPath); err == nil {
       config.SetConfigFile(configPath)
       config.ReadInConfig()
   }
   ```

3. **如何处理配置变更时的服务重启？**

   ```go
   config.OnConfigChange(func() {
       // 检查关键配置是否变更
       if configRequiresRestart() {
           log.Println("检测到关键配置变更，准备重启服务...")
           // 触发优雅重启流程
           gracefulRestart()
       } else {
           // 仅重新加载受影响的组件
           reloadAffectedComponents()
       }
   })
   ```
