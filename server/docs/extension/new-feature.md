# 新功能开发指南

本指南将帮助开发者在现有架构基础上开发新功能，确保与现有代码风格和架构保持一致。

## 开发流程概述

1. **分析需求**: 明确新功能的需求和边界
2. **设计接口**: 设计符合项目规范的 API 接口
3. **实现功能**: 按照分层架构实现功能
4. **测试**: 进行单元测试和集成测试
5. **文档**: 更新 API 文档和开发文档

## 实战案例：添加订单查询功能

下面以添加"订单查询"功能为例，演示完整开发流程：

### 1. 定义 API 接口

在`internal/api/api/payment/v1/`目录下创建或修改接口定义：

```go
// internal/api/api/payment/v1/query.go
package v1

import (
    "client-app/internal/model/input/sysin"
    "context"
)

// 请求参数
type OrderQueryReq struct {
    OrderQueryInp
}

// 响应结构
type OrderQueryRes struct {
    OrderQueryModel
}
```

### 2. 定义输入参数和响应模型

在`internal/model/input/sysin/`目录下添加模型定义：

```go
// internal/model/input/sysin/order_query.go
package sysin

import "context"

// 查询输入参数
type OrderQueryInp struct {
    OrderID string `json:"orderId" v:"required#订单ID不能为空"`
}

// 参数过滤器
func (in *OrderQueryInp) Filter(ctx context.Context) (err error) {
    // 参数验证逻辑
    return nil
}

// 查询结果模型
type OrderQueryModel struct {
    OrderID     string `json:"orderId"`     // 订单ID
    Amount      string `json:"amount"`      // 订单金额
    Status      int    `json:"status"`      // 订单状态
    CreatedTime string `json:"createdTime"` // 创建时间
}
```

### 3. 更新服务接口

在`internal/service/api.payment.go`中添加新方法：

```go
// internal/service/api.payment.go
type IPayment interface {
    Submit(ctx context.Context, in *sysin.PaymentInp) (res *sysin.PaymentModel, err error)
    // 添加新方法
    QueryOrder(ctx context.Context, in *sysin.OrderQueryInp) (res *sysin.OrderQueryModel, err error)
}
```

### 4. 实现业务逻辑

在`internal/logic/api/payment.go`中实现新方法：

```go
// internal/logic/api/payment.go
func (s *sPayment) QueryOrder(ctx context.Context, req *sysin.OrderQueryInp) (*sysin.OrderQueryModel, error) {
    // 查询订单逻辑
    result := &sysin.OrderQueryModel{
        OrderID:     req.OrderID,
        Amount:      "100.00",
        Status:      1,
        CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
    }
    return result, nil
}
```

### 5. 添加控制器方法

在`internal/controller/api/payment.go`中添加新方法：

```go
// internal/controller/api/payment.go
func (c *Payment) QueryOrder(ctx context.Context, req *v1.OrderQueryReq) (res *v1.OrderQueryRes, err error) {
    out, err := service.Payment().QueryOrder(ctx, &req.OrderQueryInp)
    if err != nil {
        return nil, err
    }
    res = new(v1.OrderQueryRes)
    if g.IsEmpty(out) {
        return res, nil
    }
    res.OrderQueryModel = *out
    return res, nil
}
```

### 6. 更新路由接口

在`internal/api/api/payment/payment.go`中添加新方法：

```go
// internal/api/api/payment/payment.go
type IPayment interface {
    Submit(ctx context.Context, req *v1.PaymentReq) (res *v1.PaymentRes, err error)
    // 添加新方法
    QueryOrder(ctx context.Context, req *v1.OrderQueryReq) (res *v1.OrderQueryRes, err error)
}
```

### 7. 测试新功能

使用 HTTP 客户端测试新添加的 API：

```
POST /api/payment/queryOrder
Content-Type: application/json

{
  "orderId": "ORDER123456789"
}
```

## 新功能开发注意事项

1. **命名规范**: 遵循项目现有命名规范
2. **参数验证**: 输入参数需实现`Filter`接口进行验证
3. **错误处理**: 统一使用项目的错误处理机制
4. **文档更新**: 完成功能后更新相关文档
5. **单元测试**: 为新功能编写充分的单元测试
