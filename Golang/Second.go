/*
Go语言并发核心模式和数据结构
包含以下主要内容：
1. 并发安全数据结构：Map、Slice、Queue等
2. Goroutine管理：工作池、优雅退出等模式
3. Channel高级模式：Fan-in、Fan-out、超时控制等
4. 原子操作与并发原语：计数器、Once等
5. 经典并发模式：生产者消费者等
*/
package主干

import (
	"container/list"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// ==================== 并发安全数据结构 ====================

// ConcurrentMap 并发安全的Map实现
// 使用读写锁(RWMutex)实现:
// - 读操作(Get)可以并发执行
// - 写操作(Set/Delete)互斥执行
// 性能考虑:
// - 适合读多写少的场景
// - 大量写操作时性能会下降
type ConcurrentMap struct {
	sync.RWMutex
	items map[string]interface{}
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		items: make(map[string]interface{}),
	}
}

func (m *ConcurrentMap) Set(key string, value interface{}) {
	m.Lock()
	defer m.Unlock()
	m.items[key] = value
}

func (m *ConcurrentMap) Get(key string) (interface{}, bool) {
	m.RLock()
	defer m.RUnlock()
	val, ok := m.items[key]
	return val, ok
}

func (m *ConcurrentMap) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.items, key)
}

// ConcurrentSlice 并发安全的Slice
type ConcurrentSlice struct {
	sync.RWMutex
	items []interface{}
}

func (cs *ConcurrentSlice) Append(item interface{}) {
	cs.Lock()
	defer cs.Unlock()
	cs.items = append(cs.items, item)
}

func (cs *ConcurrentSlice) Len() int {
	cs.RLock()
	defer cs.RUnlock()
	return len(cs.items)
}

// ConcurrentQueue 并发安全的FIFO队列
// 实现特点:
// - 基于container/list链表实现
// - 无容量限制(注意内存使用)
// - 使用互斥锁保证线程安全
// 适用场景:
// - 任务队列
// - 事件处理
type ConcurrentQueue struct {
	sync.Mutex
	queue *list.List
}

func NewConcurrentQueue() *ConcurrentQueue {
	return &ConcurrentQueue{
		queue: list.New(),
	}
}

func (q *ConcurrentQueue) Enqueue(item interface{}) {
	q.Lock()
	defer q.Unlock()
	q.queue.PushBack(item)
}

func (q *ConcurrentQueue) Dequeue() (interface{}, bool) {
	q.Lock()
	defer q.Unlock()
	if q.queue.Len() == 0 {
		return nil, false
	}
	element := q.queue.Front()
	q.queue.Remove(element)
	return element.Value, true
}

// ==================== Goroutine 管理与模式 ====================

// WorkerPool 工作池模式实现
// 核心组件:
// - tasks chan: 带缓冲的任务通道
// - wg WaitGroup: 等待所有worker完成
// 工作流程:
// 1. 初始化时创建指定数量的worker goroutine
// 2. 每个worker从tasks通道获取并执行任务
// 3. 调用Wait()关闭通道并等待所有worker退出
// 注意事项:
// - 任务通道缓冲区大小建议为worker数量的2倍
// - 任务不应阻塞，否则会影响其他任务执行
type WorkerPool struct {
	tasks    chan func()
	wg       sync.WaitGroup
	poolSize int
}

func NewWorkerPool(poolSize int) *WorkerPool {
	wp := &WorkerPool{
		tasks:    make(chan func(), poolSize*2),
		poolSize: poolSize,
	}
	wp.wg.Add(poolSize)
	for i := 0; i < poolSize; i++ {
		go wp.worker()
	}
	return wp
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for task := range wp.tasks {
		task()
	}
}

func (wp *WorkerPool) Submit(task func()) {
	wp.tasks <- task
}

func (wp *WorkerPool) Wait() {
	close(wp.tasks)
	wp.wg.Wait()
}

