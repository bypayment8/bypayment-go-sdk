package bypayment

import (
	"net/http"
	"strconv"
)

// GetDepositList 查询充值记录列表
// userID: 用户ID（可选，用于查询特定用户的记录）
// page: 页码（可选，默认为1）
// pageSize: 每页数量（可选，默认为20）
func (c *Client) GetDepositList(userID *int64, page, pageSize int) (*PageResponse, error) {
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

	resp, err := c.doRequest(http.MethodGet, "/api/payment/deposits/list", params, nil)
	if err != nil {
		return nil, err
	}

	// 解析分页响应
	var pageData struct {
		List       []*DepositListResponse `json:"list"`
		Total      int64                  `json:"total"`
		Page       int                    `json:"page"`
		PageSize   int                    `json:"pageSize"`
		TotalPages int                    `json:"totalPages"`
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

// GetDepositDetail 查询充值记录详情
// orderID: 订单ID
func (c *Client) GetDepositDetail(orderID string) (*DepositDetailResponse, error) {
	params := make(map[string]interface{})
	params["id"] = orderID

	resp, err := c.doRequest(http.MethodGet, "/api/payment/deposits/detail", params, nil)
	if err != nil {
		return nil, err
	}

	var result DepositDetailResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetDepositAddress 获取用户的充值地址
// req: 获取充值地址请求
func (c *Client) GetDepositAddress(req *DepositAddressRequest) (*DepositAddressResponse, error) {
	body := map[string]interface{}{
		"userId":  req.UserID,
		"network": req.Network,
		"coin":    req.Coin,
	}
	// 如果提供了 username，添加到请求体中
	if req.Username != "" {
		body["username"] = req.Username
	}

	resp, err := c.doRequest(http.MethodPost, "/api/payment/deposit_addresses", nil, body)
	if err != nil {
		return nil, err
	}

	var result DepositAddressResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
