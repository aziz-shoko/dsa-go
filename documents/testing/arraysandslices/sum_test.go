package arraysandslices

import (
	"fmt"
	"slices"
	"testing"
)

func TestSum(t *testing.T) {
	t.Run("collection of 5 numbers", func(t *testing.T) {
		given := []int{1, 2, 3, 4, 5}
		sum := Sum(given)
		expected := 15

		if sum != expected {
			t.Errorf("given %v, expected %d but got %d", given, expected, sum)
		}
	})
}

func BenchmarkRepeat(b *testing.B) {
	given := []int{1, 2, 3, 4, 5}

	for b.Loop() {
		Sum(given)
	}
}

func ExampleSum() {
	given := []int{1, 2, 3, 4, 5}
	fmt.Printf("%d", Sum(given))
	// Output: 15
}

func TestSumAll(t *testing.T) {
	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if !slices.Equal(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestSumAllTails(t *testing.T) {
	checkSums := func(t testing.TB, got, want []int) {
		t.Helper()
		if !slices.Equal(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}

	t.Run("sum of two populated slices", func(t *testing.T) {
		got := SumAllTails([]int{1, 2, 3}, []int{4, 5, 6})
		want := []int{5, 11}
		checkSums(t, got, want)
	})

	t.Run("sum of empty slices", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{4, 5, 6})
		want := []int{0, 11}
		checkSums(t, got, want)
	})
}
