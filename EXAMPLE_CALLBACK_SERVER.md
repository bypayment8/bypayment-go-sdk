# 回调服务器使用说明

## 方式一：使用独立可执行文件（推荐）

### 1. 编译可执行文件

```bash
cd bypayment-go-sdk
go build -o callback_server ./cmd/callback_server/main.go
```

### 2. 运行回调服务器

```bash
./callback_server -port=9999 -secret=your-secret-key
```

参数说明：
- `-port`: 回调服务器端口（默认: 9999）
- `-secret`: 商户 Secret Key（必填）

### 3. 配置服务器回调 URL

在服务器配置中将回调 URL 设置为：
```
http://127.0.0.1:9999/callback
```

### 4. 停止服务器

按 `Ctrl+C` 停止服务器，服务器会自动打印收到的所有回调统计信息。

## 方式二：使用测试函数

### 运行测试

```bash
cd bypayment-go-sdk
go test -v -run TestCallbackServer -timeout 10m
```

测试会在超时后自动停止，并打印收到的回调统计信息。

## 功能特性

1. **自动验证签名**：使用 SDK 的 `VerifyAndRespond` 函数自动验证回调签名
2. **实时打印回调**：收到回调时立即打印详细信息
3. **状态监控**：每30秒打印一次等待状态和已收到的回调数量
4. **最终统计**：停止时打印所有收到的回调请求的详细信息
5. **正确响应**：自动返回正确的回调响应格式

## 回调数据格式

回调服务器会打印以下信息：
- 请求方法和 URL
- 请求头
- 完整的回调数据（JSON 格式）
- 订单号、商户ID、用户ID
- 币种、网络、交易ID
- 金额、手续费、状态
- 订单类型、确认时间等

## 响应格式

服务器会自动返回标准的回调响应：

成功响应：
```json
{
  "code": 200,
  "message": "success"
}
```

失败响应（签名验证失败）：
```json
{
  "code": 40001,
  "message": "签名验证失败"
}
```

