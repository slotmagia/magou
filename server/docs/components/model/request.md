# 请求模型

## 请求模型概述

请求模型定义了 API 接口接收的数据结构，负责参数验证和数据预处理。在本项目中，请求模型位于两个位置：

- API 接口定义：`internal/api/api/*/v1/`
- 输入参数定义：`internal/model/input/sysin/`

## 请求模型层次结构

请求模型采用嵌套结构设计，分为两层：

1. **API 接口请求结构体**：定义 API 接口参数，用于路由绑定
2. **输入参数结构体**：定义业务逻辑输入参数，用于服务层处理

```go
// API接口请求结构体（internal/api/api/payment/v1/payment.go）
type PaymentReq struct {
    PaymentInp // 嵌套输入参数结构体
}

// 输入参数结构体（internal/model/input/sysin/payment.go）
type PaymentInp struct {
    Order string `json:"order" v:"required#订单不能为空"`
}
```

## 参数验证机制

项目使用 GoFrame 的验证框架，通过结构体标签定义验证规则：

```go
type PaymentInp struct {
    Order string `json:"order" v:"required#订单不能为空"`
    Amount float64 `json:"amount" v:"required|min:0.01#金额不能为空#金额必须大于0.01"`
    PaymentMethod string `json:"paymentMethod" v:"required|in:alipay,wechat,creditcard#支付方式不能为空#不支持的支付方式"`
}
```

常用验证规则：

| 规则     | 说明       | 示例                                     |
| -------- | ---------- | ---------------------------------------- |
| required | 必填项     | `v:"required#用户名不能为空"`            |
| min      | 最小值     | `v:"min:0.01#金额必须大于0.01"`          |
| max      | 最大值     | `v:"max:100#金额不能超过100"`            |
| length   | 长度要求   | `v:"length:6,20#密码长度必须在6-20之间"` |
| in       | 枚举值     | `v:"in:alipay,wechat#不支持的支付方式"`  |
| regex    | 正则表达式 | `v:"regex:\\d{11}#手机号格式不正确"`     |

## 自定义参数验证

除了标签验证外，还可以通过实现`Filter`接口进行自定义验证：

```go
// 参数过滤器接口
type Filter interface {
    Filter(ctx context.Context) error
}

// 实现Filter接口
func (in *PaymentInp) Filter(ctx context.Context) (err error) {
    // 自定义验证逻辑
    if in.Amount <= 0 {
        return errors.New("金额必须大于0")
    }
    return nil
}
```

## 参数预处理

在`Filter`方法中，除了验证参数，还可以进行参数预处理：

```go
func (in *OrderQueryInp) Filter(ctx context.Context) (err error) {
    // 参数验证
    if len(in.OrderID) == 0 {
        return errors.New("订单ID不能为空")
    }

    // 参数预处理
    in.OrderID = strings.TrimSpace(in.OrderID)
    return nil
}
```

## 命名规范

请求模型的命名规范：

1. **API 请求结构体**：功能名+Req，如`PaymentReq`
2. **输入参数结构体**：功能名+Inp，如`PaymentInp`
3. **JSON 标签**：使用小驼峰命名，如`json:"orderNumber"`
4. **验证标签**：使用竖线分隔多个规则，使用井号分隔错误消息，如`v:"required|min:0.01#订单号不能为空#金额必须大于0.01"`

## 请求模型示例

完整的请求模型定义示例：

```go
// API请求结构体（internal/api/api/payment/v1/query.go）
package v1

import (
    "client-app/internal/model/input/sysin"
    "context"
)

type OrderQueryReq struct {
    sysin.OrderQueryInp
}

// 输入参数结构体（internal/model/input/sysin/order_query.go）
package sysin

import (
    "context"
    "errors"
    "strings"
)

type OrderQueryInp struct {
    OrderID string `json:"orderId" v:"required#订单ID不能为空"`
    QueryType int `json:"queryType" v:"in:0,1,2#查询类型不正确"` // 0:全部, 1:基础信息, 2:详细信息
}

func (in *OrderQueryInp) Filter(ctx context.Context) (err error) {
    // 参数验证与预处理
    in.OrderID = strings.TrimSpace(in.OrderID)
    if len(in.OrderID) == 0 {
        return errors.New("订单ID不能为空")
    }

    // 附加业务规则验证
    if in.QueryType == 2 {
        // 如果是查询详细信息，需要进行额外的权限检查
        // ...
    }

    return nil
}
```

## 请求参数解析流程

1. 客户端发送请求
2. 框架解析请求数据并绑定到请求结构体
3. 执行标签验证
4. 如果实现了`Filter`接口，调用`Filter`方法进行自定义验证
5. 验证通过后，请求参数传递给控制器方法
6. 控制器调用服务层处理业务逻辑

这个流程由中间件`PreFilter`自动处理，开发者不需要手动调用验证逻辑。
