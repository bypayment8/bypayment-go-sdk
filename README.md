# ByPayment Go SDK

ByPayment Go SDK 是一个用于与 ByPayment API 服务器交互的 Go 语言客户端库。

## 功能特性

- ✅ 完整的 API 接口支持
- ✅ 自动签名生成和验证
- ✅ 防重放攻击保护（nonce + timestamp）
- ✅ 类型安全的请求和响应
- ✅ 简洁易用的 API 设计
- ✅ 回调签名验证和响应处理

## 安装

### 安装指定版本

```bash
go get github.com/bypayment8/bypayment-go-sdk@v1.0.0
```

## 接入指南

详细的接入指南请参考 [INTEGRATION.md](./INTEGRATION.md)，包含：

- 📦 安装方式（go get、go mod、本地开发）
- 🔌 不同场景的接入示例（Web 服务、后台任务、CLI、依赖注入）
- ⚙️ 配置管理（环境变量、配置文件、配置中心）
- 🛡️ 错误处理和重试机制
- ✨ 最佳实践和常见问题

快速查看接入示例：

```go
// 1. 导入 SDK
import "github.com/bypayment8/bypayment-go-sdk"

// 2. 创建客户端
client := bypayment.NewClient(
    "https://api.example.com",
    "your-api-key",
    "your-secret-key",
)

// 3. 调用 API
resp, err := client.GetDepositAddress(&bypayment.DepositAddressRequest{
    UserID:  12345,
    Network: "BSC",
    Coin:    "USDT",
})
```

## 快速开始

### 初始化客户端

```go
package main

import (
    "fmt"
    "github.com/bypayment8/bypayment-go-sdk"
)

func main() {
    // 创建客户端
    client := bypayment.NewClient(
        "https://api.example.com",  // API 服务器地址
        "your-api-key",              // 商户 API Key
        "your-secret-key",           // 商户 Secret Key
    )
    
    // 使用客户端调用 API...
}
```

## API 接口

### 币种网络配置

#### 获取币种网络配置列表

```go
// 获取所有币种的网络配置
networks, err := client.GetNetworksByCoin("")
if err != nil {
    log.Fatal(err)
}

// 获取指定币种的网络配置
networks, err := client.GetNetworksByCoin("USDT")
if err != nil {
    log.Fatal(err)
}

for _, network := range networks {
    fmt.Printf("Network: %s, Coin: %s, Name: %s\n", 
        network.Network, network.Coin, network.Name)
}
```

### 充值相关接口

#### 获取充值地址

```go
req := &bypayment.DepositAddressRequest{
    UserID:  12345,
    Network: "BSC",
    Coin:    "USDT",
}

resp, err := client.GetDepositAddress(req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Deposit Address: %s\n", resp.DepositAddress)
fmt.Printf("Expired Timestamp: %d\n", resp.ExpiredTimestamp)
```

#### 查询充值记录列表

```go
userID := int64(12345)
page := 1
pageSize := 20

result, err := client.GetDepositList(&userID, page, pageSize)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total: %d, Page: %d, PageSize: %d\n", 
    result.Total, result.Page, result.PageSize)

if deposits, ok := result.List.([]*bypayment.DepositListResponse); ok {
    for _, deposit := range deposits {
        fmt.Printf("ID: %s, Amount: %s, Status: %d\n", 
            deposit.ID, deposit.DepositAmount, deposit.Status)
    }
}
```

#### 查询充值记录详情

```go
detail, err := client.GetDepositDetail("order-id-123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Deposit Detail: ID=%s, Amount=%s, Status=%d\n", 
    detail.ID, detail.DepositAmount, detail.Status)
```

### 提现相关接口

#### 发起提现

```go
req := &bypayment.WithdrawRequest{
    MerchantOrderNo: "merchant-order-123",
    UserID:          12345,
    Network:         "BSC",
    Coin:            "USDT",
    Address:         "0x1234567890123456789012345678901234567890",
    Amount:          "100.00",
}

resp, err := client.CreateWithdraw(req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Withdraw Order ID: %s\n", resp.ID)
```

#### 查询提现记录列表

```go
userID := int64(12345)
page := 1
pageSize := 20

result, err := client.GetWithdrawList(&userID, page, pageSize)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total: %d, Page: %d, PageSize: %d\n", 
    result.Total, result.Page, result.PageSize)

if withdraws, ok := result.List.([]*bypayment.WithdrawListResponse); ok {
    for _, withdraw := range withdraws {
        fmt.Printf("ID: %s, Amount: %s, Status: %d\n", 
            withdraw.ID, withdraw.TransferAmount, withdraw.Status)
    }
}
```

#### 查询提现订单详情

```go
detail, err := client.GetWithdrawDetail("order-id-123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Withdraw Detail: ID=%s, Amount=%s, Status=%d\n", 
    detail.ID, detail.TransferAmount, detail.Status)
```

## 高级配置

### 自定义 HTTP 客户端

```go
client := bypayment.NewClient(
    "https://api.example.com",
    "your-api-key",
    "your-secret-key",
)

// 设置自定义 HTTP 客户端（例如设置超时时间）
customClient := &http.Client{
    Timeout: 60 * time.Second,
}
client.SetHTTPClient(customClient)
```

