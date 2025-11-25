# Go Heap Mastery Guide - 15 Minutes to Forever

## 🎯 THE GOLDEN TEMPLATE (Memorize This!)

```go
import "container/heap"

type MyHeap []ElementType

// The Magic 5 Methods (ALWAYS the same structure)
func (h MyHeap) Len() int           { return len(h) }
func (h MyHeap) Less(i, j int) bool { return h[i] < h[j] } // ⭐ KEY: This line determines min/max
func (h MyHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MyHeap) Push(x interface{}) { *h = append(*h, x.(ElementType)) }
func (h *MyHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

// Usage Pattern
h := &MyHeap{}
heap.Init(h)              // Only needed if heap has initial data
heap.Push(h, value)       // Add elements
top := heap.Pop(h)        // Remove and get top element
peek := (*h)[0]           // Peek at top without removing
```

## 🧠 Memory Tricks

### The "Less is More" Rule
**If `Less(i, j)` returns true, element `i` goes HIGHER in the heap**

- **Min Heap:** `h[i] < h[j]` → smaller values rise to top
- **Max Heap:** `h[i] > h[j]` → larger values rise to top

### The 5-Method Mantra
1. **Len** - how many?
2. **Less** - who's higher? (this determines min/max)
3. **Swap** - trade places
4. **Push** - add to end (heap will fix itself)
5. **Pop** - remove from end and return it

## 📚 Common Patterns & Examples

### Pattern 1: Simple Integer Heap
```go
// Min heap of integers
type IntHeap []int
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }

// Max heap of integers  
type IntHeap []int
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }
```

### Pattern 2: Struct with Custom Priority
```go
type Task struct {
    name     string
    priority int
}

// Max heap by priority (higher priority = higher in heap)
type TaskHeap []Task
func (h TaskHeap) Less(i, j int) bool { return h[i].priority > h[j].priority }

// Min heap by priority (lower priority = higher in heap)
type TaskHeap []Task  
func (h TaskHeap) Less(i, j int) bool { return h[i].priority < h[j].priority }
```

### Pattern 3: Multiple Criteria Sorting
```go
type Student struct {
    name  string
    grade int
    age   int
}

// Sort by grade (descending), then by age (ascending) as tiebreaker
type StudentHeap []Student
func (h StudentHeap) Less(i, j int) bool {
    if h[i].grade != h[j].grade {
        return h[i].grade > h[j].grade  // Higher grade first
    }
    return h[i].age < h[j].age          // Younger first for ties
}
```

## 🎯 LeetCode Problem Templates

### Template 1: "Top K" Problems
```go
func topKSomething(arr []int, k int) []int {
    h := &MinHeap{}  // Use min heap of size k
    heap.Init(h)
    
    for _, val := range arr {
        heap.Push(h, val)
        if h.Len() > k {
            heap.Pop(h)  // Remove smallest
        }
    }
    
    result := make([]int, k)
    for i := k-1; i >= 0; i-- {
        result[i] = heap.Pop(h).(int)
    }
    return result
}
```

### Template 2: "Find Median" Problems
```go
type MedianFinder struct {
    maxHeap *MaxHeap  // Stores smaller half
    minHeap *MinHeap  // Stores larger half
}

func (m *MedianFinder) AddNum(num int) {
    if m.maxHeap.Len() == 0 || num <= (*m.maxHeap)[0] {
        heap.Push(m.maxHeap, num)
    } else {
        heap.Push(m.minHeap, num)
    }
    
    // Balance the heaps
    if m.maxHeap.Len() > m.minHeap.Len() + 1 {
        heap.Push(m.minHeap, heap.Pop(m.maxHeap))
    } else if m.minHeap.Len() > m.maxHeap.Len() + 1 {
        heap.Push(m.maxHeap, heap.Pop(m.minHeap))
    }
}
```

### Template 3: "Merge K" Problems
```go
func mergeKSortedLists(lists []*ListNode) *ListNode {
    h := &NodeHeap{}
    heap.Init(h)
    
    // Add first node from each list
    for _, list := range lists {
        if list != nil {
            heap.Push(h, list)
        }
    }
    
    dummy := &ListNode{}
    current := dummy
    
    for h.Len() > 0 {
        node := heap.Pop(h).(*ListNode)
        current.Next = node
        current = current.Next
        
        if node.Next != nil {
            heap.Push(h, node.Next)
        }
    }
    
    return dummy.Next
}
```

## ⚡ Quick Reference Card

| Want | Less() Method | Example |
|------|---------------|---------|
| **Min Heap** | `h[i] < h[j]` | Smallest number on top |
| **Max Heap** | `h[i] > h[j]` | Largest number on top |
| **Top K Largest** | Use Min Heap of size K | Keep largest K elements |
| **Top K Smallest** | Use Max Heap of size K | Keep smallest K elements |

## 🚨 Common Mistakes & How to Avoid Them

### ❌ DON'T DO THIS:
```go
h.Push(x)           // Wrong! This calls your method directly
h.Pop()             // Wrong! Doesn't maintain heap property
x := h[0]; h = h[1:] // Wrong! Breaks heap structure
```

### ✅ DO THIS:
```go
heap.Push(h, x)     // Right! Uses heap package
heap.Pop(h)         // Right! Maintains heap property  
x := (*h)[0]        // Right! Peek without modifying
```

### ❌ Forgetting Type Assertions:
```go
val := heap.Pop(h)  // val is interface{}
```
### ✅ Proper Type Assertion:
```go
val := heap.Pop(h).(int)  // Convert to actual type
```

## 🎯 15-Minute Practice Drill

**Memorize the template by writing it 3 times without looking:**

1. **Write the 5 methods** for `type IntHeap []int` (min heap)
2. **Write the 5 methods** for a max heap of structs
3. **Write a complete solution** to "Kth Largest Element" using the template

**Time yourself:** Can you write a working heap solution to any "top k" problem in under 2 minutes?

## 🏆 Graduation Test

You've mastered Go heaps when you can:
- [ ] Write the 5 heap methods from memory in 30 seconds
- [ ] Instantly know: min heap uses `<`, max heap uses `>`
- [ ] Solve any "top k" problem using the template
- [ ] Remember: always use `heap.Push/Pop`, never call methods directly

**Congratulations! You now have a superpower that works on 90% of priority queue problems! 🎉**