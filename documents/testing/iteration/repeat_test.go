package iteration

import (
	"testing"
	"fmt"
)

func TestRepeat(t *testing.T) {
	repeated := Repeat("a", 6)
	expected := "aaaaaa"

	if repeated != expected {
		t.Errorf("repeated %q expected %q", repeated, expected)
	}
}

// func Benchmark(b *testing.B) {
// 	//... setup ...
// 	for b.Loop() {
// 		//... code to measure ...
// 	}
// 	//... cleanup ...
// }
func BenchmarkRepeat(b *testing.B) {
	for b.Loop() {
		Repeat("a", 5)
	}	
}

func ExampleRepeat() {
	repeated := Repeat("a", 3)
	fmt.Printf("%q", repeated)
	// Output: "aaa"
}