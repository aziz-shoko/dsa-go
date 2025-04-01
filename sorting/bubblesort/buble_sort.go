// package bubblesort provides an implementation of the buble sort algorithm
package bubblesort

import (
	"golang.org/x/exp/constraints" 	// constraints is an experimental package from golang x/exp/ and is used for writing generic code in Go
)

// Sort perofrms  an in-place buble srot on the provided slice.
// It returns the sorted slice for convenience
// Time Complexity: O(n^2) where n is hte length of the slice
// Space complexity: O(1) as sorting is done in-place
func Sort[T constraints.Ordered](items []T) []T {
	n := len(items)
	if n <= 1 {
		return items
	}

	swapped := true
	// Continue until we go through the entire slice without swapping
	// [3,1,4,5,2]
	for swapped {
		swapped = false
		for i:=0; i < n-1; i++ {
			if items[i] > items[i+1] {
				items[i], items[i+1] = items[i+1], items[i]
				swapped = true
			}
		}
		// After each iteration, the largest element has bubbled to the end
		// So we can reduce the range we check by 1
		n--
	}

	return items
}

// SortWithComparator sorts the slice using a custom comparison function
// The comparator function should return:
// - negative value if a < b
// - zero if a == b
// - positive value if a > b
// This allows sorting of complex types or custom orderring
func SortWithComparator[T any](items []T, comparator func(a, b T) int) []T {
	n := len(items)	
	if n <= 1 {
		return items
	}

	swapped := true
	for swapped {
		swapped = false
		for i := 0; i < n -1; i++ {
			if comparator(items[i], items[i+1]) > 0 {
				items[i], items[i+1] = items[i+1], items[i]
				swapped = true
			}
		}
		n--
	}

	return items
}

