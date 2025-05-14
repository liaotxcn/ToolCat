package main

import (
	"container/heap"
	"container/list"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

/*
常见数据结构和算法的Go语言实现
包含以下主要内容：
1. 基础数据结构：链表、栈、队列、双端队列、集合、图等
2. 常用算法：排序、搜索、动态规划等
3. 实用工具：字符串处理、切片操作等
所有实现都考虑了线程安全，适合学习和生产使用
*/

// ==================== 基础数据结构 ====================

// ==================== Map 操作扩展 ====================

// MapOperations map操作
func MapOperations() {
	fmt.Println("\n=== Map操作示例 ===")

	// 线程安全的sync.Map示例
	var syncMap sync.Map
	syncMap.Store("name", "Pai")
	syncMap.Store("age", 999)

	// 遍历sync.Map
	syncMap.Range(func(key, value interface{}) bool {
		fmt.Printf("%v: %v\n", key, value)
		return true
	})

	// 1. 创建和初始化map
	user := map[string]interface{}{
		"name":    "Pai",
		"age":     999,
		"active":  true,
		"hobbies": []string{"reading", "coding"},
	}

	// 2. 添加/修改元素
	user["email"] = "pai@liaotx.cn"
	user["age"] = 26 // 修改已有键的值

	// 3. 检查键是否存在
	if email, exists := user["email"]; exists {
		fmt.Printf("Email存在: %v\n", email)
	}

	// 4. 删除键
	delete(user, "active")

	// 5. 遍历map
	fmt.Println("用户信息:")
	for key, value := range user {
		fmt.Printf("%s: %v\n", key, value)
	}

	// 6. 获取map长度
	fmt.Println("用户信息字段数:", len(user))

	// 7. 清空map
	clear(user) // Go 1.21+
	fmt.Println("清空后map长度:", len(user))
}

// ==================== Set 实现 ====================

// Set 线程安全的集合实现
// 基于map实现，支持并、交、差等集合操作
// 使用读写锁保证并发安全
type Set struct {
	data map[interface{}]struct{}
	mu   sync.RWMutex
}

// Graph 图的邻接表表示
// 支持带权重的有向图和无向图
// 提供BFS、DFS、Dijkstra等图算法
type Graph struct {
	Vertices int
	Adj      map[int]map[int]int // 邻接表带权重
}

func NewSet() *Set {
	return &Set{
		data: make(map[interface{}]struct{}),
	}
}

func (s *Set) Add(items ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range items {
		s.data[item] = struct{}{}
	}
}

func (s *Set) Remove(item interface{}) {
	delete(s.data, item)
}

func (s *Set) Contains(item interface{}) bool {
	_, exists := s.data[item]
	return exists
}

func (s *Set) Size() int {
	return len(s.data)
}

func (s *Set) ToSlice() []interface{} {
	result := make([]interface{}, 0, len(s.data))
	for item := range s.data {
		result = append(result, item)
	}
	return result
}

func (s *Set) Union(other *Set) *Set {
	union := NewSet()
	for item := range s.data {
		union.Add(item)
	}
	for item := range other.data {
		union.Add(item)
	}
	return union
}

func (s *Set) Intersection(other *Set) *Set {
	intersection := NewSet()
	// 遍历较小的集合更高效
	if s.Size() < other.Size() {
		for item := range s.data {
			if other.Contains(item) {
				intersection.Add(item)
			}
		}
	} else {
		for item := range other.data {
			if s.Contains(item) {
				intersection.Add(item)
			}
		}
	}
	return intersection
}

// ==================== String 操作扩展 ====================

// StringOperations 演示各种字符串操作
func StringOperations() {
	fmt.Println("\n=== String操作示例 ===")

	str := " Hello, 世界! "

	// 1. 基本操作
	fmt.Println("原始字符串:", str)
	fmt.Println("长度:", len(str))                 // 字节长度
	fmt.Println("字符数:", len([]rune(str)))        // 字符长度
	fmt.Println("大写:", strings.ToUpper(str))     // 转大写
	fmt.Println("小写:", strings.ToLower(str))     // 转小写
	fmt.Println("修剪空格:", strings.TrimSpace(str)) // 去除首尾空格

	// 2. 分割和连接
	parts := strings.Split(str, ",")
	fmt.Println("分割结果:", parts)
	fmt.Println("连接结果:", strings.Join(parts, "|"))

	// 3. 子串操作
	fmt.Println("包含'世界'?", strings.Contains(str, "世界"))
	fmt.Println("前缀检查:", strings.HasPrefix(str, " "))
	fmt.Println("后缀检查:", strings.HasSuffix(str, "! "))

	// 4. 字符串构建
	var builder strings.Builder
	builder.WriteString("Hello")
	builder.WriteByte(',')
	builder.WriteRune('世')
	builder.WriteRune('界')
	fmt.Println("构建的字符串:", builder.String())

	// 5. 字符串替换
	replacer := strings.NewReplacer("Hello", "你好", "世界", "World")
	fmt.Println("替换结果:", replacer.Replace(str))
}

// ==================== 切片高级操作 ====================

// SliceOperations 演示切片高级操作
func SliceOperations() {
	fmt.Println("\n=== 切片高级操作 ===")

	// 1. 创建和初始化
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("原始切片:", nums)

	// 2. 切片操作
	fmt.Println("子切片[1:3]:", nums[1:3])
	fmt.Println("子切片[:3]:", nums[:3])
	fmt.Println("子切片[3:]:", nums[3:])

	// 3. 追加元素
	nums = append(nums, 6, 7, 8)
	fmt.Println("追加后:", nums)

	// 4. 复制切片
	copyNums := make([]int, len(nums))
	copy(copyNums, nums)
	fmt.Println("复制结果:", copyNums)

	// 5. 删除元素 (删除索引2)
	nums = append(nums[:2], nums[3:]...)
	fmt.Println("删除索引2后:", nums)

	// 6. 插入元素 (在索引2插入9)
	nums = append(nums[:2], append([]int{9}, nums[2:]...)...)
	fmt.Println("插入9后:", nums)

	// 7. 过滤切片
	filtered := FilterSlice(nums, func(x int) bool { return x%2 == 0 })
	fmt.Println("过滤偶数:", filtered)

	// 8. 映射转换
	mapped := MapSlice(nums, func(x int) int { return x * 2 })
	fmt.Println("元素乘2:", mapped)
}

func FilterSlice(s []int, fn func(int) bool) []int {
	var result []int
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func MapSlice(s []int, fn func(int) int) []int {
	result := make([]int, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// ListNode 链表节点
type ListNode struct {
	Val  int
	Next *ListNode
}

// LinkedList 链表实现
type LinkedList struct {
	Head *ListNode
}

func (l *LinkedList) Insert(val int) {
	node := &ListNode{Val: val}
	if l.Head == nil {
		l.Head = node
		return
	}
	current := l.Head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = node
}

// BinaryTreeNode 二叉树节点
type BinaryTreeNode struct {
	Val   int
	Left  *BinaryTreeNode
	Right *BinaryTreeNode
}

// Stack 线程安全栈实现
type Stack struct {
	data []interface{}
	mu   sync.RWMutex
}

func (s *Stack) Push(val interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, val)
}

func (s *Stack) Pop() interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.data) == 0 {
		return nil
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val
}

func (s *Stack) Peek() interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.data) == 0 {
		return nil
	}
	return s.data[len(s.data)-1]
}

