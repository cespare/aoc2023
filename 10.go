package main

import (
	"fmt"
	"strings"
)

func init() {
	addSolutions(10, problem10)
}

func problem10(ctx *problemContext) {
	var g pipeGrid
	scanner := ctx.scanner()
	for scanner.scan() {
		g.g.addRow([]byte(scanner.text()))

	}
	ctx.reportLoad()

	g.init()

	ctx.reportPart1(g.part1())

	g.fill()
	g.print()
	ctx.reportPart2(g.count(g.innerRegion))
}

type pipeGrid struct {
	g           grid[byte]
	loop        []vec2
	loopSet     map[vec2]struct{}
	start       vec2
	innerRegion byte // X or Y
}

func (g *pipeGrid) print() {
	r := strings.NewReplacer(
		"-", "─",
		"|", "│",
		"F", "┌",
		"J", "┘",
		"L", "└",
		"7", "┐",
		"7", "┐",
	)
	for _, row := range g.g.g {
		fmt.Println(r.Replace(string(row)))
	}
}

var (
	northv = vec2{0, -1}
	eastv  = vec2{1, 0}
	southv = vec2{0, 1}
	westv  = vec2{-1, 0}
)

type cardinal uint8

const (
	north cardinal = 1 << iota
	east
	south
	west
)

func (c cardinal) GoString() string {
	var ss []string
	if c&north != 0 {
		ss = append(ss, "north")
	}
	if c&east != 0 {
		ss = append(ss, "east")
	}
	if c&south != 0 {
		ss = append(ss, "south")
	}
	if c&west != 0 {
		ss = append(ss, "west")
	}
	return strings.Join(ss, "|")
}

func cardToVec(c cardinal) vec2 {
	switch c {
	case north:
		return northv
	case east:
		return eastv
	case south:
		return southv
	case west:
		return westv
	default:
		panic("unreached")
	}
}

func vecToCard(v vec2) cardinal {
	switch v {
	case northv:
		return north
	case eastv:
		return east
	case southv:
		return south
	case westv:
		return west
	default:
		panic("unreached")
	}
}

func pipeExits(c byte) cardinal {
	switch c {
	case '-':
		return east | west
	case '7':
		return west | south
	case '|':
		return north | south
	case 'J':
		return north | west
	case 'L':
		return east | north
	case 'F':
		return south | east
	case 'S':
		return north | east | south | west
	default:
		fmt.Println(string(c))
		panic("unreached")
	}
}

func (g *pipeGrid) init() {
	g.g.forEach(func(v vec2, b byte) bool {
		if b == 'S' {
			g.start = v
			return false
		}
		return true
	})

	// Map out the loop.
	cur := g.start
	g.loop = []vec2{cur}
	g.loopSet = map[vec2]struct{}{cur: {}}
	for {
		var connected []vec2
		p0 := g.g.at(cur)
		for _, n := range cur.neighbors4() {
			if !g.g.contains(n) {
				continue
			}
			if _, ok := g.loopSet[n]; ok {
				continue
			}
			p1 := g.g.at(n)
			if p1 == '.' {
				continue
			}
			if pipesConnected(p0, p1, n.sub(cur)) {
				connected = append(connected, n)
			}
		}
		if len(connected) == 0 {
			break
		}
		if cur == g.start && len(connected) != 2 {
			fmt.Printf("\033[01;34m>>>> len(connected): %v\x1B[m\n", len(connected))
			panic("bad")
		}
		if cur != g.start && len(connected) != 1 {
			panic("bad")
		}
		cur = connected[0]
		g.loop = append(g.loop, cur)
		g.loopSet[cur] = struct{}{}
	}

}

func pipesConnected(p0, p1 byte, dv vec2) bool {
	if pipeExits(p0)&vecToCard(dv) == 0 {
		return false
	}
	if pipeExits(p1)&vecToCard(dv.mul(-1)) == 0 {
		return false
	}
	return true
}

func (g *pipeGrid) part1() int {
	return len(g.loop) / 2
}

func (g *pipeGrid) fill() {
	for i, cur := range g.loop {
		var prev vec2
		if i == 0 {
			prev = g.loop[len(g.loop)-1]
		} else {
			prev = g.loop[i-1]
		}
		dv := cur.sub(prev)
		c := g.g.at(cur)
		for _, n := range cur.neighbors4() {
			if _, ok := g.loopSet[n]; ok {
				continue
			}
			g.flood(n, pipePolarity(n.sub(cur), dv, c))
		}
	}
}

func (g *pipeGrid) flood(v vec2, c byte) {
	q := []vec2{v}
	visited := make(map[vec2]struct{})
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		if _, ok := visited[v]; ok {
			continue
		}
		if !g.g.contains(v) {
			if c == 'X' {
				g.innerRegion = 'Y'
			} else {
				g.innerRegion = 'X'
			}
			continue
		}
		c0 := g.g.at(v)
		if c0 == 'X' || c0 == 'Y' {
			if c0 != c {
				panic("inconsistent")
			}
			continue
		}
		g.g.set(v, c)
		for _, n := range v.neighbors4() {
			if _, ok := g.loopSet[n]; !ok {
				q = append(q, n)
			}
		}
	}
}

func pipePolarity(dn, dv vec2, c byte) byte {
	// Rotate from the entrance side clockwise to the neighbor in question
	// and see if we pass over the exit side on the way.
	dir := vecToCard(dv.mul(-1))
	exits := pipeExits(c)
	ndir := vecToCard(dn)
	result := byte('X')
	for {
		// left-rotate a 4-bit pattern by 1.
		dir = ((dir << 1) | (dir >> 3)) & 0xf
		if dir == ndir {
			return result
		}
		if exits&dir != 0 {
			result = 'Y'
		}
	}
}

func (g *pipeGrid) count(c byte) int64 {
	var n int64
	g.g.forEach(func(_ vec2, c1 byte) bool {
		if c1 == c {
			n++
		}
		return true
	})
	return n
}
