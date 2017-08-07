package main

import "fmt"

type pair struct {
	p1 *point
	p2 *point
}

func newPair(p1, p2 *point) *pair {
	return &pair{p1: p1, p2: p2}
}

func (p *pair) String() string {
	return fmt.Sprintf("p1: %s, p2: %s\n", p.p1, p.p2)
}
