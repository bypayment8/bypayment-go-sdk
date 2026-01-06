package bypayment

import (
	"net/http"
	"strconv"
)

// GetWithdrawList 查询提现记录列表
// userID: 用户ID（可选，用于查询特定用户的记录）
// page: 页码（可选，默认为1）
// pageSize: 每页数量（可选，默认为20）
func (c *Client) GetWithdrawList(userID *int64, page, pageSize int) (*PageResponse, error) {
	params := make(map[string]interface{})
	if userID != nil {
		params["userId"] = strconv.FormatInt(*userID, 10)
	}
	if page > 0 {
		params["page"] = page
	}
	if pageSize > 0 {
		params["pageSize"] = pageSize
	}

	resp, err := c.doRequest(http.MethodGet, "/api/payment/withdraws/list", params, nil)
	if err != nil {
		return nil, err
	}

	// 解析分页响应
	var pageData struct {
		List       []*WithdrawListResponse `json:"list"`
		Total      int64                   `json:"total"`
		Page       int                     `json:"page"`
		PageSize   int                     `json:"pageSize"`
		TotalPages int                     `json:"totalPages"`
	}
	if err := parseResponse(resp, &pageData); err != nil {
		return nil, err
	}

	return &PageResponse{
		List:     pageData.List,
		Total:    pageData.Total,
		Page:     pageData.Page,
		PageSize: pageData.PageSize,
	}, nil
}

// GetWithdrawDetail 查询提现订单详情
// orderID: 订单ID
func (c *Client) GetWithdrawDetail(orderID string) (*WithdrawDetailResponse, error) {
	params := make(map[string]interface{})
	params["id"] = orderID

	resp, err := c.doRequest(http.MethodGet, "/api/payment/withdraws/detail", params, nil)
	if err != nil {
		return nil, err
	}

	var result WithdrawDetailResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateWithdraw 发起提现
// req: 提现请求
func (c *Client) CreateWithdraw(req *WithdrawRequest) (*WithdrawResponse, error) {
	body := map[string]interface{}{
		"merchantOrderNo": req.MerchantOrderNo,
		"userId":          req.UserID,
		"network":         req.Network,
		"coin":            req.Coin,
		"address":         req.Address,
		"amount":          req.Amount,
	}
	// 如果提供了 username，添加到请求体中
	if req.Username != "" {
		body["username"] = req.Username
	}

	resp, err := c.doRequest(http.MethodPost, "/api/payment/withdraws", nil, body)
	if err != nil {
		return nil, err
	}

	var result WithdrawResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
