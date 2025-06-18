package structsmethodsinterfaces

import (
	"math"
)

// Interface shape
type Shape interface {
	Area() float64
}

// Rectangle methods
func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Height + r.Width)
}

func (r *Rectangle) Area() float64 {
	return r.Height * r.Width
}

// Circle methods
func (r *Circle) Perimeter() float64 {
	return 2 * r.Radius * math.Pi
}

func (r *Circle) Area() float64 {
	return math.Pi * r.Radius * r.Radius
}

// Triangle methods
func (t *Triangle) Area() float64 {
	return (t.Base * t.Height) / 2
}