## 回调处理

SDK 提供了完整的回调处理功能，包括签名验证和响应生成。

### 基本用法

```go
import (
    "encoding/json"
    "net/http"
    "github.com/bypayment8/bypayment-go-sdk"
)

// HTTP Handler 示例
func callbackHandler(w http.ResponseWriter, r *http.Request) {
    // 1. 解析回调请求
    var callbackData map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&callbackData); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // 2. 验证签名
    secretKey := "your-secret-key"
    if !bypayment.VerifyCallbackSignature(callbackData, secretKey) {
        resp := bypayment.NewCallbackErrorResponse(40001, "签名验证失败")
        bypayment.WriteCallbackResponse(w, resp)
        return
    }

    // 3. 解析回调数据
    req, err := bypayment.ParseCallbackRequest(r)
    if err != nil {
        resp := bypayment.NewCallbackErrorResponse(40002, "解析回调数据失败")
        bypayment.WriteCallbackResponse(w, resp)
        return
    }

    // 4. 处理业务逻辑
    // ... 处理订单状态变更 ...

    // 5. 返回成功响应
    resp := bypayment.NewCallbackSuccessResponse("success")
    bypayment.WriteCallbackResponse(w, resp)
}
```

### 便捷函数

使用 `VerifyAndRespond` 可以简化回调处理：

```go
func callbackHandler(w http.ResponseWriter, r *http.Request) {
    var callbackData map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&callbackData); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    secretKey := "your-secret-key"
    
    // 验证并处理回调
    verified, resp := bypayment.VerifyAndRespond(
        callbackData,
        secretKey,
        func(req *bypayment.CallbackRequest) (code int, message string) {
            // 处理回调业务逻辑
            fmt.Printf("Order: %s, Status: %s\n", req.OrderNo, req.Status)
            
            // 返回处理结果
            return 200, "success"
        },
    )

    if !verified {
        // 签名验证失败
        bypayment.WriteCallbackResponse(w, resp)
        return
    }

    // 写入响应
    bypayment.WriteCallbackResponse(w, resp)
}
```

### 回调数据结构

```go
type CallbackRequest struct {
    OrderNo      string `json:"orderNo"`      // 订单号
    MerchantID   int64  `json:"merchantId"`   // 商户ID
    UserID       int64  `json:"userId"`       // 用户ID
    Coin         string `json:"coin"`         // 币种
    Network      string `json:"network"`      // 网络
    TxID         string `json:"txId"`         // 交易ID
    FromAddress  string `json:"fromAddress"`  // 发送地址
    ToAddress    string `json:"toAddress"`    // 接收地址
    Amount       string `json:"amount"`       // 金额
    ActualAmount string `json:"actualAmount"` // 实际金额
    Fee          string `json:"fee"`          // 手续费
    Status       string `json:"status"`       // 订单状态: confirming, confirmed, successful, failed
    OrderType    string `json:"orderType"`    // 订单类型: deposit, withdraw
    ConfirmedAt  *int64 `json:"confirmedAt,omitempty"` // 确认时间戳（秒）
    Sign         string `json:"sign"`         // 签名
    SignType     string `json:"signType"`    // 签名类型
}
```

## 错误处理

所有 API 方法都会返回错误。错误可能包括：

- 网络错误
- API 错误（错误码非 200）
- 签名验证失败
- 参数错误

```go
result, err := client.GetDepositList(nil, 1, 20)
if err != nil {
    // 处理错误
    fmt.Printf("Error: %v\n", err)
    return
}
```

API 响应格式统一为：
```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

当 `code` 不为 200 时，表示请求失败，错误信息在 `message` 字段中。

## 签名机制

SDK 会自动处理所有签名相关的逻辑：

1. 自动生成时间戳（timestamp）和随机数（nonce）
2. 将所有参数（包括 apiKey、timestamp、nonce）按 ASCII 码排序
3. 拼接成 `a=b&c=d` 格式
4. 使用 Secret Key 进行 MD5 签名
5. 自动添加必要的 HTTP Headers

## 安全特性

- **签名验证**：所有请求都使用 MD5 签名进行验证
- **防重放攻击**：使用 timestamp + nonce 机制防止重放攻击
- **时间戳验证**：请求在 5 分钟内有效
- **频率限制**：服务器端对每个 API Key 进行频率限制（默认 200 次/分钟）

## 接口列表

### 币种网络配置
- `GetNetworksByCoin(coin string)` - 获取币种网络配置列表

### 充值相关
- `GetDepositAddress(req *DepositAddressRequest)` - 获取充值地址
- `GetDepositList(userID *int64, page, pageSize int)` - 查询充值记录列表
- `GetDepositDetail(orderID string)` - 查询充值记录详情

### 提现相关
- `CreateWithdraw(req *WithdrawRequest)` - 发起提现
- `GetWithdrawList(userID *int64, page, pageSize int)` - 查询提现记录列表
- `GetWithdrawDetail(orderID string)` - 查询提现订单详情

## 许可证

MIT License

## 支持

如有问题或建议，请联系技术支持。

