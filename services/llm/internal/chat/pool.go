package chat

import (
	"sync"
	"time"

	"github.com/tmc/langchaingo/llms"
)

// 管理LLM连接资源池
type LLMPool struct {
	pool     chan llms.LLM            // 存放LLM连接缓冲通道
	creator  func() (llms.LLM, error) // 创建新连接回调函数
	mu       sync.Mutex               // 保护并发访问互斥锁
	timeout  time.Duration            // 获取连接超时时间
	maxRetry int                      // 创建连接最大重试次数
}

// 创建新LLM连接池
// size: 连接池容量
// creator: 创建新连接的函数
func NewLLMPool(size int, creator func() (llms.LLM, error)) *LLMPool {
	return &LLMPool{
		pool:     make(chan llms.LLM, size),
		creator:  creator,
		timeout:  5 * time.Second, // 默认5秒超时
		maxRetry: 3,               // 默认重试3次
	}
}

// 从池中获取一个LLM连接
// 如果池中没有可用连接，会等待直到超时或获取成功
// 超时后会尝试创建新连接
func (p *LLMPool) Get() (llms.LLM, error) {
	select {
	case llm := <-p.pool:
		return llm, nil
	case <-time.After(p.timeout):
		return p.createWithRetry() // 超时后重试创建
	}
}

// 带重试机制的连接创建方法
// 使用指数退避策略进行重试
func (p *LLMPool) createWithRetry() (llms.LLM, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var lastErr error
	for i := 0; i < p.maxRetry; i++ {
		llm, err := p.creator()
		if err == nil {
			return llm, nil
		}
		lastErr = err
		time.Sleep(time.Second * time.Duration(i+1)) // 指数退避
	}
	return nil, lastErr
}

// 将连接放回池中
// 如果池已满，会安全关闭连接
func (p *LLMPool) Put(llm llms.LLM) {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case p.pool <- llm: // 尝试放回池中
	default:
		// 池已满，安全关闭连接
		if closer, ok := llm.(interface{ Close() error }); ok {
			_ = closer.Close()
		}
	}
}

// 设置获取连接的超时时间
func (p *LLMPool) SetTimeout(d time.Duration) {
	p.timeout = d
}

// 设置创建连接的最大重试次数
func (p *LLMPool) SetMaxRetry(n int) {
	p.maxRetry = n
}

// 获取连接池的统计信息
func (p *LLMPool) Stats() map[string]interface{} {
	return map[string]interface{}{
		"capacity":  cap(p.pool), // 池总容量
		"available": len(p.pool), // 当前可用连接数
		"timeout":   p.timeout,   // 当前超时设置
	}
}
