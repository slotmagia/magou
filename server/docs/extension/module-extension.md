# 现有功能扩展指南

本指南用于帮助开发者对现有功能进行扩展或修改，确保与原有架构和设计理念保持一致。

## 扩展流程概述

1. **理解现有功能**: 深入了解要扩展的功能模块
2. **确定扩展点**: 识别合适的扩展点
3. **设计扩展方案**: 设计与现有架构兼容的扩展方案
4. **实现扩展**: 按照分层架构实现扩展功能
5. **测试与文档**: 进行测试并更新文档

## 实战案例：扩展支付功能

以下以扩展支付功能为例，演示如何在现有基础上增加支付方式：

### 1. 理解现有功能

首先分析现有支付功能的实现：

```go
// internal/model/input/sysin/payment.go
type PaymentInp struct {
    Order string `json:"order" v:"required#订单不能为空"`
}

// internal/logic/api/payment.go
func (s *sPayment) Submit(ctx context.Context, req *sysin.PaymentInp) (*sysin.PaymentModel, error) {
    fmt.Println("Submit(ctx context.Context, req *sysin.PaymentInp)")
    return nil, nil
}
```

### 2. 扩展输入参数

修改输入参数，增加支付方式字段：

```go
// internal/model/input/sysin/payment.go
type PaymentInp struct {
    Order      string `json:"order" v:"required#订单不能为空"`
    PaymentMethod string `json:"paymentMethod" v:"required#支付方式不能为空"`
}

func (in *PaymentInp) Filter(ctx context.Context) (err error) {
    // 验证支付方式是否合法
    validMethods := []string{"alipay", "wechat", "creditcard"}
    if !validate.InSlice(validMethods, in.PaymentMethod) {
        return errors.New("不支持的支付方式")
    }
    return nil
}
```

### 3. 扩展业务逻辑

修改支付业务逻辑，支持不同支付方式的处理：

```go
// internal/logic/api/payment.go
func (s *sPayment) Submit(ctx context.Context, req *sysin.PaymentInp) (*sysin.PaymentModel, error) {
    // 根据支付方式选择不同的处理逻辑
    switch req.PaymentMethod {
    case "alipay":
        return s.processAlipayPayment(ctx, req)
    case "wechat":
        return s.processWechatPayment(ctx, req)
    case "creditcard":
        return s.processCreditCardPayment(ctx, req)
    default:
        return nil, errors.New("不支持的支付方式")
    }
}

// 支付宝支付处理
func (s *sPayment) processAlipayPayment(ctx context.Context, req *sysin.PaymentInp) (*sysin.PaymentModel, error) {
    // 处理支付宝支付逻辑
    return &sysin.PaymentModel{
        // 填充支付结果
    }, nil
}

// 微信支付处理
func (s *sPayment) processWechatPayment(ctx context.Context, req *sysin.PaymentInp) (*sysin.PaymentModel, error) {
    // 处理微信支付逻辑
    return &sysin.PaymentModel{
        // 填充支付结果
    }, nil
}

// 信用卡支付处理
func (s *sPayment) processCreditCardPayment(ctx context.Context, req *sysin.PaymentInp) (*sysin.PaymentModel, error) {
    // 处理信用卡支付逻辑
    return &sysin.PaymentModel{
        // 填充支付结果
    }, nil
}
```

### 4. 扩展响应模型

修改响应模型，增加支付结果信息：

```go
// internal/model/input/sysin/payment.go
type PaymentModel struct {
    TradeNo       string `json:"tradeNo"`       // 交易号
    PaymentMethod string `json:"paymentMethod"` // 支付方式
    Status        int    `json:"status"`        // 支付状态
    Message       string `json:"message"`       // 状态描述
}
```

### 5. 测试扩展功能

使用 HTTP 客户端测试扩展后的 API：

```
POST /api/payment/submit
Content-Type: application/json

{
  "order": "ORDER123456",
  "paymentMethod": "alipay"
}
```

## 功能扩展注意事项

1. **向后兼容**: 尽量保持对现有 API 的向后兼容
2. **参数验证**: 新增参数需要添加相应的验证规则
3. **错误处理**: 统一处理新增的错误情况
4. **日志记录**: 为重要操作添加日志记录
5. **性能考虑**: 评估扩展对性能的影响
6. **单元测试**: 为扩展功能添加单元测试用例
