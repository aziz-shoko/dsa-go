package bubblesort

import (
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "single element",
			input:    []int{1},
			expected: []int{1},
		},
		{
			name:     "already sorted",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "reverse sorted",
			input:    []int{5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "random order",
			input:    []int{3, 1, 4, 1, 5, 9, 2, 6, 5},
			expected: []int{1, 1, 2, 3, 4, 5, 5, 6, 9},
		},
		{
			name:     "with duplicates",
			input:    []int{3, 1, 3, 1, 5, 5, 2},
			expected: []int{1, 1, 2, 3, 3, 5, 5},
		},
		{
			name:     "negative numbers",
			input:    []int{-3, -1, -4, 1, -5, 9, -2},
			expected: []int{-5, -4, -3, -2, -1, 1, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy to avoid modifying the test data
			input := make([]int, len(tt.input))
			copy(input, tt.input)

			result := Sort(input)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Sort() = %v, want %v", result, tt.expected)
			}
		})
	}

	// Test with other types
	t.Run("string slice", func(t *testing.T) {
		input := []string{"banana", "apple", "cherry", "date"}
		expected := []string{"apple", "banana", "cherry", "date"}

		result := Sort(input)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Sort()=%v, want %v", result, expected)
		}
	})

	t.Run("float slice", func(t *testing.T) {
		input := []float64{3.14, 1.41, 2.71, 1.73}
		expected := []float64{1.41, 1.73, 2.71, 3.14}

		result := Sort(input)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Sort()=%v, want %v", result, expected)
		}
	})
}

func TestSortWithComparator(t *testing.T) {
	type Person struct {
		Name string
		Age int
	}

	people := []Person{
		{"Alice",30},
		{"Bob",25},
		{"Charlie",35},
		{"David",20},
	}

	t.Run("sort by age", func(t *testing.T) {
		expected := []Person{
			{"David", 20},
			{"Bob", 25},
			{"Alice", 30},
			{"Charlie", 35},
		}

		result := SortWithComparator(people, func(a, b Person) int {
			return a.Age - b.Age
		})

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("SortWithComparator()=%v, want %v", result, expected)
		}
	})

	t.Run("sort by name", func(t *testing.T) {
		expected := []Person{
			{"Alice",30},
			{"Bob", 25},
			{"Charlie", 35},
			{"David", 20},
		}

		result := SortWithComparator(people, func(a, b Person) int {
			if a.Name < b.Name {
				return -1
			} else if a.Name > b.Name {
				return 1
			} else {
				return 0
			}
		})

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("SortWithComparator()=%v, want %v", result, expected)
		}
	})
}

func BenchmarkXxx(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run("size="+string(rune(size)), func(b *testing.B) {
			// Create a worse-case scenario (reverse sorted)
			input := make([]int, size)
			for i:=0; i < size; i++ {
				input[i] = size - i
			}

			b.ResetTimer()
			for i:=0; i < b.N; i++ {
				// Make a copy so we don't benefit from previous sorts
				data := make([]int, len(input))
				copy(data, input)
				Sort(data)
			}
		})
	}
}