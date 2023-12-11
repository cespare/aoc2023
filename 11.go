package main

func init() {
	addSolutions(11, problem11)
}

func problem11(ctx *problemContext) {
	var g galaxy
	scanner := ctx.scanner()
	for scanner.scan() {
		g.g.addRow([]byte(scanner.text()))

	}
	ctx.reportLoad()

	ctx.reportPart1(g.expandAndSum(1))
	ctx.reportPart2(g.expandAndSum(1e6 - 1))
}

type galaxy struct {
	g grid[byte]
}

func (g *galaxy) expandAndSum(inc int64) int64 {
	var vs []vec2
	g.g.forEach(func(v vec2, c byte) bool {
		if c == '#' {
			vs = append(vs, v)
		}
		return true
	})
	cols := g.emptyCols()
	for i := len(cols) - 1; i >= 0; i-- {
		x := cols[i]
		for j, v := range vs {
			if v.x > x {
				v.x += inc
			}
			vs[j] = v
		}
	}
	rows := g.emptyRows()
	for i := len(rows) - 1; i >= 0; i-- {
		y := rows[i]
		for j, v := range vs {
			if v.y > y {
				v.y += inc
			}
			vs[j] = v
		}
	}
	var sum int64
	for i, v0 := range vs {
		for _, v1 := range vs[i+1:] {
			sum += v1.sub(v0).mag()
		}
	}
	return sum
}

func (g *galaxy) emptyCols() []int64 {
	empty := make(map[int64]struct{})
	for x := int64(0); x < g.g.cols; x++ {
		empty[x] = struct{}{}
	}
	g.g.forEach(func(v vec2, c byte) bool {
		if c != '.' {
			delete(empty, v.x)
		}
		return true
	})
	var cols []int64
	for x := int64(0); x < g.g.cols; x++ {
		if _, ok := empty[x]; ok {
			cols = append(cols, x)
		}
	}
	return cols
}

func (g *galaxy) emptyRows() []int64 {
	empty := make(map[int64]struct{})
	for y := int64(0); y < g.g.rows; y++ {
		empty[y] = struct{}{}
	}
	g.g.forEach(func(v vec2, c byte) bool {
		if c != '.' {
			delete(empty, v.y)
		}
		return true
	})
	var rows []int64
	for y := int64(0); y < g.g.rows; y++ {
		if _, ok := empty[y]; ok {
			rows = append(rows, y)
		}
	}
	return rows
}
