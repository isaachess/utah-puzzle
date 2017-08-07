package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

var (
	pzl   = make(puzzle)
	X     = "X"
	O     = "O"
	Blank = "_"
	pairs []*pair
	N     int
)

func main() {
	for {
		if err := run(); err != nil {
			log.Print(err)
		}
	}
}

func run() error {

	//processLine("X__ ___ ___", 0)
	//processLine("XX_ _O_ ___", 1)
	//processLine("OXO ___ __X", 2)

	//processLine("___ OOO ___", 0)
	//processLine("___ _OX ___", 1)
	//processLine("___ _X_ _X_", 2)

	if err := processStdinLines(); err != nil {
		return err
	}

	// find sixlet we want by analyzing pairs
	sixletX, _ := findSixlet(pairs, X)
	sixletO, _ := findSixlet(pairs, O)

	best, err := bestMove(sixletX, sixletO)
	if err != nil {
		return err
	}

	pzl.Set(best, X)

	fmt.Println(pzl)

	return nil
}

func processStdinLines() error {
	var lines int
	for {
		r := bufio.NewReader(os.Stdin)
		line, _, err := r.ReadLine()
		if err != nil {
			return err
		}
		processLine(string(line), lines)

		lines++
		if lines == N {
			break
		}
	}

	return nil
}

func processLine(ln string, y int) {
	zLevels := strings.Split(ln, " ")

	if N == 0 {
		N = len(zLevels)
	}

	for z, zVal := range zLevels {
		for x := 0; x < len(zVal); x++ {
			char := string(zVal[x])
			pt := newPoint(x, y, z)
			pzl.Set(pt, char)
			if !isBlank(char) {
				pairs = append(pairs, findPairs(pt, char)...)
			}
		}
	}
}

func findPairs(pt *point, val string) []*pair {
	xs := validPairPointsForAxis(pt.x)
	ys := validPairPointsForAxis(pt.y)
	zs := validPairPointsForAxis(pt.z)
	return checkPairs(pt, xs, ys, zs, val)
}

func checkPairs(pt *point, xs, ys, zs []int, val string) []*pair {
	// check x pairs
	x, y, z := pt.x, pt.y, pt.z
	var pairs []*pair
	for _, newX := range xs {
		newP := newPoint(newX, y, z)
		if pointValMatchesVal(newP, val) {
			pairs = append(pairs, newPair(pt, newP))
		}
	}

	// check y pairs
	for _, newY := range ys {
		newP := newPoint(x, newY, z)
		if pointValMatchesVal(newP, val) {
			pairs = append(pairs, newPair(pt, newP))
		}
	}

	// check z pairs
	for _, newZ := range zs {
		newP := newPoint(x, y, newZ)
		if pointValMatchesVal(newP, val) {
			pairs = append(pairs, newPair(pt, newP))
		}
	}

	return pairs
}

func pointValMatchesVal(pt *point, val string) bool {
	return pzl.Get(pt) == val
}

func validPairPointsForAxis(val int) []int {
	var vals []int
	lessVal := val - 1
	moreVal := val + 1

	if lessVal >= 0 {
		vals = append(vals, lessVal)
	}
	if moreVal < N {
		vals = append(vals, moreVal)
	}

	return vals
}

func isBlank(val string) bool {
	return val == Blank
}

func findSixlet(pairs []*pair, team string) (*sixlet, error) {
	for _, p := range pairs {
		sl, err := findSixletForPair(p, team)
		if err == nil {
			return sl, nil
		}
	}

	return nil, errors.New("No sixlet found for any pair")
}

func findSixletForPair(p *pair, team string) (*sixlet, error) {
	if p.p1.x != p.p2.x {
		return varySixlet(p, team, modder{tran: pointModY, si: newByXY}, modder{tran: pointModZ, si: newByXZ})
	}
	if p.p1.y != p.p2.y {
		return varySixlet(p, team, modder{tran: pointModX, si: newByXY}, modder{tran: pointModZ, si: newByYZ})
	}
	if p.p1.z != p.p2.z {
		return varySixlet(p, team, modder{tran: pointModX, si: newByXZ}, modder{tran: pointModY, si: newByYZ})
	}

	return nil, errors.New("Pair is badly formed, cannot determine what transforms to do")
}

