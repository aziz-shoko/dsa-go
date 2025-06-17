package integer

import (
	"testing"
	"fmt"
)

func TestAddr(t *testing.T) {
	sum := Add(3, 8)
	expected := 11

	if sum != expected {
		t.Errorf("Expected %d but got %d", expected, sum)
	}
}

func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	// Output: 6
}

