package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"toolcat/services/llm/internal/chat"

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
			"htp://localhost:8080": true,
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

	// 静态文件配置
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

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
