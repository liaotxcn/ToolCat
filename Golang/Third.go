package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// ==================== 高级并发安全数据结构 ====================

// ConcurrentMapWithShard 分片并发安全Map，减少锁竞争
type ConcurrentMapWithShard struct {
	shards     []*concurrentMapShard
	shardCount int
}

type concurrentMapShard struct {
	sync.RWMutex
	items map[string]interface{}
}

func NewConcurrentMapWithShard(shardCount int) *ConcurrentMapWithShard {
	shards := make([]*concurrentMapShard, shardCount)
	for i := 0; i < shardCount; i++ {
		shards[i] = &concurrentMapShard{
			items: make(map[string]interface{}),
		}
	}
	return &ConcurrentMapWithShard{
		shards:     shards,
		shardCount: shardCount,
	}
}

func (m *ConcurrentMapWithShard) getShard(key string) *concurrentMapShard {
	// 使用FNV哈希算法分配shard
	hash := fnv32(key)
	return m.shards[hash%uint32(m.shardCount)]
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func (m *ConcurrentMapWithShard) Set(key string, value interface{}) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

func (m *ConcurrentMapWithShard) Get(key string) (interface{}, bool) {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()
	val, ok := shard.items[key]
	return val, ok
}

// ConcurrentRingBuffer 并发安全环形缓冲区
type ConcurrentRingBuffer struct {
	buffer   []interface{}
	size     int
	head     int
	tail     int
	count    int
	mu       sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
}

func NewConcurrentRingBuffer(size int) *ConcurrentRingBuffer {
	rb := &ConcurrentRingBuffer{
		buffer: make([]interface{}, size),
		size:   size,
	}
	rb.notFull = sync.NewCond(&rb.mu)
	rb.notEmpty = sync.NewCond(&rb.mu)
	return rb
}

func (rb *ConcurrentRingBuffer) Put(item interface{}) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	for rb.count == rb.size {
		rb.notFull.Wait()
	}

	rb.buffer[rb.tail] = item
	rb.tail = (rb.tail + 1) % rb.size
	rb.count++
	rb.notEmpty.Signal()
}

func (rb *ConcurrentRingBuffer) Get() interface{} {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	for rb.count == 0 {
		rb.notEmpty.Wait()
	}

	item := rb.buffer[rb.head]
	rb.head = (rb.head + 1) % rb.size
	rb.count--
	rb.notFull.Signal()
	return item
}

// ==================== 高级Goroutine模式 ====================

// DynamicWorkerPool 动态调整大小的worker pool
type DynamicWorkerPool struct {
	tasks       chan func()
	workerCount int32
	minWorkers  int32
	maxWorkers  int32
	idleTimeout time.Duration
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

func NewDynamicWorkerPool(min, max int, idleTimeout time.Duration) *DynamicWorkerPool {
	pool := &DynamicWorkerPool{
		tasks:       make(chan func(), max*2),
		minWorkers:  int32(min),
		maxWorkers:  int32(max),
		idleTimeout: idleTimeout,
		stopChan:    make(chan struct{}),
	}

	// 启动最小数量的worker
	for i := 0; i < min; i++ {
		pool.addWorker()
	}

	// 启动监控goroutine动态调整worker数量
	go pool.monitor()

	return pool
}

func (p *DynamicWorkerPool) addWorker() {
	atomic.AddInt32(&p.workerCount, 1)
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()
		defer atomic.AddInt32(&p.workerCount, -1)

		idleTimer := time.NewTimer(p.idleTimeout)
		defer idleTimer.Stop()

		for {
			select {
			case task, ok := <-p.tasks:
				if !ok {
					return
				}
				task()
				// 重置空闲计时器
				if !idleTimer.Stop() {
					<-idleTimer.C
				}
				idleTimer.Reset(p.idleTimeout)
			case <-idleTimer.C:
				// 空闲超时，检查是否可以退出
				if atomic.LoadInt32(&p.workerCount) > p.minWorkers {
					return
				}
				idleTimer.Reset(p.idleTimeout)
			case <-p.stopChan:
				return
			}
		}
	}()
}

func (p *DynamicWorkerPool) monitor() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentWorkers := atomic.LoadInt32(&p.workerCount)
			queueLen := len(p.tasks)

			// 如果任务积压且worker不足，增加worker
			if queueLen > 0 && currentWorkers < p.maxWorkers {
				toAdd := int32(queueLen/2) + 1
				if toAdd > p.maxWorkers-currentWorkers {
					toAdd = p.maxWorkers - currentWorkers
				}
				for i := int32(0); i < toAdd; i++ {
					p.addWorker()
				}
			}
		case <-p.stopChan:
			return
		}
	}
}

func (p *DynamicWorkerPool) Submit(task func()) {
	p.tasks <- task
}

func (p *DynamicWorkerPool) Stop() {
	close(p.stopChan)
	close(p.tasks)
	p.wg.Wait()
}

// ==================== Context高级用法 ====================

func ContextAdvancedUsage() {
	// 创建带有超时和取消的context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 创建带有值的context
	ctx = context.WithValue(ctx, "requestID", "12345")

	// 启动多个goroutine处理任务
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		processTask1(ctx)
	}()

	go func() {
		defer wg.Done()
		processTask2(ctx)
	}()

	wg.Wait()
}

func processTask1(ctx context.Context) {
	// 检查context是否已取消
	select {
	case <-ctx.Done():
		fmt.Println("Task1 canceled:", ctx.Err())
		return
	default:
	}

	// 获取context中的值
	requestID := ctx.Value("requestID").(string)
	fmt.Println("Task1 processing with requestID:", requestID)

	// 模拟耗时操作
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("Task1 completed")
	case <-ctx.Done():
		fmt.Println("Task1 canceled during processing:", ctx.Err())
	}
}