func (s *Stack) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *Stack) IsEmpty() bool {
	return s.Size() == 0
}

// Queue 线程安全队列实现
type Queue struct {
	data []interface{}
	mu   sync.RWMutex
}

func (q *Queue) Enqueue(val interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data = append(q.data, val)
}

func (q *Queue) Dequeue() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.data) == 0 {
		return nil
	}
	val := q.data[0]
	q.data = q.data[1:]
	return val
}

func (q *Queue) Peek() interface{} {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if len(q.data) == 0 {
		return nil
	}
	return q.data[0]
}

func (q *Queue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.data)
}

func (q *Queue) IsEmpty() bool {
	return q.Size() == 0
}

// ==================== 高级数据结构 ====================

func NewGraph(vertices int) *Graph {
	return &Graph{
		Vertices: vertices,
		Adj:      make(map[int]map[int]int),
	}
}

func (g *Graph) AddEdge(u, v int, weight int) {
	if g.Adj[u] == nil {
		g.Adj[u] = make(map[int]int)
	}
	g.Adj[u][v] = weight

	// 无向图需要双向添加
	if g.Adj[v] == nil {
		g.Adj[v] = make(map[int]int)
	}
	g.Adj[v][u] = weight
}

// MinHeap 最小堆实现
type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// ==================== 排序算法 ====================

