package main

import (
	"fmt"
	"math"
)

type rectangle struct {
	width, height float64
}

type circle struct {
	radius float64
}

type geometry interface {
	area() float64
	perim() float64
}

func (r rectangle) perim() float64 {
	return 2*r.width + 2*r.height
}

func (r rectangle) area() float64 {
	return r.width * r.height
}

func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func measure(g geometry) {
	fmt.Println("area", g.area())
	fmt.Println("perim", g.perim())
}

func main() {
	rq := rectangle{10, 5}
	qp := circle{22}

	measure(rq)
	measure(qp)
}
