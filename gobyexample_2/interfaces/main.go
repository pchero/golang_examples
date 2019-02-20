// interfaces are named collections of method
// signatures

package main

import (
	"fmt"
	"math"
)

// here's a basic interface for geometric shapes.
type geometry interface {
	area() float64
	perim() float64
}

// for our example we'll implement this interface on
// 'rect' and 'circle' types.
type rect struct {
	width  float64
	height float64
}
type circle struct {
	radius float64
}

// to implement an interface in go, we just need to
// implement all the methods in the interface. here we
// implement 'geometry' on 'rect's.
func (r rect) area() float64 {
	return r.width * r.height
}

func (r rect) perim() float64 {
	return 2*r.width + 2*r.height
}

// the implementation for 'circle's
func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}

// if a variable has an interface type, then we can call
// methods that are in the named interface. here's a
// generic 'measure' function taking advantage of this
// to work on any 'geometry'.
func measure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perim())
}

func main() {
	r := rect{width: 3, height: 4}
	c := circle{radius: 5}

	// the 'circle' and 'rect' struct types both
	// implement the 'geometry' interface so we can use
	// instance of
	// these structs as arguments to 'measure'.
	measure(r)
	measure(c)
}