func processTask2(ctx context.Context) {
	// 使用context的超时控制
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("Task2 completed (this shouldn't happen)")
	case <-ctx.Done():
		fmt.Println("Task2 canceled:", ctx.Err())
	}
}

// ==================== 高级Channel模式 ====================

// Multiplexing 多路复用模式
func ChannelMultiplexing() {
	ch1 := make(chan int)
	ch2 := make(chan string)

	go func() {
		for i := 0; i < 5; i++ {
			ch1 <- i
			time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
		}
		close(ch1)
	}()

	go func() {
		for _, s := range []string{"A", "B", "C", "D", "E"} {
			ch2 <- s
			time.Sleep(time.Duration(rand.Intn(400)) * time.Millisecond)
		}
		close(ch2)
	}()

	for {
		select {
		case num, ok := <-ch1:
			if !ok {
				ch1 = nil // 设置为nil后不会再被select选中
				continue
			}
			fmt.Println("Received number:", num)
		case str, ok := <-ch2:
			if !ok {
				ch2 = nil
				continue
			}
			fmt.Println("Received string:", str)
		}

		// 所有channel都关闭后退出
		if ch1 == nil && ch2 == nil {
			break
		}
	}
}

// RateLimiter 基于令牌桶的速率限制器
type RateLimiter struct {
	tokens      int32
	maxTokens   int32
	refillRate  time.Duration
	stopChan    chan struct{}
	mu          sync.Mutex
	waitingReqs []chan struct{}
}

func NewRateLimiter(maxTokens int, refillRate time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:     int32(maxTokens),
		maxTokens:  int32(maxTokens),
		refillRate: refillRate,
		stopChan:   make(chan struct{}),
	}

	go rl.refillTokens()
	return rl
}

func (rl *RateLimiter) refillTokens() {
	ticker := time.NewTicker(rl.refillRate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			if rl.tokens < rl.maxTokens {
				rl.tokens++
				// 如果有等待的请求，唤醒一个
				if len(rl.waitingReqs) > 0 {
					ch := rl.waitingReqs[0]
					rl.waitingReqs = rl.waitingReqs[1:]
					close(ch) // 通知可以继续执行
				}
			}
			rl.mu.Unlock()
		case <-rl.stopChan:
			return
		}
	}
}

func (rl *RateLimiter) Wait() {
	rl.mu.Lock()
	if rl.tokens > 0 {
		rl.tokens--
		rl.mu.Unlock()
		return
	}

	// 没有可用令牌，加入等待队列
	ch := make(chan struct{})
	rl.waitingReqs = append(rl.waitingReqs, ch)
	rl.mu.Unlock()

	// 等待通知
	<-ch
}

func (rl *RateLimiter) Stop() {
	close(rl.stopChan)
}

// ==================== 并发实用工具 ====================

// ErrorGroupWithContext 增强的ErrorGroup
func ErrorGroupWithContext() {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("Task1 completed")
			return nil
		case <-ctx.Done():
			fmt.Println("Task1 canceled")
			return ctx.Err()
		}
	})

	g.Go(func() error {
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("Task2 completed")
			return nil
		case <-ctx.Done():
			fmt.Println("Task2 canceled")
			return ctx.Err()
		}
	})

	g.Go(func() error {
		time.Sleep(500 * time.Millisecond)
		return fmt.Errorf("simulated error")
	})

	if err := g.Wait(); err != nil {
		fmt.Println("Error occurred:", err)
	}
}

// Semaphore 基于channel的信号量实现
type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(max int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, max),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}

func (s *Semaphore) TryAcquire(timeout time.Duration) bool {
	select {
	case s.ch <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

// ==================== 主函数 ====================

func main() {
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU()) // 使用所有CPU核心

	fmt.Println("=== 分片并发Map示例 ===")
	shardedMap := NewConcurrentMapWithShard(16)
	shardedMap.Set("key1", "value1")
	if val, ok := shardedMap.Get("key1"); ok {
		fmt.Println("获取到的值:", val)
	}

	fmt.Println("\n=== 动态Worker Pool示例 ===")
	pool := NewDynamicWorkerPool(2, 10, 10*time.Second)
	for i := 0; i < 20; i++ {
		taskID := i
		pool.Submit(func() {
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			fmt.Printf("任务 %d 由worker完成\n", taskID)
		})
	}
	time.Sleep(2 * time.Second) // 等待任务执行
	pool.Stop()

	fmt.Println("\n=== Context高级用法 ===")
	ContextAdvancedUsage()

	fmt.Println("\n=== Channel多路复用 ===")
	ChannelMultiplexing()

	fmt.Println("\n=== 速率限制器 ===")
	limiter := NewRateLimiter(5, time.Second)
	for i := 0; i < 10; i++ {
		go func(id int) {
			limiter.Wait()
			fmt.Printf("请求 %d 通过\n", id)
		}(i)
	}
	time.Sleep(3 * time.Second)
	limiter.Stop()

	fmt.Println("\n=== 增强的ErrorGroup ===")
	ErrorGroupWithContext()

	fmt.Println("\n=== 并发环形缓冲区 ===")
	rb := NewConcurrentRingBuffer(3)
	go func() {
		for i := 0; i < 5; i++ {
			rb.Put(i)
			fmt.Println("生产:", i)
		}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			item := rb.Get()
			fmt.Println("消费:", item)
		}
	}()

	time.Sleep(2 * time.Second)
}
