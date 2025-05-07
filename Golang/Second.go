package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// ==================== 并发安全数据结构 ====================

// ConcurrentMap 并发安全的Map
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

// ConcurrentQueue 并发安全的队列
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

// WorkerPool 工作池模式
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

// Fan-out 模式：一个channel分发给多个worker
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
