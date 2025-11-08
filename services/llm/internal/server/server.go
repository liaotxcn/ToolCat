package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"weave/services/llm/internal/chat"
	"weave/services/llm/internal/models"

	"golang.org/x/time/rate"
)

// 全局限流器-限制每秒最多10个请求
var limiter = rate.NewLimiter(rate.Every(time.Second), 10)

// 限流中间件
func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// XSS防护中间件
func xssMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置安全相关的HTTP头
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		next.ServeHTTP(w, r)
	})
}

// CORS跨域中间件
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 允许的域名列表
		allowedOrigins := map[string]bool{
			"http://localhost:8080": true,
			// 添加其他允许的域名t
		}

		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// 启动HTTP服务器
// pool: LLM连接池实例
func StartWebServer(pool *chat.LLMPool) {
	// 构建中间件链: XSS防护 -> CORS控制 -> 请求限流
	handlerChain := xssMiddleware(
		corsMiddleware(
			rateLimitMiddleware(
				http.DefaultServeMux,
			),
		),
	)

	// 对API路由应用限流中间件
	// 统一使用respondJSON函数
	http.Handle("/api/chat", rateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		chat, err := chat.NewChat(pool)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer chat.Close()

		prompt := chat.BuildPrompt(req.Message)
		response, err := chat.GetLLM().Call(context.Background(), prompt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"response": response,
		})
	})))

	// 获取对话历史记录的HTTP端点
	http.Handle("/api/history", rateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 获取LLM实例
		llm, err := pool.Get()
		if err != nil {
			http.Error(w, "Failed to get LLM instance", http.StatusInternalServerError)
			return
		}
		defer pool.Put(llm)

		// 创建chat服务实例
		repo := MemoryRepository()
		chatService := chat.NewChatService(repo, llm)

		// 获取历史记录
		histories, err := chatService.GetHistory(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 返回成功响应
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(APIResponse{
			Success: true,
			Data:    histories,
		})
	})))

	// 清空对话历史记录的HTTP端点
	http.Handle("/api/clear-history", rateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 获取LLM实例
		llm, err := pool.Get()
		if err != nil {
			http.Error(w, "Failed to get LLM instance", http.StatusInternalServerError)
			return
		}
		defer pool.Put(llm)

		// 创建chat服务实例
		repo := MemoryRepository()
		chatService := chat.NewChatService(repo, llm)

		// 清空历史记录
		if err := chatService.ClearHistory(context.Background()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 返回成功响应
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(APIResponse{
			Success: true,
			Data:    nil,
			Error:   "",
		})
	})))

	// 初始化内存历史记录
	memoryHistory = make([]*models.ChatHistory, 0)

	fmt.Println("Web server started at http://localhost:8080")
	// 启动服务器监听8080端口
	http.ListenAndServe(":8080", handlerChain)
}

// 标准API响应结构
type APIResponse struct {
	Success bool        `json:"success"`         // 请求是否成功
	Data    interface{} `json:"data"`            // 响应数据
	Error   string      `json:"error,omitempty"` // 错误信息(可选)
}

// 内存存储库实现
var (
	memoryHistory []*models.ChatHistory
	memoryMutex   sync.RWMutex
)

// MemoryRepository 内存存储库实现
// 实现ChatRepository接口
func MemoryRepository() chat.ChatRepository {
	return &memoryRepo{}
}

type memoryRepo struct{}

func (r *memoryRepo) SaveHistory(history *models.ChatHistory) error {
	memoryMutex.Lock()
	defer memoryMutex.Unlock()
	memoryHistory = append(memoryHistory, history)
	return nil
}

func (r *memoryRepo) GetHistories() ([]*models.ChatHistory, error) {
	memoryMutex.RLock()
	defer memoryMutex.RUnlock()
	// 返回历史记录的副本
	histories := make([]*models.ChatHistory, len(memoryHistory))
	copy(histories, memoryHistory)
	return histories, nil
}

func (r *memoryRepo) ClearHistories() error {
	memoryMutex.Lock()
	defer memoryMutex.Unlock()
	memoryHistory = make([]*models.ChatHistory, 0)
	return nil
}
