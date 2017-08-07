package main

import "fmt"

type puzzle map[string]string

func (p puzzle) String() string {
	var str string
	for y := 0; y < N; y++ {
		for z := 0; z < N; z++ {
			for x := 0; x < N; x++ {
				str += p.Get(newPoint(x, y, z))
			}
			str += " "
		}
		str += "\n"
	}
	return str
}

func (p puzzle) Get(pt *point) string {
	return p[p.Key(pt)]
}

func (p puzzle) Set(pt *point, val string) {
	p[p.Key(pt)] = val
}

func (p puzzle) Key(pt *point) string {
	return fmt.Sprintf("x%dy%dz%d", pt.x, pt.y, pt.z)
}
