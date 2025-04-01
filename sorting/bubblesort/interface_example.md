```Go
// In GO, only package mains can be directly executed by the cli, or else it will error:
// package command-line-arguments is not a main package
// package main

import (
	"fmt"
)

// This interface defines what it means to be "sortable"
type Sortable interface {
	Less(i, j int) bool
	Swap(i, j int)
	Length() int
}

// This function can sort ANYTHING that satisfies the Sortable interface
func BubbleSort(s Sortable) {
	for i := 0; i < s.Length(); i++ {
		for j := 0; j < s.Length()-i-1; j++ {
			if s.Less(j+1, j) {
				s.Swap(j, j+1)
			}
		}
	}
}

// let's create a concrete type: a slice of ints
type IntSlice []int

// Implement the Sortable interface for IntSlice
func (s IntSlice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s IntSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s IntSlice) Length() int {
	return len(s)
}

// Now we can also create a complete different type
type PersonSlice []struct {
	Name string
	Age int
}

// Implement the sortable interface for PersonSlice
func (s PersonSlice) Less(i, j int) bool {
	return s[i].Age < s[j].Age	// Sort by age
}

func (s PersonSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s PersonSlice) Length() int {
	return len(s)
}

// Usage
func main() {
	numbers := IntSlice{5,2,6,3,1,4}
	BubbleSort(numbers)
	fmt.Println(numbers)

	people := PersonSlice{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}
	BubbleSort(people)
	fmt.Println(people)
}
```