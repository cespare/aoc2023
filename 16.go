package main

func init() {
	addSolutions(16, problem16)
}

func problem16(ctx *problemContext) {
	var c contraption
	scanner := ctx.scanner()
	for scanner.scan() {
		c.g.addRow([]byte(scanner.text()))
	}
	ctx.reportLoad()

	ctx.reportPart1(c.countEnergized(vec2{0, 0}, eastv))

	var best int
	var starts [][2]vec2
	for x := int64(0); x < c.g.cols; x++ {
		starts = append(
			starts,
			[2]vec2{{x, 0}, southv},
			[2]vec2{{x, c.g.rows - 1}, northv},
		)
	}
	for y := int64(0); y < c.g.rows; y++ {
		starts = append(
			starts,
			[2]vec2{{0, y}, eastv},
			[2]vec2{{c.g.cols - 1, y}, westv},
		)
	}
	for _, start := range starts {
		best = max(best, c.countEnergized(start[0], start[1]))
	}
	ctx.reportPart2(best)
}

type contraption struct {
	g grid[byte]
}

func (c *contraption) countEnergized(p, v vec2) int {
	seen := make(map[vec2]cardinal) // bitset of seen dirs
	c.energize(p, v, seen)
	return len(seen)
}

func (c *contraption) energize(p, v vec2, seen map[vec2]cardinal) {
	if !c.g.contains(p) {
		return
	}
	dir := vecToCard(v)
	dirs := seen[p]
	if dirs&dir != 0 {
		return
	}
	seen[p] = dirs | dir
	switch c.g.at(p) {
	case '.':
	case '\\':
		switch v {
		case northv:
			v = westv
		case eastv:
			v = southv
		case southv:
			v = eastv
		case westv:
			v = northv
		default:
			panic("bad")
		}
	case '/':
		switch v {
		case northv:
			v = eastv
		case eastv:
			v = northv
		case southv:
			v = westv
		case westv:
			v = southv
		default:
			panic("bad")
		}
	case '-':
		if v == northv || v == southv {
			c.energize(p.add(eastv), eastv, seen)
			c.energize(p.add(westv), westv, seen)
			return
		}
	case '|':
		if v == eastv || v == westv {
			c.energize(p.add(northv), northv, seen)
			c.energize(p.add(southv), southv, seen)
			return
		}
	}
	c.energize(p.add(v), v, seen)
}
