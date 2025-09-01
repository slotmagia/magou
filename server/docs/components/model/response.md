# 响应模型

## 响应模型概述

响应模型定义了 API 接口返回的数据结构，用于统一格式化 HTTP 响应。在本项目中，响应模型分为两层：

1. **统一响应封装**：`internal/model/response.go`中定义的`Response`结构体
2. **业务响应模型**：定义在 API 接口和输入参数中的具体业务模型

## 统一响应结构

所有 API 响应都使用统一的`Response`结构体进行封装：

```go
// internal/model/response.go
type Response struct {
    Code      int         `json:"code" example:"0" description:"状态码"`
    Message   string      `json:"message,omitempty" example:"操作成功" description:"提示消息"`
    Data      interface{} `json:"data,omitempty" description:"数据集"`
    Error     interface{} `json:"error,omitempty" description:"错误信息"`
    Timestamp int64       `json:"timestamp" example:"1640966400" description:"服务器时间戳"`
    TraceID   string      `json:"traceID" v:"0" example:"d0bb93048bc5c9164cdee845dcb7f820" description:"链路ID"`
}
```

响应字段说明：

| 字段名    | 类型   | 说明                               |
| --------- | ------ | ---------------------------------- |
| code      | int    | 状态码，0 表示成功，其他值表示失败 |
| message   | string | 提示消息，成功或失败的描述         |
| data      | object | 业务数据，成功时返回               |
| error     | object | 错误详情，失败时返回               |
| timestamp | int64  | 服务器时间戳                       |
| traceID   | string | 链路追踪 ID，用于问题排查          |

## 响应模型层次结构

响应模型采用与请求模型类似的嵌套结构设计：

1. **API 接口响应结构体**：定义 API 接口返回值，用于路由绑定
2. **输出参数结构体**：定义业务逻辑输出参数，包含具体数据

```go
// API接口响应结构体（internal/api/api/payment/v1/payment.go）
type PaymentRes struct {
    PaymentModel // 嵌套输出参数结构体
}

// 输出参数结构体（internal/model/input/sysin/payment.go）
type PaymentModel struct {
    TradeNo       string `json:"tradeNo"`       // 交易号
    PaymentMethod string `json:"paymentMethod"` // 支付方式
    Status        int    `json:"status"`        // 支付状态
    Message       string `json:"message"`       // 状态描述
}
```

## 响应处理机制

响应处理由中间件`ResponseHandler`自动完成，流程如下：

1. 控制器返回业务数据或错误
2. 中间件捕获返回值
3. 根据返回值类型进行判断：
   - 如果是错误，封装为错误响应
   - 如果是业务数据，封装为成功响应
4. 将封装后的响应写入 HTTP 响应体

```go
// 中间件处理示例
func (s *sMiddleware) ResponseHandler(r *ghttp.Request) {
    r.Middleware.Next()

    // 判断是否已经有响应处理
    if r.Response.BufferLength() > 0 {
        return
    }

    // 获取返回值
    var (
        err  = r.GetError()
        res  = r.GetHandlerResponse()
        code = gerror.Code(err)
    )

    // 构造响应
    response := &model.Response{
        Code:      code.Code(),
        Message:   gerror.Current(err).Error(),
        Data:      res,
        Timestamp: gtime.Timestamp(),
        TraceID:   r.GetTraceId(),
    }

    // 写入响应
    r.Response.WriteJson(response)
}
```

## 成功响应示例

成功响应的 JSON 格式：

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "tradeNo": "T2023032712345678",
    "paymentMethod": "alipay",
    "status": 1,
    "message": "支付成功"
  },
  "timestamp": 1679898000,
  "traceID": "abcdef123456"
}
```

## 错误响应示例

错误响应的 JSON 格式：

```json
{
  "code": 1001,
  "message": "订单不存在",
  "error": {
    "details": "未找到订单记录: ORDER123456"
  },
  "timestamp": 1679898000,
  "traceID": "abcdef123456"
}
```

## 空响应处理

如果业务逻辑返回空值，响应处理器会返回一个带有空`data`字段的成功响应：

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {},
  "timestamp": 1679898000,
  "traceID": "abcdef123456"
}
```

## 命名规范

响应模型的命名规范：

1. **API 响应结构体**：功能名+Res，如`PaymentRes`
2. **输出参数结构体**：功能名+Model，如`PaymentModel`
3. **JSON 标签**：使用小驼峰命名，如`json:"tradeNo"`
4. **字段注释**：字段后加注释说明用途，如`TradeNo string // 交易号`

## 响应模型示例

完整的响应模型定义示例：

```go
// API响应结构体（internal/api/api/payment/v1/query.go）
package v1

import (
    "client-app/internal/model/input/sysin"
)

type OrderQueryRes struct {
    sysin.OrderQueryModel
}

// 输出参数结构体（internal/model/input/sysin/order_query.go）
package sysin

type OrderQueryModel struct {
    OrderID     string `json:"orderId"`     // 订单ID
    Amount      string `json:"amount"`      // 订单金额
    Status      int    `json:"status"`      // 订单状态（0:未支付, 1:已支付, 2:已取消）
    CreatedTime string `json:"createdTime"` // 创建时间
    PaymentInfo struct {
        PaymentMethod string `json:"paymentMethod"` // 支付方式
        PaymentTime   string `json:"paymentTime"`   // 支付时间
        TradeNo       string `json:"tradeNo"`       // 交易号
    } `json:"paymentInfo,omitempty"` // 支付信息，仅在已支付状态返回
}
```

## 错误码规范

项目定义了统一的错误码，在`internal/consts`包中：

| 错误码 | 说明           |
| ------ | -------------- |
| 0      | 成功           |
| 1      | 失败(通用错误) |
| 401    | 未授权         |
| 403    | 禁止访问       |
| 404    | 资源不存在     |
| 500    | 服务器内部错误 |
| 1001   | 订单不存在     |
| 1002   | 订单状态错误   |
| ...    | 其他业务错误码 |

使用错误码示例：

```go
if order == nil {
    return nil, gerror.NewCode(consts.CodeOrderNotFound, "订单不存在")
}
```
