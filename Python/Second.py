"""
Python 数据结构与算法(2)
包含：基础数据结构 + 进阶算法实现 + 实战模板
"""


# ==================== 基础数据结构增强版 ====================

class AdvancedList:
    """列表进阶操作"""

    @staticmethod
    def list_comprehension():
        """列表生成式高级用法"""
        # 生成二维矩阵
        matrix = [[i * j for j in range(5)] for i in range(5)]
        print("5x5矩阵:", matrix)

        # 带条件的生成式
        even_squares = [x ** 2 for x in range(10) if x % 2 == 0]
        print("偶数的平方:", even_squares)

    @staticmethod
    def slicing_techniques():
        """切片高级技巧"""
        lst = list(range(10))
        print("原始列表:", lst)
        print("每隔2个取一个:", lst[::2])
        print("反转列表:", lst[::-1])
        print("部分反转:", lst[5:2:-1])


class Graph:
    """图结构实现"""

    def __init__(self):
        self.graph = {}

    def add_edge(self, u, v, directed=False):
        """添加边"""
        if u not in self.graph:
            self.graph[u] = []
        self.graph[u].append(v)

        if not directed:  # 无向图需要添加反向边
            if v not in self.graph:
                self.graph[v] = []
            self.graph[v].append(u)

    def bfs(self, start):
        """广度优先搜索"""
        visited = set()
        queue = [start]
        result = []

        while queue:
            vertex = queue.pop(0)
            if vertex not in visited:
                visited.add(vertex)
                result.append(vertex)
                queue.extend(self.graph.get(vertex, []))
        return result

    def dfs(self, start):
        """深度优先搜索"""
        visited = set()
        result = []

        def _dfs(node):
            visited.add(node)
            result.append(node)
            for neighbor in self.graph.get(node, []):
                if neighbor not in visited:
                    _dfs(neighbor)

        _dfs(start)
        return result


# ==================== 常用算法增强版 ====================

class DPAlgorithms:
    """动态规划算法"""

    @staticmethod
    def fibonacci(n):
        """斐波那契数列(DP版)"""
        if n == 0:
            return 0
        dp = [0] * (n + 1)
        dp[1] = 1
        for i in range(2, n + 1):
            dp[i] = dp[i - 1] + dp[i - 2]
        return dp[n]

    @staticmethod
    def knapsack(weights, values, capacity):
        """0-1背包问题"""
        n = len(weights)
        dp = [[0] * (capacity + 1) for _ in range(n + 1)]

        for i in range(1, n + 1):
            for w in range(1, capacity + 1):
                if weights[i - 1] <= w:
                    dp[i][w] = max(dp[i - 1][w], values[i - 1] + dp[i - 1][w - weights[i - 1]])
                else:
                    dp[i][w] = dp[i - 1][w]
        return dp[n][capacity]


class DivideAndConquer:
    """分治算法"""

    @staticmethod
    def merge_sort(arr):
        """归并排序"""
        if len(arr) <= 1:
            return arr

        mid = len(arr) // 2
        left = DivideAndConquer.merge_sort(arr[:mid])
        right = DivideAndConquer.merge_sort(arr[mid:])

        return DivideAndConquer.merge(left, right)

    @staticmethod
    def merge(left, right):
        """合并两个有序数组"""
        result = []
        i = j = 0

        while i < len(left) and j < len(right):
            if left[i] < right[j]:
                result.append(left[i])
                i += 1
            else:
                result.append(right[j])
                j += 1

        result.extend(left[i:])
        result.extend(right[j:])
        return result


class GreedyAlgorithms:
    """贪心算法"""

    @staticmethod
    def coin_change(coins, amount):
        """硬币找零问题"""
        coins.sort(reverse=True)
        count = 0
        result = []

        for coin in coins:
            while amount >= coin:
                amount -= coin
                count += 1
                result.append(coin)

        return result if amount == 0 else []


# ==================== 实用算法模板增强版 ====================

class AdvancedTemplates:
    """进阶算法模板"""

    @staticmethod
    def dijkstra(graph, start):
        """Dijkstra最短路径算法"""
        import heapq
        distances = {vertex: float('inf') for vertex in graph}
        distances[start] = 0
        heap = [(0, start)]

        while heap:
            current_dist, current_vertex = heapq.heappop(heap)

            if current_dist > distances[current_vertex]:
                continue

            for neighbor, weight in graph[current_vertex].items():
                distance = current_dist + weight
                if distance < distances[neighbor]:
                    distances[neighbor] = distance
                    heapq.heappush(heap, (distance, neighbor))

        return distances

    @staticmethod
    def topological_sort(graph):
        """拓扑排序(Kahn算法)"""
        in_degree = {u: 0 for u in graph}
        for u in graph:
            for v in graph[u]:
                in_degree[v] += 1

        queue = [u for u in graph if in_degree[u] == 0]
        result = []

        while queue:
            u = queue.pop(0)
            result.append(u)
            for v in graph[u]:
                in_degree[v] -= 1
                if in_degree[v] == 0:
                    queue.append(v)

        if len(result) == len(graph):
            return result
        return []  # 存在环


# ==================== 算法实战案例 ====================

class AlgorithmCases:
    """算法实战案例"""

    @staticmethod
    def lru_cache():
        """LRU缓存实现"""
        from collections import OrderedDict

        class LRUCache:
            def __init__(self, capacity):
                self.cache = OrderedDict()
                self.capacity = capacity

            def get(self, key):
                if key not in self.cache:
                    return -1
                self.cache.move_to_end(key)
                return self.cache[key]

            def put(self, key, value):
                if key in self.cache:
                    self.cache.move_to_end(key)
                self.cache[key] = value
                if len(self.cache) > self.capacity:
                    self.cache.popitem(last=False)

        # 测试用例
        cache = LRUCache(2)
        cache.put(1, 1)
        cache.put(2, 2)
        print("LRU Get 1:", cache.get(1))
        cache.put(3, 3)
        print("LRU Get 2:", cache.get(2))

    @staticmethod
    def word_break(s, word_dict):
        """单词拆分问题(DP解法)"""
        n = len(s)
        dp = [False] * (n + 1)
        dp[0] = True

        for i in range(1, n + 1):
            for j in range(i):
                if dp[j] and s[j:i] in word_dict:
                    dp[i] = True
                    break
        return dp[n]


# ==================== 测试代码 ====================

if __name__ == "__main__":
    print("===== 列表进阶操作 =====")
    AdvancedList.list_comprehension()
    AdvancedList.slicing_techniques()

    print("\n===== 图结构操作 =====")
    g = Graph()
    g.add_edge(0, 1)
    g.add_edge(0, 2)
    g.add_edge(1, 3)
    g.add_edge(2, 4)
    print("BFS遍历:", g.bfs(0))
    print("DFS遍历:", g.dfs(0))

    print("\n===== 动态规划算法 =====")
    print("斐波那契(10):", DPAlgorithms.fibonacci(10))
    weights = [2, 3, 4, 5]
    values = [3, 4, 5, 6]
    print("背包问题(容量=5):", DPAlgorithms.knapsack(weights, values, 5))

    print("\n===== 分治算法 =====")
    arr = [38, 27, 43, 3, 9, 82, 10]
    print("归并排序:", DivideAndConquer.merge_sort(arr))

    print("\n===== 贪心算法 =====")
    coins = [1, 5, 10, 25]
    print("硬币找零(63分):", GreedyAlgorithms.coin_change(coins, 63))

    print("\n===== 算法实战 =====")
    AlgorithmCases.lru_cache()
    print("单词拆分:", AlgorithmCases.word_break("leetcode", {"leet", "code"}))