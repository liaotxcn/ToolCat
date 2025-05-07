"""
Python 数据结构与算法(1)
包含：基础数据结构实现、常用算法模板和示例
"""


# ==================== Part 1: 基础数据结构 ====================

class ListOperations:
    """列表(list)操作大全"""

    @staticmethod
    def demo():
        # 初始化
        lst = [12, 56, 43, 23, 'abc', 87]
        print(f"原始列表: {lst}")

        # 增删改查
        lst.append(34)  # 末尾添加
        lst.insert(1, 'new')  # 指定位置插入
        lst.remove(43)  # 删除第一个匹配项
        popped = lst.pop(2)  # 删除并返回指定位置元素
        lst[0] = 'updated'  # 修改元素

        # 其他操作
        lst.reverse()  # 反转列表
        lst_copy = lst.copy()  # 浅拷贝
        lst.clear()  # 清空列表

        # 遍历方式
        for item in lst:
            print(item, end=' ')
        print()

        for i in range(len(lst)):
            print(lst[i], end=' ')
        print()

        for i, val in enumerate(lst):
            print(f"Index {i}: {val}")


class StringOperations:
    """字符串操作大全"""

    @staticmethod
    def demo():
        s = "  Hello World!  "

        # 常用方法
        print(s.strip())  # "Hello World!"
        print(s.lower())  # "  hello world!  "
        print(s.upper())  # "  HELLO WORLD!  "
        print(s.replace('l', 'X'))  # "  HeXXo WorXd!  "
        print(s.split())  # ['Hello', 'World!']

        # 判断方法
        print("123".isdigit())  # True
        print("abc".isalpha())  # True

        # 字符串格式化
        name, age = "Alice", 25
        print(f"{name} is {age} years old")  # f-string (Python 3.6+)


class LinkedList:
    """链表实现"""

    class Node:
        def __init__(self, val):
            self.val = val
            self.next = None

    def __init__(self):
        self.head = None

    def append(self, val):
        if not self.head:
            self.head = self.Node(val)
        else:
            curr = self.head
            while curr.next:
                curr = curr.next
            curr.next = self.Node(val)

    def print_list(self):
        curr = self.head
        while curr:
            print(curr.val, end=" -> ")
            curr = curr.next
        print("None")


class Stack:
    """栈实现 (LIFO)"""

    def __init__(self):
        self.items = []

    def push(self, item):
        self.items.append(item)

    def pop(self):
        if not self.is_empty():
            return self.items.pop()

    def peek(self):
        if not self.is_empty():
            return self.items[-1]

    def is_empty(self):
        return len(self.items) == 0

    def size(self):
        return len(self.items)


class Queue:
    """队列实现 (FIFO)"""

    def __init__(self):
        self.items = []

    def enqueue(self, item):
        self.items.insert(0, item)

    def dequeue(self):
        if not self.is_empty():
            return self.items.pop()

    def is_empty(self):
        return len(self.items) == 0

    def size(self):
        return len(self.items)


class TreeNode:
    """二叉树节点"""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class BinaryTree:
    """二叉树遍历"""

    @staticmethod
    def preorder(root):
        """前序遍历: 根->左->右"""
        return [root.val] + BinaryTree.preorder(root.left) + BinaryTree.preorder(root.right) if root else []

    @staticmethod
    def inorder(root):
        """中序遍历: 左->根->右"""
        return BinaryTree.inorder(root.left) + [root.val] + BinaryTree.inorder(root.right) if root else []

    @staticmethod
    def postorder(root):
        """后序遍历: 左->右->根"""
        return BinaryTree.postorder(root.left) + BinaryTree.postorder(root.right) + [root.val] if root else []


# ==================== Part 2: 常用算法 ====================

class SortingAlgorithms:
    """排序算法"""

    @staticmethod
    def bubble_sort(arr):
        """冒泡排序 O(n^2)"""
        n = len(arr)
        for i in range(n):
            for j in range(0, n - i - 1):
                if arr[j] > arr[j + 1]:
                    arr[j], arr[j + 1] = arr[j + 1], arr[j]
        return arr

    @staticmethod
    def quick_sort(arr):
        """快速排序 O(n log n)"""
        if len(arr) <= 1:
            return arr
        pivot = arr[len(arr) // 2]
        left = [x for x in arr if x < pivot]
        middle = [x for x in arr if x == pivot]
        right = [x for x in arr if x > pivot]
        return SortingAlgorithms.quick_sort(left) + middle + SortingAlgorithms.quick_sort(right)


class SearchAlgorithms:
    """搜索算法"""

    @staticmethod
    def binary_search(arr, target):
        """二分查找 O(log n)"""
        left, right = 0, len(arr) - 1
        while left <= right:
            mid = left + (right - left) // 2
            if arr[mid] == target:
                return mid
            elif arr[mid] < target:
                left = mid + 1
            else:
                right = mid - 1
        return -1


class TwoPointers:
    """双指针技巧"""

    @staticmethod
    def has_cycle(head):
        """快慢指针判断链表是否有环"""
        slow = fast = head
        while fast and fast.next:
            slow = slow.next
            fast = fast.next.next
            if slow == fast:
                return True
        return False

    @staticmethod
    def two_sum(nums, target):
        """对撞指针求两数之和"""
        left, right = 0, len(nums) - 1
        while left < right:
            s = nums[left] + nums[right]
            if s == target:
                return [left, right]
            elif s < target:
                left += 1
            else:
                right -= 1
        return []


# ==================== Part 3: 算法模板 ====================

class AlgorithmTemplates:
    """常用算法模板"""

    @staticmethod
    def sliding_window(s, k):
        """滑动窗口模板"""
        from collections import defaultdict
        window = defaultdict(int)
        left = right = 0
        res = 0

        while right < len(s):
            # 右移窗口
            c = s[right]
            window[c] += 1
            right += 1

            # 满足条件时收缩窗口
            while len(window) > k:
                d = s[left]
                window[d] -= 1
                if window[d] == 0:
                    del window[d]
                left += 1

            # 更新结果
            res = max(res, right - left)
        return res

    @staticmethod
    def backtrack_template(nums):
        """回溯法模板"""
        res = []

        def backtrack(path, choices):
            if 满足结束条件:
                res.append(path.copy())
                return

            for choice in choices:
                if 选择不合法:
                    continue

                # 做选择
                path.append(choice)

                # 进入下一层决策
                backtrack(path, 新选择列表)

                # 撤销选择
                path.pop()

        backtrack([], nums)
        return res


# ==================== 测试代码 ====================

if __name__ == "__main__":
    print("===== 列表操作演示 =====")
    ListOperations.demo()

    print("\n===== 字符串操作演示 =====")
    StringOperations.demo()

    print("\n===== 链表操作演示 =====")
    ll = LinkedList()
    for i in [1, 2, 3, 4]:
        ll.append(i)
    ll.print_list()

    print("\n===== 排序算法演示 =====")
    arr = [64, 34, 25, 12, 22, 11, 90]
    print("冒泡排序:", SortingAlgorithms.bubble_sort(arr.copy()))
    print("快速排序:", SortingAlgorithms.quick_sort(arr.copy()))

    print("\n===== 二分查找演示 =====")
    sorted_arr = [11, 12, 22, 25, 34, 64, 90]
    target = 22
    print(f"在{sorted_arr}中查找{target}: 索引={SearchAlgorithms.binary_search(sorted_arr, target)}")

    print("\n===== 双指针演示 =====")
    nums = [2, 7, 11, 15]
    target = 9
    print(f"两数之和{nums}, target={target}: {TwoPointers.two_sum(nums, target)}")