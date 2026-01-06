package bypayment

import "github.com/shopspring/decimal"

// NetworkConfigResponse 网络配置响应
type NetworkConfigResponse struct {
	Network                 string `json:"network"`                      // 网络名称
	Coin                    string `json:"coin"`                         // 币种
	CoinName                string `json:"coinName"`                     // 币种名称
	NetworkName             string `json:"networkName"`                  // 网络名称
	DepositDust             string `json:"depositDust"`                  // 充值最小金额（粉尘值）
	DepositFee              string `json:"depositFee"`                   // 充值手续费
	WithdrawFee             string `json:"withdrawFee"`                  // 提现手续费
	WithdrawMin             string `json:"withdrawMin"`                  // 最小提现金额
	WithdrawIntegerMultiple string `json:"withdrawIntegerMultiple"`      // 提现金额整数倍数
	WithdrawMax             string `json:"withdrawMax"`                  // 最大提现金额
	ContractAddress         string `json:"contractAddress,omitempty"`    // 合约地址（仅ERC20等代币）
	ContractAddressURL      string `json:"contractAddressUrl,omitempty"` // 合约地址浏览器链接
	BlockInterval           int    `json:"blockInterval,omitempty"`      // 区块间隔（秒）
	BlockURL                string `json:"blockUrl,omitempty"`           // 区块浏览器链接模板
	TxURL                   string `json:"txUrl,omitempty"`              // 交易浏览器链接模板
	AddressURL              string `json:"addressUrl,omitempty"`         // 地址浏览器链接模板
	ConfirmTimes            int    `json:"confirmTimes"`                 // 确认次数
	UnlockConfirmTimes      int    `json:"unlockConfirmTimes"`           // 解锁确认次数
	InsertTime              int64  `json:"insertTime"`                   // 创建时间戳（秒）
	UpdateTime              int64  `json:"updateTime"`                   // 更新时间戳（秒）
}

// DepositAddressRequest 获取充值地址请求
type DepositAddressRequest struct {
	UserID   int64  `json:"userId"`             // 用户ID
	Username string `json:"username,omitempty"` // 用户名（可选，用于控制台展示）
	Network  string `json:"network"`            // 网络
	Coin     string `json:"coin"`               // 币种
}

// DepositAddressResponse 充值地址响应
type DepositAddressResponse struct {
	DepositAddress   string `json:"depositAddress"`   // 充值地址
	ExpiredTimestamp int64  `json:"expiredTimestamp"` // 过期时间戳（秒）
	Tip              string `json:"tip,omitempty"`    // 提示信息，如"请在30分钟内完成充值"
}

// DepositListResponse 充值记录列表响应
type DepositListResponse struct {
	ID                   string          `json:"id"`                             // 订单ID
	UserID               int64           `json:"userId"`                         // 用户ID
	Coin                 string          `json:"coin"`                           // 币种
	Network              string          `json:"network"`                        // 网络
	DepositAmount        decimal.Decimal `json:"depositAmount"`                  // 充值金额
	TxID                 string          `json:"txId,omitempty"`                 // 交易哈希
	TxURL                string          `json:"txUrl,omitempty"`                // 交易浏览器链接
	Address              string          `json:"address"`                        // 充值地址
	AddressURL           string          `json:"addressUrl,omitempty"`           // 地址浏览器链接
	CurConfirmTimes      int             `json:"curConfirmTimes"`                // 当前确认次数
	ConfirmTimes         int             `json:"confirmTimes"`                   // 所需确认次数
	UnlockConfirm        int             `json:"unlockConfirm"`                  // 解锁确认次数
	EstimatedArrivalTime uint64          `json:"estimatedArrivalTime,omitempty"` // 预计到账时间戳（秒）
	EstimatedUnlockTime  uint64          `json:"estimatedUnlockTime,omitempty"`  // 预计解锁时间戳（秒）
	Status               string          `json:"status"`                         // 订单状态
	InsertTime           uint64          `json:"insertTime,omitempty"`           // 创建时间戳（秒）
}

// DepositDetailResponse 充值记录详情响应
type DepositDetailResponse struct {
	DepositListResponse
}

// WithdrawRequest 提现请求
type WithdrawRequest struct {
	MerchantOrderNo string `json:"merchantOrderNo"`    // 商户订单号
	UserID          int64  `json:"userId"`             // 用户ID
	Username        string `json:"username,omitempty"` // 用户名（可选，用于控制台展示）
	Network         string `json:"network"`            // 网络
	Coin            string `json:"coin"`               // 币种
	Address         string `json:"address"`            // 提现地址
	Amount          string `json:"amount"`             // 提现金额
}

// WithdrawResponse 提现响应
type WithdrawResponse struct {
	ID string `json:"id"` // 订单ID
}

// WithdrawListResponse 提现记录列表响应
type WithdrawListResponse struct {
	ID                   string          `json:"id"`                             // 订单ID
	MerchantOrderNo      string          `json:"merchantOrderNo"`                // 商户订单号
	TransferAmount       decimal.Decimal `json:"transferAmount"`                 // 转账金额
	TransactionFee       decimal.Decimal `json:"transactionFee"`                 // 交易手续费
	ApplyTime            int64           `json:"applyTime"`                      // 申请时间戳（秒）
	Status               string          `json:"status"`                         // 订单状态
	AddressURL           string          `json:"addressUrl,omitempty"`           // 地址浏览器链接
	UserID               int64           `json:"userId"`                         // 用户ID
	Network              string          `json:"network"`                        // 网络
	Coin                 string          `json:"coin"`                           // 币种
	Address              string          `json:"address"`                        // 提现地址
	TxID                 string          `json:"txId,omitempty"`                 // 交易哈希
	TxURL                string          `json:"txUrl,omitempty"`                // 交易浏览器链接
	CurConfirmTimes      int             `json:"curConfirmTimes"`                // 当前确认次数
	ConfirmTimes         int             `json:"confirmTimes"`                   // 所需确认次数
	EstimatedArrivalTime uint64          `json:"estimatedArrivalTime,omitempty"` // 预计到账时间戳（秒）
}

// WithdrawDetailResponse 提现订单详情响应
type WithdrawDetailResponse struct {
	WithdrawListResponse
}

// PageResponse 分页响应
type PageResponse struct {
	List     interface{} `json:"list"`     // 数据列表
	Total    int64       `json:"total"`    // 总记录数
	Page     int         `json:"page"`     // 当前页码
	PageSize int         `json:"pageSize"` // 每页大小
}
