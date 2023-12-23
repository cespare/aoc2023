package main

import (
	"fmt"
)

func init() {
	addSolutions(21, problem21)
}

func problem21(ctx *problemContext) {
	var g garden
	scanner := ctx.scanner()
	for scanner.scan() {
		g.g.addRow([]byte(scanner.text()))
	}
	g.init()
	ctx.reportLoad()

	// ig := &infinityGarden{
	// 	outer: map[vec2]*outerGardenState{
	// 		{0, 0}: {
	// 			g:        g.g.clone(),
	// 			frontier: map[vec2]struct{}{g.start: {}},
	// 		},
	// 	},
	// }

	// frontier := map[vec2]struct{}{g.start: {}}
	// for i := 0; i < 1000; i++ {
	// 	frontier = g.advance(frontier)
	// }
	// ctx.reportPart1(len(frontier))

	frontier := map[vec2]struct{}{g.start: {}}
	for i := 0; i < 501; i++ {
		switch i {
		case 65, 65 + 131, 65 + 131*2, 65 + 131*3:
			fmt.Println(i, len(frontier))
			// case 1, 2, 3, 6, 10, 50, 100, 500:
			// 	fmt.Println(i, len(frontier))
		}
		frontier = g.advanceInfinite(frontier)
	}

	// TODO: describe how the input file (and target step number) are special.
	// Describe how to use quadratic extrapolation based on the first few
	// sequence terms.
}

type garden struct {
	g     grid[byte]
	start vec2
}

func (g *garden) print(frontier map[vec2]struct{}) {
	g1 := g.g.clone()
	for v := range frontier {
		g1.set(v, 'O')
	}
	for _, row := range g1.g {
		fmt.Println(string(row))
	}
}

func (g *garden) init() {
	g.g.forEach(func(v vec2, c byte) bool {
		if c == 'S' {
			g.start = v
			g.g.set(v, '.')
			return false
		}
		return true
	})
}

func (g *garden) advance(frontier map[vec2]struct{}) map[vec2]struct{} {
	next := make(map[vec2]struct{}, len(frontier))
	for v := range frontier {
		for _, n := range v.neighbors4() {
			if !g.g.contains(n) {
				continue
			}
			if g.g.at(n) != '.' {
				continue
			}
			next[n] = struct{}{}
		}
	}
	return next
}

func (g *garden) advanceInfinite(frontier map[vec2]struct{}) map[vec2]struct{} {
	next := make(map[vec2]struct{}, len(frontier))
	for v := range frontier {
		for _, n := range v.neighbors4() {
			gv := vec2{mod(n.x, g.g.cols), mod(n.y, g.g.rows)}
			if g.g.at(gv) == '.' {
				next[n] = struct{}{}
			}
		}
	}
	return next
}

func mod(n, m int64) int64 {
	if m < 0 {
		panic("no")
	}
	r := n % m
	if r < 0 {
		r += m
	}
	return r
}
