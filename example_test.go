package bypayment_test

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/bypayment8/bypayment-go-sdk"
)

// 测试配置常量
const (
	testBaseURL   = "https://api.bypayment.io"
	testAPIKey    = "your-api-key"
	testSecretKey = "your-secret-key"
)

// 代理配置常量
const (
	// ClashProxyURL Clash 代理地址（如果不需要代理，设置为空字符串 ""）
	//ClashProxyURL = "http://127.0.0.1:7897"
	ClashProxyURL = "" // 不使用代理时设置为空字符串
)

// newTestClient 创建测试客户端并配置代理（如果需要）
func newTestClient() *bypayment.Client {
	client := bypayment.NewClient(testBaseURL, testAPIKey, testSecretKey)

	// 配置使用 Clash 代理（如果设置了代理地址）
	// 方法1: 通过环境变量设置（推荐）
	// export HTTP_PROXY=http://127.0.0.1:7890
	// export HTTPS_PROXY=http://127.0.0.1:7890
	// 然后运行测试: go test

	// 方法2: 在代码中配置代理
	if ClashProxyURL != "" {
		proxyURL, err := url.Parse(ClashProxyURL)
		if err == nil {
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
			client.SetHTTPClient(&http.Client{
				Transport: transport,
				Timeout:   30 * time.Second,
			})
		}
	}

	return client
}

func TestClient_GetNetworksByCoin(t *testing.T) {
	client := newTestClient()

	// 获取所有币种的网络配置
	networks, err := client.GetNetworksByCoin("")
	if err != nil {
		t.Fatalf("获取网络配置失败: %v", err)
	}

	for _, network := range networks {
		fmt.Printf("Network: %s, Coin: %s deposit_fee: %v withdraw_fee: %v \n", network.Network, network.Coin, network.DepositFee, network.WithdrawFee)
	}
}

