package code

// package for a basic library
const LIB string = `
package mylib

import (
	"fmt"
	"math"
)

// Shape is an interface that defines the methods a shape should have
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle is a struct that implements the Shape interface
type Rectangle struct {
	Width  float64
	Height float64
}

// Area returns the area of a rectangle
func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter returns the perimeter of a rectangle
func (r *Rectangle) Perimeter() float64 {
	return 2*r.Width + 2*r.Height
}

// Circle is a struct that implements the Shape interface
type Circle struct {
	Radius float64
}

// Area returns the area of a circle
func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter returns the perimeter of a circle
func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// PrintAreaAndPerimeter takes a Shape as an argument and prints its area and perimeter
func PrintAreaAndPerimeter(s Shape) {
	fmt.Println("Area:", s.Area())
	fmt.Println("Perimeter:", s.Perimeter())
}
`
