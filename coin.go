package bypayment

import "net/http"

// GetNetworksByCoin 获取币种网络配置列表
// coin: 币种（可选，如果不传则返回所有，如果传了则过滤指定币种）
func (c *Client) GetNetworksByCoin(coin string) ([]*NetworkConfigResponse, error) {
	params := make(map[string]interface{})
	if coin != "" {
		params["coin"] = coin
	}

	resp, err := c.doRequest(http.MethodGet, "/api/payment/coins/networks/list", params, nil)
	if err != nil {
		return nil, err
	}

	var result []*NetworkConfigResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}
