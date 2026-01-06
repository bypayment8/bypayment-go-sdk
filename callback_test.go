package bypayment_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bypayment8/bypayment-go-sdk"
)

// TestCallbackServer 测试回调服务器
// 运行方式: go test -v -run TestCallbackServer -timeout 10m
func TestCallbackServer(t *testing.T) {
	// 创建回调服务器
	server := bypayment.NewCallbackServer("9999", testSecretKey)

	// 启动服务器
	if err := server.Start(); err != nil {
		t.Fatalf("启动回调服务器失败: %v", err)
	}
	defer server.Stop()

	// 打印回调 URL
	fmt.Printf("========================================\n")
	fmt.Printf("回调服务器已启动\n")
	fmt.Printf("回调 URL: %s\n", server.GetCallbackURL())
	fmt.Printf("等待回调请求...\n")
	fmt.Printf("请在服务器配置中将回调 URL 设置为: %s\n", server.GetCallbackURL())
	fmt.Printf("测试将在超时后自动停止\n")
	fmt.Printf("========================================\n\n")

	// 每30秒打印一次状态
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// 定期打印状态
	go func() {
		for range ticker.C {
			received := server.GetReceived()
			fmt.Printf("[%s] 等待回调中... 已收到 %d 个回调\n", time.Now().Format("15:04:05"), len(received))
		}
	}()

	// 等待测试超时（测试超时时间由 -timeout 参数控制，例如 -timeout 10m）
	// 这里等待足够长的时间，让测试可以持续运行
	<-time.After(30 * time.Minute)

	// 最终检查
	received := server.GetReceived()
	if len(received) > 0 {
		fmt.Printf("\n=== 最终统计 ===\n")
		fmt.Printf("共收到 %d 个回调请求\n", len(received))
		for i, req := range received {
			fmt.Printf("\n回调 %d:\n", i+1)
			fmt.Printf("  订单号: %s\n", req.OrderNo)
			fmt.Printf("  用户ID: %d\n", req.UserID)
			fmt.Printf("  币种: %s\n", req.Coin)
			fmt.Printf("  网络: %s\n", req.Network)
			fmt.Printf("  状态: %s\n", req.Status)
			fmt.Printf("  类型: %s\n", req.OrderType)
			fmt.Printf("  金额: %s\n", req.Amount)
			fmt.Printf("  手续费: %s\n", req.Fee)
			fmt.Printf("  交易ID: %v\n", req.TxID)
			fmt.Printf("  发送地址: %v\n", req.FromAddress)
			fmt.Printf("  接收地址: %v\n", req.ToAddress)
		}
		fmt.Printf("===============\n")
	} else {
		fmt.Printf("\n未收到任何回调请求\n")
		fmt.Printf("提示: 确保服务器已配置回调 URL 为: %s\n", server.GetCallbackURL())
	}
}