// Goroutine 优雅退出模式
// gracefulShutdown goroutine优雅退出模式
// 参数：
// - stopChan: 接收停止信号的channel
// - doneChan: 通知清理完成的channel
// 特性：
// - 支持周期性任务和清理操作
// - 可响应外部停止信号
func gracefulShutdown(stopChan <-chan struct{}, doneChan chan<- struct{}) {
	defer close(doneChan)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			fmt.Println("接收到停止信号，开始清理...")
			// 执行清理逻辑
			time.Sleep(1 * time.Second) // 模拟清理操作
			fmt.Println("清理完成，退出goroutine")
			return
		case <-ticker.C:
			fmt.Println("执行周期性任务...")
			// 正常业务逻辑
		}
	}
}

// ==================== Channel 高级模式 ====================

// Fan-in 模式：多个输入channel合并为一个
// fanIn Channel扇入模式
// 功能：合并多个输入channel到一个输出channel
// 参数：
// - inputs: 多个输入channel
// 返回值：
// - 合并后的输出channel
// 注意：当所有输入channel关闭后，输出channel会自动关闭
func fanIn(inputs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, in := range inputs {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for n := range ch {
				out <- n
			}
		}(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// fanOut Channel扇出模式
// 功能：将一个输入channel分发给多个worker处理
// 参数：
// - input: 输入channel
// - workers: worker数量
// - process: 处理函数
// 注意：会等待所有worker完成任务
func fanOut(input <-chan int, workers int, process func(int)) {
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for n := range input {
				process(n)
			}
		}()
	}

	wg.Wait()
}

// 带超时的Channel操作
func channelWithTimeout() {
	ch := make(chan string, 1)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "操作结果"
	}()

	select {
	case res := <-ch:
		fmt.Println("收到结果:", res)
	case <-time.After(1 * time.Second):
		fmt.Println("操作超时")
	}
}

// ==================== 原子操作与并发原语 ====================

// AtomicCounter 原子计数器
type AtomicCounter struct {
	value int64
}

func (c *AtomicCounter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

// OnceDemo sync.Once示例
func OnceDemo() {
	var once sync.Once
	initFn := func() {
		fmt.Println("初始化操作，只执行一次")
	}

	for i := 0; i < 5; i++ {
		go func() {
			once.Do(initFn)
		}()
	}
	time.Sleep(500 * time.Millisecond)
}

// ==================== 并发模式示例 ====================

// ProducerConsumer 生产者消费者模式
func ProducerConsumer() {
	const (
		producers = 3
		consumers = 2
		items     = 10
	)

	ch := make(chan int, 5)
	var wg sync.WaitGroup

	// 生产者
	for i := 0; i < producers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < items; j++ {
				item := id*100 + j
				ch <- item
				fmt.Printf("生产者%d 生产: %d\n", id, item)
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			}
		}(i)
	}

	// 消费者
	for i := 0; i < consumers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < items*producers/consumers; j++ {
				item := <-ch
				fmt.Printf("消费者%d 消费: %d\n", id, item)
				time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	close(ch)
}

// ==================== 主函数 ====================

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== 并发安全数据结构示例 ===")
	concurrentMap := NewConcurrentMap()
	concurrentMap.Set("key1", "value1")
	if val, ok := concurrentMap.Get("key1"); ok {
		fmt.Println("获取到的值:", val)
	}

	fmt.Println("\n=== WorkerPool 示例 ===")
	wp := NewWorkerPool(3)
	for i := 0; i < 10; i++ {
		taskID := i
		wp.Submit(func() {
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("任务 %d 完成\n", taskID)
		})
	}
	wp.Wait()

	fmt.Println("\n=== 优雅退出示例 ===")
	stopChan := make(chan struct{})
	doneChan := make(chan struct{})
	go gracefulShutdown(stopChan, doneChan)
	time.Sleep(2 * time.Second)
	close(stopChan)
	<-doneChan

	fmt.Println("\n=== 原子计数器示例 ===")
	counter := &AtomicCounter{}
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				counter.Increment()
			}
		}()
	}
	time.Sleep(1 * time.Second)
	fmt.Println("最终计数值:", counter.Value())

	fmt.Println("\n=== 生产者消费者模式 ===")
	ProducerConsumer()

	fmt.Println("\n=== sync.Once 示例 ===")
	OnceDemo()

	fmt.Println("\n=== 带超时的Channel操作 ===")
	channelWithTimeout()
}

