package bypayment

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Client 是 ByPayment SDK 的客户端
type Client struct {
	baseURL    string
	apiKey     string
	secretKey  string
	httpClient *http.Client
}

// NewClient 创建一个新的 ByPayment 客户端
// baseURL: API 服务器的基础 URL，例如 "https://api.example.com" 或 "http://localhost:8080"
// apiKey: 商户的 API Key
// secretKey: 商户的 Secret Key（用于签名）
func NewClient(baseURL, apiKey, secretKey string) *Client {
	return &Client{
		baseURL:   strings.TrimSuffix(baseURL, "/"),
		apiKey:    apiKey,
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetHTTPClient 设置自定义的 HTTP 客户端
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// generateNonce 生成随机数（用于防重放）
func generateNonce() string {
	return uuid.New().String()
}

// generateTimestamp 生成当前时间戳（秒级）
func generateTimestamp() int64 {
	return time.Now().Unix()
}

// MD5Sign 生成 MD5 签名（公开函数，供回调验证使用）
// 1. 将发送或接收到的所有参数按照参数名ASCII码从小到大排序（a-z），sign、signType、和空值不参与签名！
// 2. 将排序后的参数拼接成URL键值对的格式，例如 a=b&c=d&e=f，参数值不要进行url编码。
// 3. 再将拼接好的字符串与商户密钥KEY进行MD5加密得出sign签名参数，sign = md5 ( a=b&c=d&e=f + KEY )
// 4. md5结果为小写。
func MD5Sign(params map[string]interface{}, key string) string {
	// 1. 过滤并排序参数
	keys := make([]string, 0, len(params))
	for k, v := range params {
		// sign、signType、和空值不参与签名
		if k == "sign" || k == "signType" {
			continue
		}
		// 空值不参与签名
		if v == nil || v == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. 拼接成 URL 键值对格式
	var builder strings.Builder
	for i, k := range keys {
		if i > 0 {
			builder.WriteString("&")
		}
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(formatValueForSign(params[k]))
	}
	queryString := builder.String()

	// 3. 拼接密钥并 MD5 加密
	signString := queryString + key
	hash := md5.Sum([]byte(signString))
	return fmt.Sprintf("%x", hash)
}

// formatValueForSign 格式化参数值用于签名计算
// JSON 解析时数字会被解析为 float64，需要正确处理避免科学计数法
func formatValueForSign(v interface{}) string {
	switch val := v.(type) {
	case float64:
		if val == math.Trunc(val) {
			return strconv.FormatInt(int64(val), 10)
		}
		return strconv.FormatFloat(val, 'f', -1, 64)
	case float32:
		if val == float32(math.Trunc(float64(val))) {
			return strconv.FormatInt(int64(val), 10)
		}
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// doRequest 执行 HTTP 请求（带签名）
func (c *Client) doRequest(method, path string, params map[string]interface{}, body interface{}) (*http.Response, error) {
	// 生成时间戳和随机数
	timestamp := generateTimestamp()
	nonce := generateNonce()

	// 构建完整 URL
	requestURL := c.baseURL + path

	var req *http.Request
	var err error

	var sign string

	if method == http.MethodGet {
		// GET 请求：业务参数放在 query string，apiKey/timestamp/nonce 在 Header 中
		if params == nil {
			params = make(map[string]interface{})
		}

		// 创建签名参数（包含业务参数和认证参数）
		signParams := make(map[string]interface{})
		// 复制业务参数到签名参数
		for k, v := range params {
			signParams[k] = v
		}
		// 添加 apiKey、timestamp、nonce 到签名参数（这些值会参与签名计算，但不在 query 中）
		signParams["apiKey"] = c.apiKey
		signParams["timestamp"] = strconv.FormatInt(timestamp, 10)
		signParams["nonce"] = nonce

		// 生成签名
		sign = MD5Sign(signParams, c.secretKey)

		// 构建 query string（只包含业务参数，不包含 apiKey/timestamp/nonce）
		queryValues := url.Values{}
		for k, v := range params {
			if v != nil && v != "" {
				queryValues.Set(k, fmt.Sprintf("%v", v))
			}
		}
		if len(queryValues) > 0 {
			requestURL += "?" + queryValues.Encode()
		}

		req, err = http.NewRequest(method, requestURL, nil)
	} else {
		// POST 请求：参数放在 body
		if body == nil {
			body = make(map[string]interface{})
		}
		bodyMap, ok := body.(map[string]interface{})
		if !ok {
			// 如果不是 map，尝试转换为 map
			bodyBytes, _ := json.Marshal(body)
			json.Unmarshal(bodyBytes, &bodyMap)
		}

		// 添加 apiKey、timestamp、nonce 到签名参数
		bodyMap["apiKey"] = c.apiKey
		bodyMap["timestamp"] = strconv.FormatInt(timestamp, 10)
		bodyMap["nonce"] = nonce

		// 生成签名
		sign = MD5Sign(bodyMap, c.secretKey)

		// 序列化 body
		bodyBytes, err := json.Marshal(bodyMap)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}

		req, err = http.NewRequest(method, requestURL, bytes.NewBuffer(bodyBytes))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置 Headers
	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("X-Timestamp", strconv.FormatInt(timestamp, 10))
	req.Header.Set("X-Nonce", nonce)
	req.Header.Set("X-Sign-Type", "MD5")
	req.Header.Set("X-Sign", sign)

	// 执行请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	return resp, nil
}

// APIResponse 是 API 的统一响应格式
type APIResponse struct {
	RequestID string      `json:"requestID,omitempty"`
	Code      int         `json:"code"`    // 200 表示成功，其他为错误码
	Message   string      `json:"message"` // 提示信息
	Data      interface{} `json:"data"`    // 数据负载
}

// parseResponse 解析响应
func parseResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		var apiResp APIResponse
		if err := json.Unmarshal(bodyBytes, &apiResp); err == nil {
			return fmt.Errorf("API 错误 (code: %d): %s", apiResp.Code, apiResp.Message)
		}
		return fmt.Errorf("API 错误 (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	// 解析为统一响应格式
	var apiResp APIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查业务码
	if apiResp.Code != 200 {
		return fmt.Errorf("API 错误 (code: %d): %s", apiResp.Code, apiResp.Message)
	}

	// 提取 data 字段
	if apiResp.Data == nil {
		return nil
	}

	// 将 data 转换为目标类型
	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %w", err)
	}

	if err := json.Unmarshal(dataBytes, result); err != nil {
		return fmt.Errorf("解析数据失败: %w", err)
	}

	return nil
}