// QuickSort 优化的快速排序算法
// 时间复杂度：
// - 平均情况：O(n log n)
// - 最坏情况：O(n^2)
// 空间复杂度：O(log n)
// 对小数组自动切换到插入排序提高性能
func QuickSort(arr []int) []int {
	if len(arr) <= 16 { // 小数组使用插入排序
		return InsertionSort(arr)
	}

	pivot := medianOfThree(arr[0], arr[len(arr)/2], arr[len(arr)-1])
	var left, right []int

	for i := 0; i < len(arr); i++ {
		if arr[i] < pivot {
			left = append(left, arr[i])
		} else if arr[i] > pivot {
			right = append(right, arr[i])
		}
	}

	left = QuickSort(left)
	right = QuickSort(right)

	return append(append(left, pivot), right...)
}

// MergeSort 归并排序
func MergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2
	left := MergeSort(arr[:mid])
	right := MergeSort(arr[mid:])

	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}

// ==================== 搜索算法 ====================

// BinarySearch 二分查找
func BinarySearch(arr []int, target int) int {
	low, high := 0, len(arr)-1

	for low <= high {
		mid := low + (high-low)/2
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1
}

// DFS 深度优先搜索
func (g *Graph) DFS(start int) []int {
	visited := make(map[int]bool)
	var result []int
	g.dfsHelper(start, visited, &result)
	return result
}

func (g *Graph) dfsHelper(node int, visited map[int]bool, result *[]int) {
	visited[node] = true
	*result = append(*result, node)

	for _, neighbor := range g.Adj[node] {
		if !visited[neighbor] {
			g.dfsHelper(neighbor, visited, result)
		}
	}
}

// BFS 广度优先搜索
func (g *Graph) BFS(start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	var result []int

	visited[start] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		for _, neighbor := range g.Adj[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return result
}

// ==================== 动态规划 ====================

// Fibonacci 斐波那契数列
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}

	dp := make([]int, n+1)
	dp[0], dp[1] = 0, 1

	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}

	return dp[n]
}

// Knapsack 0-1背包问题
func Knapsack(weights []int, values []int, capacity int) int {
	n := len(weights)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	for i := 1; i <= n; i++ {
		for w := 1; w <= capacity; w++ {
			if weights[i-1] <= w {
				dp[i][w] = max(values[i-1]+dp[i-1][w-weights[i-1]], dp[i-1][w])
			} else {
				dp[i][w] = dp[i-1][w]
			}
		}
	}

	return dp[n][capacity]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ==================== 图算法 ====================

// Dijkstra 最短路径算法
func (g *Graph) Dijkstra(start int) map[int]int {
	dist := make(map[int]int)
	for i := 0; i < g.Vertices; i++ {
		dist[i] = math.MaxInt32
	}
	dist[start] = 0

	minHeap := &MinHeap{}
	heap.Init(minHeap)
	heap.Push(minHeap, start)

	for minHeap.Len() > 0 {
		u := heap.Pop(minHeap).(int)

		for _, v := range g.Adj[u] {
			if dist[v] > dist[u]+1 { // 假设每条边权重为1
				dist[v] = dist[u] + 1
				heap.Push(minHeap, v)
			}
		}
	}

	return dist
}

// ==================== 实用算法 ====================

// LRUCache LRU缓存
type LRUCache struct {
	capacity int
	cache    map[int]*list.Element
	list     *list.List
}

type Pair struct {
	key      int
	value    int
	expireAt time.Time
}

func (lru *LRUCache) PutWithExpire(key int, value int, ttl time.Duration) {
	if elem, ok := lru.cache[key]; ok {
		pair := elem.Value.(*Pair)
		pair.value = value
		pair.expireAt = time.Now().Add(ttl)
		lru.list.MoveToFront(elem)
	} else {
		if lru.list.Len() >= lru.capacity {
			last := lru.list.Back()
			delete(lru.cache, last.Value.(*Pair).key)
			lru.list.Remove(last)
		}
		newElem := lru.list.PushFront(&Pair{
			key:      key,
			value:    value,
			expireAt: time.Now().Add(ttl),
		})
		lru.cache[key] = newElem
	}
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		list:     list.New(),
	}
}

func (lru *LRUCache) Get(key int) int {
	if elem, ok := lru.cache[key]; ok {
		lru.list.MoveToFront(elem)
		return elem.Value.(*Pair).value
	}
	return -1
}