// ConcurrentRingBuffer 并发安全的环形缓冲区
// 特性：
// - 固定容量，循环使用空间
// - 支持多生产者多消费者场景
// - 提供阻塞和非阻塞操作
type ConcurrentRingBuffer struct {
	buf      []interface{}
	size     int
	capacity int
	head     int
	tail     int
	mu       sync.RWMutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
}

func NewConcurrentRingBuffer(capacity int) *ConcurrentRingBuffer {
	rb := &ConcurrentRingBuffer{
		buf:      make([]interface{}, capacity),
		capacity: capacity,
	}
	rb.notEmpty = sync.NewCond(&rb.mu)
	rb.notFull = sync.NewCond(&rb.mu)
	return rb
}

func (rb *ConcurrentRingBuffer) Put(item interface{}) bool {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.size == rb.capacity {
		return false
	}

	rb.buf[rb.tail] = item
	rb.tail = (rb.tail + 1) % rb.capacity
	rb.size++
	rb.notEmpty.Signal()
	return true
}

func (rb *ConcurrentRingBuffer) Take() (interface{}, bool) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.size == 0 {
		return nil, false
	}

	item := rb.buf[rb.head]
	rb.head = (rb.head + 1) % rb.capacity
	rb.size--
	rb.notFull.Signal()
	return item, true
}

// PubSub 发布订阅模式实现
// 核心机制:
// - subs map: 主题到订阅者channel的映射
// - RWMutex: 保护订阅关系的并发访问
// 消息流程:
// 1. 订阅者通过Subscribe()注册接收channel
// 2. 发布者通过Publish()向主题发布消息
// 3. 消息会广播给该主题所有订阅者
// 特性:
// - 支持多主题
// - 每个订阅者有独立带缓冲的channel(100)
// - 线程安全的订阅和发布操作
type PubSub struct {
	mu     sync.RWMutex
	subs   map[string][]chan interface{}
	closed bool
}

func NewPubSub() *PubSub {
	return &PubSub{
		subs: make(map[string][]chan interface{}),
	}
}

func (ps *PubSub) Subscribe(topic string) <-chan interface{} {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		return nil
	}

	ch := make(chan interface{}, 100)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch
}

func (ps *PubSub) Publish(topic string, msg interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	for _, ch := range ps.subs[topic] {
		ch <- msg
	}
}

func (ps *PubSub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, subs := range ps.subs {
			for _, ch := range subs {
				close(ch)
			}
		}
	}
}

// RateLimiter 令牌桶限流算法实现
// 算法原理:
// - 令牌以固定速率(r tokens/sec)添加到桶中
// - 桶容量最大为c个令牌
// - 请求消耗令牌，无令牌时拒绝
// 实现细节:
// - 使用mutex保护令牌计数
// - 每次Allow()计算时间差并补充令牌
// 参数建议:
// - rate: 根据系统负载设置
// - capacity: 应对突发流量的能力
type RateLimiter struct {
	capacity int64
	tokens   int64
	rate     float64 // tokens per second
	lastTime time.Time
	mu       sync.Mutex
}

func NewRateLimiter(rate float64, capacity int64) *RateLimiter {
	return &RateLimiter{
		capacity: capacity,
		tokens:   capacity,
		rate:     rate,
		lastTime: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastTime).Seconds()
	rl.lastTime = now

	// 添加新令牌
	rl.tokens = rl.tokens + int64(elapsed*rl.rate)
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}

	// 检查是否有足够令牌
	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

// Barrier 并发屏障
type Barrier struct {
	count    int
	n        int
	mu       sync.Mutex
	cond     *sync.Cond
	released bool
}

func NewBarrier(n int) *Barrier {
	b := &Barrier{n: n}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *Barrier) Wait() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.released {
		return
	}

	b.count++
	if b.count < b.n {
		for !b.released {
			b.cond.Wait()
		}
	} else {
		b.released = true
		b.cond.Broadcast()
	}
}

func (b *Barrier) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.count = 0
	b.released = false
}