func TestClient_GetDepositAddress(t *testing.T) {
	client := newTestClient()

	// 获取充值地址
	req := &bypayment.DepositAddressRequest{
		UserID:   12345,
		Username: "test2",
		Network:  "TRON",
		Coin:     "USDT",
	}

	resp, err := client.GetDepositAddress(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deposit Address: %s\n", resp.DepositAddress)
}

func TestClient_GetDepositList(t *testing.T) {
	client := newTestClient()

	// 查询充值记录列表
	//userID := int64(12345)
	page := 1
	pageSize := 20

	result, err := client.GetDepositList(nil, page, pageSize)
	if err != nil {
		t.Fatalf("获取充值列表失败: %v", err)
	}

	fmt.Printf("Total: %d, Page: %d, PageSize: %d\n", result.Total, result.Page, result.PageSize)

	if deposits, ok := result.List.([]*bypayment.DepositListResponse); ok {
		for _, deposit := range deposits {
			fmt.Printf("ID: %s, UserID: %d, Coin: %s, Network: %s, Amount: %v, Status: %s\n",
				deposit.ID, deposit.UserID, deposit.Coin, deposit.Network, deposit.DepositAmount, deposit.Status)
		}
	}
}

func TestClient_GetDepositDetail(t *testing.T) {
	client := newTestClient()

	// 先获取充值列表，使用第一个订单的ID来查询详情
	//userID := int64(12345)
	result, err := client.GetDepositList(nil, 1, 1)
	if err != nil {
		t.Fatalf("获取充值列表失败: %v", err)
	}

	// 如果没有订单，跳过测试
	if result.Total == 0 {
		t.Skip("没有充值订单，跳过详情测试")
		return
	}

	// 获取第一个订单的ID
	deposits, ok := result.List.([]*bypayment.DepositListResponse)
	if !ok || len(deposits) == 0 {
		t.Skip("没有充值订单，跳过详情测试")
		return
	}

	orderID := deposits[0].ID

	// 查询充值记录详情
	detail, err := client.GetDepositDetail(orderID)
	if err != nil {
		t.Fatalf("获取充值详情失败: %v", err)
	}

	fmt.Printf("Deposit Detail:\n")
	fmt.Printf("  ID: %s\n", detail.ID)
	fmt.Printf("  UserID: %d\n", detail.UserID)
	fmt.Printf("  Coin: %s\n", detail.Coin)
	fmt.Printf("  Network: %s\n", detail.Network)
	fmt.Printf("  Amount: %v\n", detail.DepositAmount)
	fmt.Printf("  Status: %s\n", detail.Status)
	fmt.Printf("  TxID: %s\n", detail.TxID)
	fmt.Printf("  Address: %s\n", detail.Address)
	fmt.Printf("  CurConfirmTimes: %d\n", detail.CurConfirmTimes)
	fmt.Printf("  ConfirmTimes: %d\n", detail.ConfirmTimes)
}

func TestClient_GetWithdrawList(t *testing.T) {
	client := newTestClient()

	// 查询提现记录列表
	//userID := int64(12345)
	page := 1
	pageSize := 20

	result, err := client.GetWithdrawList(nil, page, pageSize)
	if err != nil {
		t.Fatalf("获取提现列表失败: %v", err)
	}

	fmt.Printf("Total: %d, Page: %d, PageSize: %d\n", result.Total, result.Page, result.PageSize)

	if withdraws, ok := result.List.([]*bypayment.WithdrawListResponse); ok {
		for _, withdraw := range withdraws {
			fmt.Printf("ID: %s, MerchantOrderNo: %s, UserID: %d, Coin: %s, Network: %s, Amount: %v, Fee: %v, Status: %s\n",
				withdraw.ID, withdraw.MerchantOrderNo, withdraw.UserID, withdraw.Coin, withdraw.Network, withdraw.TransferAmount, withdraw.TransactionFee, withdraw.Status)
		}
	}
}

func TestClient_GetWithdrawDetail(t *testing.T) {
	client := newTestClient()

	// 先获取提现列表，使用第一个订单的ID来查询详情
	//userID := int64(12345)
	result, err := client.GetWithdrawList(nil, 1, 1)
	if err != nil {
		t.Fatalf("获取提现列表失败: %v", err)
	}

	// 如果没有订单，跳过测试
	if result.Total == 0 {
		t.Skip("没有提现订单，跳过详情测试")
		return
	}

	// 获取第一个订单的ID
	withdraws, ok := result.List.([]*bypayment.WithdrawListResponse)
	if !ok || len(withdraws) == 0 {
		t.Skip("没有提现订单，跳过详情测试")
		return
	}

	orderID := withdraws[0].ID

	// 查询提现订单详情
	detail, err := client.GetWithdrawDetail(orderID)
	if err != nil {
		t.Fatalf("获取提现详情失败: %v", err)
	}

	fmt.Printf("Withdraw Detail:\n")
	fmt.Printf("  ID: %s\n", detail.ID)
	fmt.Printf("  MerchantOrderNo: %s\n", detail.MerchantOrderNo)
	fmt.Printf("  UserID: %d\n", detail.UserID)
	fmt.Printf("  Coin: %s\n", detail.Coin)
	fmt.Printf("  Network: %s\n", detail.Network)
	fmt.Printf("  TransferAmount: %v\n", detail.TransferAmount)
	fmt.Printf("  TransactionFee: %v\n", detail.TransactionFee)
	fmt.Printf("  Status: %s\n", detail.Status)
	fmt.Printf("  TxID: %s\n", detail.TxID)
	fmt.Printf("  Address: %s\n", detail.Address)
	fmt.Printf("  CurConfirmTimes: %d\n", detail.CurConfirmTimes)
	fmt.Printf("  ConfirmTimes: %d\n", detail.ConfirmTimes)
}

func TestClient_CreateWithdraw(t *testing.T) {
	client := newTestClient()

	// 发起提现
	req := &bypayment.WithdrawRequest{
		MerchantOrderNo: "1547043348858710026",
		UserID:          12345,
		Network:         "TRON",
		Coin:            "USDT",
		Address:         "to address",
		Amount:          "10",
	}

	resp, err := client.CreateWithdraw(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Withdraw Order ID: %s\n", resp.ID)
}
