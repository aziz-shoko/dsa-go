```go
package bubblesort_test

import (
	"fmt"
	"github.com/<username>/dsa-go/sorting/bubblesort"
)

// This example dmeonstrates how to use the Sort function
// with different types of slices
func Example_basic() {
	// Sort integers 
	numbers := []int{5, 2, 6, 3, 1, 4}
	sorted := bubblesort.Sort(numbers)
	fmt.Println(sorted)

	// Sort strings
	words := []string{"banana", "apple", "cherry", "date"}
	sortedWords := bubblesort.Sort(words)
	fmt.Println(sortedWords)
	
	// Output:
	// [1 2 3 4 5 6]
	// [apple banana cherry date]
}

// This example shows how to use the SortWithComparator function
// to sort custom types or with custom ordering.
func Example_withComparator() {
	// Define a custom type
	type Person struct {
		Name string
		Age  int
	}
	
	// Create a slice of our custom type
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
		{"David", 20},
	}
	
	// Sort by age (ascending)
	byAge := bubblesort.SortWithComparator(people, func(a, b Person) int {
		return a.Age - b.Age
	})
	
	// Print the sorted slice
	for _, p := range byAge {
		fmt.Printf("%s: %d years\n", p.Name, p.Age)
	}
	
	// Output:
	// David: 20 years
	// Bob: 25 years
	// Alice: 30 years
	// Charlie: 35 years
}
```