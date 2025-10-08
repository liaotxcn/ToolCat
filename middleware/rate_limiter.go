package middleware

import (
	"net/http"
	"sync"
	"time"

	"toolcat/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RateLimiter 限流中间件
// rate: 每秒生成的令牌数
// burst: 最大令牌桶容量
func RateLimiter(rate float64, burst int) gin.HandlerFunc {
	// 创建一个令牌桶管理器，按IP地址区分不同客户端
	bucketManager := NewTokenBucketManager(rate, burst)

	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()

		// 尝试从令牌桶中获取令牌
		if !bucketManager.Allow(clientIP) {
			// 记录限流日志
			pkg.With(
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("client_ip", clientIP),
			).Info("Rate limit exceeded")

			// 返回429 Too Many Requests状态码
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}

// TokenBucket 实现令牌桶算法
type TokenBucket struct {
	rate       float64    // 每秒生成的令牌数
	capacity   int        // 令牌桶容量
	tokens     float64    // 当前令牌数量
	lastRefill time.Time  // 上次填充令牌的时间
	mtx        sync.Mutex // 互斥锁，保证线程安全
}

// NewTokenBucket 创建一个新的令牌桶
func NewTokenBucket(rate float64, capacity int) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		capacity:   capacity,
		tokens:     float64(capacity),
		lastRefill: time.Now(),
	}
}

// Allow 尝试从令牌桶中获取一个令牌
// 如果获取成功返回true，否则返回false
func (tb *TokenBucket) Allow() bool {
	return tb.Take(1)
}

// Take 尝试从令牌桶中获取指定数量的令牌
func (tb *TokenBucket) Take(count int) bool {
	tb.mtx.Lock()
	defer tb.mtx.Unlock()

	// 计算自上次填充以来应该生成的令牌数
	now := time.Now()
	duration := now.Sub(tb.lastRefill).Seconds()
	newTokens := duration * tb.rate

	// 更新令牌数量和上次填充时间
	if newTokens > 0 {
		tb.tokens = min(float64(tb.capacity), tb.tokens+newTokens)
		tb.lastRefill = now
	}

	// 检查是否有足够的令牌
	if tb.tokens >= float64(count) {
		tb.tokens -= float64(count)
		return true
	}

	return false
}

// TokenBucketManager 管理多个客户端的令牌桶
type TokenBucketManager struct {
	rate     float64
	capacity int
	buckets  map[string]*TokenBucket
	mtx      sync.RWMutex
}

// NewTokenBucketManager 创建一个新的令牌桶管理器
func NewTokenBucketManager(rate float64, capacity int) *TokenBucketManager {
	return &TokenBucketManager{
		rate:     rate,
		capacity: capacity,
		buckets:  make(map[string]*TokenBucket),
	}
}

// Allow 检查指定客户端是否可以继续请求
func (tbm *TokenBucketManager) Allow(clientID string) bool {
	// 先尝试读取锁获取令牌桶
	tbm.mtx.RLock()
	bucket, exists := tbm.buckets[clientID]
	tbm.mtx.RUnlock()

	// 如果令牌桶不存在，创建一个新的
	if !exists {
		tbm.mtx.Lock()
		// 双重检查，防止并发创建
		bucket, exists = tbm.buckets[clientID]
		if !exists {
			bucket = NewTokenBucket(tbm.rate, tbm.capacity)
			tbm.buckets[clientID] = bucket
		}
		tbm.mtx.Unlock()
	}

	// 尝试从令牌桶中获取令牌
	return bucket.Allow()
}

// 辅助函数：返回两个数中的较小值
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
