package main

import "fmt"

type point struct {
	x int
	y int
	z int
}

func newPoint(x, y, z int) *point {
	return &point{x: x, y: y, z: z}
}

func (p *point) String() string {
	return fmt.Sprintf("x: %d, y: %d, z: %d", p.x, p.y, p.z)
}

func (p *point) Coords() []int {
	return []int{p.x, p.y, p.z}
}
