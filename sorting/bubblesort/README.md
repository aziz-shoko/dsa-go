// Disclaimer: This bubblesort package was done by ClaudeAI as a way for me to see how packages in Go or in programming general work.
// The other packages will be implemented by me (which will be shit but its for learning purposes)
# Bubble Sort

A generic implementation of the bubble sort algorithm in Go.

## Description

Bubble sort is a simple comparison-based sorting algorithm. It repeatedly steps through the list, compares adjacent elements, and swaps them if they are in the wrong order. The pass through the list is repeated until the list is sorted.

### Characteristics:

- **Time Complexity**: O(n²) in worst and average cases, O(n) in best case (when the list is already sorted)
- **Space Complexity**: O(1) as sorting is done in-place
- **Stable**: Yes (equal elements maintain their relative order)

## Usage

```go
import "github.com/yourusername/dsa-go/sorting/bubblesort"

// Sort a slice of integers
numbers := []int{5, 2, 6, 3, 1, 4}
sorted := bubblesort.Sort(numbers)
// sorted: [1, 2, 3, 4, 5, 6]

// Sort a slice of strings
words := []string{"banana", "apple", "cherry", "date"}
sortedWords := bubblesort.Sort(words)
// sortedWords: ["apple", "banana", "cherry", "date"]

// Sort custom types using a comparator function
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {"Alice", 30},
    {"Bob", 25},
    {"Charlie", 35},
    {"David", 20},
}

// Sort by age
byAge := bubblesort.SortWithComparator(people, func(a, b Person) int {
    return a.Age - b.Age
})
// byAge: [{"David", 20}, {"Bob", 25}, {"Alice", 30}, {"Charlie", 35}]
```

## Features

- Generic implementation works with any ordered type (numbers, strings, etc.)
- Custom comparator function for complex types or custom ordering
- Efficient implementation with early termination when the list becomes sorted

## Testing

Run tests with:

```bash
go test
```

Run benchmarks with:

```bash
go test -bench=.
```
