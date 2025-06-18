package structsmethodsinterfaces

import (
	"testing"
)

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := rectangle.Perimeter()
	want := 40.0
	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

// func TestArea(t *testing.T) {
// 	checkArea := func(t testing.TB, shape Shape, want float64) {
// 		t.Helper()
// 		got := shape.Area()
// 		if got != want {
// 			t.Errorf("got %g want %g", got, want)
// 		}
// 	}
// 	t.Run("rectangles", func(t *testing.T) {
// 		rectangle := Rectangle{10.0, 10.0}
// 		checkArea(t, &rectangle, 100.0)
// 	})

// 	t.Run("circle", func(t *testing.T) {
// 		circle := Circle{10.0}
// 		checkArea(t, &circle, 314.1592653589793)
// 	})
// }

// Table Driven Test example below
func TestArea(t *testing.T) {
	areaTests := []struct {
		name  string
		shape Shape
		want  float64
	}{
		{name: "Rectangle", shape: &Rectangle{10.0, 10.0}, want: 100.0},
		{name: "Circle", shape: &Circle{10}, want: 314.1592653589793},
		{name: "Triangle", shape: &Triangle{10.0, 5.0}, want: 25.0},
	}

	for _, tt := range areaTests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.want {
				t.Errorf("%#v got %g want %g", tt.shape, got, tt.want)
			}
		})
	}
}