func varySixlet(p *pair, team string, mods ...modder) (*sixlet, error) {
	for _, mod := range mods {
		// pair with -1 and -2
		sltop := modSixlet(p.p1, p.p2, -1, -2, mod.tran, mod.si)

		// pair with -1 and +1
		slmid := modSixlet(p.p1, p.p2, -1, 1, mod.tran, mod.si)

		// pair with +1 and +2
		slbottom := modSixlet(p.p1, p.p2, 1, 2, mod.tran, mod.si)

		for _, sl := range []*sixlet{sltop, slmid, slbottom} {
			if validSixlet(sl) && sumSixlet(sl, team) >= 4 {
				return sl, nil
			}
		}
	}

	return nil, errors.New("No suitable sixlet found")
}

func modSixlet(p1, p2 *point, mod1, mod2 int, trans transformer, si sorter) *sixlet {
	return newSixlet(
		p1,
		p2,
		trans(p1, mod1),
		trans(p2, mod1),
		trans(p1, mod2),
		trans(p2, mod2),
		si,
	)
}

func pointModX(pt *point, mod int) *point {
	return newPoint(pt.x+mod, pt.y, pt.z)
}

func pointModY(pt *point, mod int) *point {
	return newPoint(pt.x, pt.y+mod, pt.z)
}

func pointModZ(pt *point, mod int) *point {
	return newPoint(pt.x, pt.y, pt.z+mod)
}

func validSixlet(sl *sixlet) bool {
	for _, pt := range sl.Points() {
		for _, val := range pt.Coords() {
			if val < 0 || val >= N {
				return false
			}
		}
	}

	return true
}

func sumSixlet(sl *sixlet, team string) int {
	var sum int
	for _, pt := range sl.Points() {
		if pointValMatchesVal(pt, team) {
			sum++
		}
	}
	return sum
}

func bestMove(slx, slo *sixlet) (*point, error) {
	if slx != nil {
		for _, pt := range slx.Points() {
			if isWinningPoint(pt, slx, X) {
				return pt, nil
			}
		}
	}

	if slo != nil {
		for _, pt := range slo.Points() {
			if isWinningPoint(pt, slo, O) {
				return pt, nil
			}
		}
	}
	return nil, errors.New("Cannot find best move")
}

func isWinningPoint(pt *point, sl *sixlet, team string) bool {
	if !pointValMatchesVal(pt, Blank) {
		return false
	}

	pzl.Set(pt, team)
	defer pzl.Set(pt, Blank)

	return sixletHasSquare(sl, team)
}

func sixletHasSquare(sl *sixlet, team string) bool {
	// check 1234
	sq1234 := true
	for _, pt := range []*point{sl.p1, sl.p2, sl.p3, sl.p4} {
		if pzl.Get(pt) != team {
			sq1234 = false
		}
	}
	if sq1234 {
		return sq1234
	}

	// check 3456
	sq3456 := true
	for _, pt := range []*point{sl.p3, sl.p4, sl.p5, sl.p6} {
		if pzl.Get(pt) != team {
			sq3456 = false
		}
	}
	return sq3456
}

type sorter func([]*point) sort.Interface

func newByXY(pts []*point) sort.Interface {
	return ByXY(pts)
}

type ByXY []*point

func (a ByXY) Len() int           { return len(a) }
func (a ByXY) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByXY) Less(i, j int) bool { return a[i].x+a[i].y < a[j].x+a[j].y }

func newByYZ(pts []*point) sort.Interface {
	return ByYZ(pts)
}

type ByYZ []*point

func (a ByYZ) Len() int           { return len(a) }
func (a ByYZ) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByYZ) Less(i, j int) bool { return a[i].y+a[i].z < a[j].y+a[j].z }

func newByXZ(pts []*point) sort.Interface {
	return ByXZ(pts)
}

type ByXZ []*point

func (a ByXZ) Len() int           { return len(a) }
func (a ByXZ) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByXZ) Less(i, j int) bool { return a[i].x+a[i].z < a[j].x+a[j].z }

type transformer func(*point, int) *point

type modder struct {
	si   sorter
	tran transformer
}
