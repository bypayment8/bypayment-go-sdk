package bypayment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// CallbackServer 回调服务器
type CallbackServer struct {
	server      *http.Server
	port        string
	secretKey   string
	received    []CallbackRequest
	mu          sync.Mutex
	callbackURL string
}

// NewCallbackServer 创建回调服务器
func NewCallbackServer(port, secretKey string) *CallbackServer {
	return &CallbackServer{
		port:        port,
		secretKey:   secretKey,
		received:    make([]CallbackRequest, 0),
		callbackURL: fmt.Sprintf("http://0.0.0.0:%s/callback", port),
	}
}

// Start 启动服务器
func (s *CallbackServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", s.handleCallback)

	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: mux,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("回调服务器启动失败: %v\n", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)
	return nil
}

// Stop 停止服务器
func (s *CallbackServer) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

// GetCallbackURL 获取回调 URL
func (s *CallbackServer) GetCallbackURL() string {
	return s.callbackURL
}

// GetReceived 获取接收到的回调请求
func (s *CallbackServer) GetReceived() []CallbackRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.received
}

// ClearReceived 清空接收到的回调请求
func (s *CallbackServer) ClearReceived() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.received = make([]CallbackRequest, 0)
}

// handleCallback 处理回调请求
func (s *CallbackServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	var callbackData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&callbackData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 打印接收到的回调数据
	fmt.Printf("\n=== 收到回调请求 ===\n")
	fmt.Printf("时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL.String())

	callbackJSON, _ := json.MarshalIndent(callbackData, "", "  ")
	fmt.Printf("Body:\n%s\n", string(callbackJSON))
	fmt.Printf("===================\n\n")

	// 验证并处理回调
	verified, resp := VerifyAndRespond(
		callbackData,
		s.secretKey,
		func(req *CallbackRequest) (code int, message string) {
			// 保存接收到的回调请求
			s.mu.Lock()
			s.received = append(s.received, *req)
			s.mu.Unlock()

			// 打印回调信息
			fmt.Printf("回调处理:\n")
			fmt.Printf("  订单号: %s\n", req.OrderNo)
			fmt.Printf("  用户ID: %d\n", req.UserID)
			fmt.Printf("  币种: %s\n", req.Coin)
			fmt.Printf("  网络: %s\n", req.Network)
			fmt.Printf("  交易ID: %v\n", req.TxID)
			fmt.Printf("  金额: %s\n", req.Amount)
			fmt.Printf("  手续费: %s\n", req.Fee)
			fmt.Printf("  状态: %s\n", req.Status)
			fmt.Printf("  订单类型: %s\n", req.OrderType)
			fmt.Printf("\n")

			// 返回成功响应
			return 200, "success"
		},
	)

	if !verified {
		fmt.Printf("⚠️  签名验证失败\n")
	} else {
		fmt.Printf("✅ 签名验证通过\n")
	}

	// 写入响应
	if err := WriteCallbackResponse(w, resp); err != nil {
		fmt.Printf("写入响应失败: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Printf("回调响应: code=%d, message=%s\n\n", resp.Code, resp.Message)
}

// Run 运行服务器（阻塞，直到收到停止信号）
func (s *CallbackServer) Run() error {
	// 启动服务器
	if err := s.Start(); err != nil {
		return fmt.Errorf("启动回调服务器失败: %w", err)
	}

	// 打印回调 URL
	fmt.Printf("========================================\n")
	fmt.Printf("回调服务器已启动\n")
	fmt.Printf("回调 URL: %s\n", s.GetCallbackURL())
	fmt.Printf("等待回调请求...\n")
	fmt.Printf("请在服务器配置中将回调 URL 设置为: %s\n", s.GetCallbackURL())
	fmt.Printf("按 Ctrl+C 停止服务器\n")
	fmt.Printf("========================================\n\n")

	// 每30秒打印一次状态
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// 定期打印状态
	go func() {
		for range ticker.C {
			received := s.GetReceived()
			fmt.Printf("[%s] 等待回调中... 已收到 %d 个回调\n", time.Now().Format("15:04:05"), len(received))
		}
	}()

	// 等待停止信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Printf("\n正在停止服务器...\n")

	// 停止服务器
	if err := s.Stop(); err != nil {
		return fmt.Errorf("停止服务器失败: %w", err)
	}

	// 打印最终统计
	received := s.GetReceived()
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
		fmt.Printf("提示: 确保服务器已配置回调 URL 为: %s\n", s.GetCallbackURL())
	}

	return nil
}
