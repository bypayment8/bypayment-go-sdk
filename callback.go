package bypayment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
)

// CallbackRequest 回调请求数据
type CallbackRequest struct {
	OrderNo                  string          `json:"orderNo"`                  // 订单号
	UserID                   *int64          `json:"userId"`                   // 用户ID
	Coin                     string          `json:"coin"`                     // 币种
	Network                  string          `json:"network"`                  // 网络
	TxID                     *string         `json:"txId"`                     // 交易ID
	FromAddress              *string         `json:"fromAddress"`              // 发送地址
	ToAddress                *string         `json:"toAddress"`                // 接收地址
	Amount                   decimal.Decimal `json:"amount"`                   // 金额
	Fee                      decimal.Decimal `json:"fee"`                      // 手续费
	Status                   string          `json:"status"`                   // 订单状态: confirming, confirmed, successful, failed
	OrderType                string          `json:"orderType"`                // 订单类型: deposit, withdraw
	BlockNumber              *uint64         `json:"blockNumber"`              // 区块高度
	BlockTimestamp           *uint64         `json:"blockTimestamp"`           // 区块时间戳（秒）
	ConfirmBlockNumber       *uint64         `json:"confirmBlockNumber"`       // 确认区块高度
	UnlockConfirmBlockNumber *uint64         `json:"unlockConfirmBlockNumber"` // 解锁确认区块高度
	ConfirmedTimestamp       *uint64         `json:"confirmedTimestamp"`       // 确认时间戳（秒）
	SuccessfulTimestamp      *uint64         `json:"successfulTimestamp"`      // 成功时间戳（秒）
	Sign                     string          `json:"sign"`                     // 签名
	SignType                 string          `json:"signType"`                 // 签名类型
}

// VerifyCallbackSignature 验证回调签名
// callbackData: 回调请求数据（map 格式）
// secretKey: 商户的 Secret Key
// 返回 true 表示签名验证通过，false 表示验证失败
func VerifyCallbackSignature(callbackData map[string]interface{}, secretKey string) bool {
	// 获取签名和签名类型
	sign, signOk := callbackData["sign"].(string)
	signType, signTypeOk := callbackData["signType"].(string)

	if !signOk || !signTypeOk {
		return false
	}

	// 目前只支持 MD5 签名
	if signType != "MD5" {
		return false
	}

	// 计算签名
	calculatedSign := MD5Sign(callbackData, secretKey)

	// 比较签名（不区分大小写）
	return strings.EqualFold(calculatedSign, sign)
}

// ParseCallbackRequest 解析回调请求
// 从 HTTP 请求中解析回调数据
func ParseCallbackRequest(r *http.Request) (*CallbackRequest, error) {
	var req CallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("解析回调请求失败: %w", err)
	}
	return &req, nil
}

// CallbackResponse 回调响应
type CallbackResponse struct {
	Code    int         `json:"code"`           // 响应码，200 表示成功
	Message string      `json:"message"`        // 响应消息
	Data    interface{} `json:"data,omitempty"` // 响应数据（可选）
}

// NewCallbackSuccessResponse 创建成功的回调响应
func NewCallbackSuccessResponse(message string) *CallbackResponse {
	if message == "" {
		message = "success"
	}
	return &CallbackResponse{
		Code:    200,
		Message: message,
	}
}

// NewCallbackErrorResponse 创建失败的回调响应
func NewCallbackErrorResponse(code int, message string) *CallbackResponse {
	if message == "" {
		message = "callback failed"
	}
	return &CallbackResponse{
		Code:    code,
		Message: message,
	}
}

// WriteCallbackResponse 写入回调响应到 HTTP ResponseWriter
// 成功情况（code == 200 && message == "success"）：返回文本 "success"
// 其他情况：返回对应的业务码和错误文本（文本格式）
func WriteCallbackResponse(w http.ResponseWriter, resp *CallbackResponse) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// 成功情况返回文本 "success"
	if resp.Code == 200 && resp.Message == "success" {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("success"))
		return err
	}

	// 其他情况返回业务码和错误文本
	statusCode := resp.Code
	if statusCode < 400 || statusCode > 599 {
		// 如果不是有效的 HTTP 状态码，使用 400
		statusCode = 400
	}
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(resp.Message))
	return err
}

// VerifyAndRespond 验证回调并返回响应
// 这是一个便捷函数，用于在 HTTP handler 中处理回调
// callbackData: 回调请求数据（map 格式）
// secretKey: 商户的 Secret Key
// handler: 回调处理函数，接收验证后的 CallbackRequest，返回响应码和消息
// 返回是否验证通过，以及响应数据
func VerifyAndRespond(
	callbackData map[string]interface{},
	secretKey string,
	handler func(*CallbackRequest) (code int, message string),
) (bool, *CallbackResponse) {
	// 验证签名
	if !VerifyCallbackSignature(callbackData, secretKey) {
		return false, NewCallbackErrorResponse(40001, "签名验证失败")
	}

	// 解析回调请求
	reqBytes, err := json.Marshal(callbackData)
	if err != nil {
		return false, NewCallbackErrorResponse(40002, "解析回调数据失败")
	}

	var req CallbackRequest
	if err := json.Unmarshal(reqBytes, &req); err != nil {
		return false, NewCallbackErrorResponse(40002, "解析回调数据失败")
	}

	// 调用处理函数
	code, message := handler(&req)

	// 构建响应
	if code == 200 {
		return true, NewCallbackSuccessResponse(message)
	}
	return true, NewCallbackErrorResponse(code, message)
}
