package main

import "sort"

type sixlet struct {
	p1 *point
	p2 *point
	p3 *point
	p4 *point
	p5 *point
	p6 *point
}

// newSixlet takes the 6 points and the plane the sixlet lies on, and ensures
// the points are sorted properly
func newSixlet(p1, p2, p3, p4, p5, p6 *point, si sorter) *sixlet {
	pts := []*point{p1, p2, p3, p4, p5, p6}
	sort.Sort(si(pts))
	return &sixlet{
		p1: pts[0],
		p2: pts[1],
		p3: pts[2],
		p4: pts[3],
		p5: pts[4],
		p6: pts[5],
	}
}

func (sl *sixlet) Points() []*point {
	return []*point{sl.p1, sl.p2, sl.p3, sl.p4, sl.p5, sl.p6}
}

func (sl *sixlet) String() string {
	var str string
	for i, pt := range sl.Points() {
		if i%2 == 0 {
			str += "\n"
		}
		str += pzl.Get(pt)
	}
	return str
}