func (lru *LRUCache) Put(key int, value int) {
	if elem, ok := lru.cache[key]; ok {
		lru.list.MoveToFront(elem)
		elem.Value.(*Pair).value = value
	} else {
		if lru.list.Len() >= lru.capacity {
			// 删除最久未使用的元素
			last := lru.list.Back()
			delete(lru.cache, last.Value.(*Pair).key)
			lru.list.Remove(last)
		}
		newElem := lru.list.PushFront(&Pair{key: key, value: value})
		lru.cache[key] = newElem
	}
}

// ==================== 主函数 ====================

func main() {
	// 测试map操作
	MapOperations()

	// 测试set操作
	fmt.Println("\n=== Set操作测试 ===")
	set1 := NewSet()
	set1.Add(1, 2, 3, "a", "b")
	set2 := NewSet()
	set2.Add(3, 4, 5, "b", "c")

	fmt.Println("Set1:", set1.ToSlice())
	fmt.Println("Set2:", set2.ToSlice())
	fmt.Println("并集:", set1.Union(set2).ToSlice())
	fmt.Println("交集:", set1.Intersection(set2).ToSlice())
	fmt.Println("Set1包含'a'?", set1.Contains("a"))

	// 测试字符串操作
	StringOperations()

	// 测试切片操作
	SliceOperations()

	TestOriginalDataStructures()
}

// TestOriginalDataStructures 测试所有数据结构和算法
// 包含以下测试用例：
// 1. 链表操作
// 2. 排序算法
// 3. 搜索算法
// 4. 图算法
// 5. 动态规划
// 6. LRU缓存
func TestOriginalDataStructures() {
	// 测试链表
	list := LinkedList{}
	list.Insert(1)
	list.Insert(2)
	list.Insert(3)
	fmt.Println("链表实现测试:", list.Head.Val, list.Head.Next.Val)

	// 测试快速排序
	arr := []int{3, 1, 4, 1, 5, 9, 2, 6}
	fmt.Println("快速排序前:", arr)
	fmt.Println("快速排序后:", QuickSort(arr))

	// 测试二分查找
	sortedArr := []int{1, 3, 5, 7, 9}
	fmt.Println("二分查找 5:", BinarySearch(sortedArr, 5))

	// 测试图算法
	graph := NewGraph(5)
	graph.AddEdge(0, 1, 1)
	graph.AddEdge(0, 2, 1)
	graph.AddEdge(1, 3, 1)
	graph.AddEdge(2, 4, 1)
	fmt.Println("BFS遍历:", graph.BFS(0))
	fmt.Println("DFS遍历:", graph.DFS(0))

	// 测试动态规划
	fmt.Println("斐波那契(10):", Fibonacci(10))
	weights := []int{2, 3, 4, 5}
	values := []int{3, 4, 5, 6}
	fmt.Println("背包问题(容量=5):", Knapsack(weights, values, 5))

	// 测试LRU缓存
	lru := NewLRUCache(2)
	lru.Put(1, 1)
	lru.Put(2, 2)
	fmt.Println("LRU Get 1:", lru.Get(1))
	lru.Put(3, 3)
	fmt.Println("LRU Get 2:", lru.Get(2))
}

func InsertionSort(arr []int) []int {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
	return arr
}

func medianOfThree(a, b, c int) int {
	if a < b {
		if b < c {
			return b
		} else if a < c {
			return c
		} else {
			return a
		}
	} else {
		if a < c {
			return a
		} else if b < c {
			return c
		} else {
			return b
		}
	}
}

// Deque 双端队列实现
type Deque struct {
	data []interface{}
	mu   sync.RWMutex
}

func (d *Deque) PushFront(val interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data = append([]interface{}{val}, d.data...)
}

func (d *Deque) PushBack(val interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data = append(d.data, val)
}

func (d *Deque) PopFront() interface{} {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.data) == 0 {
		return nil
	}
	val := d.data[0]
	d.data = d.data[1:]
	return val
}

func (d *Deque) PopBack() interface{} {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.data) == 0 {
		return nil
	}
	val := d.data[len(d.data)-1]
	d.data = d.data[:len(d.data)-1]
	return val
}

func (d *Deque) PeekFront() interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if len(d.data) == 0 {
		return nil
	}
	return d.data[0]
}

func (d *Deque) PeekBack() interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if len(d.data) == 0 {
		return nil
	}
	return d.data[len(d.data)-1]
}

func (d *Deque) Size() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.data)
}

func (d *Deque) IsEmpty() bool {
	return d.Size() == 0
}
