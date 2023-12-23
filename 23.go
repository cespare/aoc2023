package main

func init() {
	addSolutions(23, problem23)
}

func problem23(ctx *problemContext) {
	var h hike
	scanner := ctx.scanner()
	for scanner.scan() {
		h.g.addRow([]byte(scanner.text()))
	}
	h.init()
	ctx.reportLoad()

	ctx.reportPart1(h.part1())
	ctx.reportPart2(h.part2())
}

type hike struct {
	g     grid[byte]
	start vec2
	end   vec2
}

func (h *hike) init() {
	for x := int64(0); x < h.g.cols; x++ {
		v := vec2{x, 0}
		if h.g.at(v) == '.' {
			h.start = v
		}
		v.y = h.g.rows - 1
		if h.g.at(v) == '.' {
			h.end = v
		}
	}
}

func (h *hike) part1() int {
	return h.longest(h.start, make(map[vec2]struct{}))
}

func (h *hike) longest(v vec2, seen map[vec2]struct{}) int {
	if v == h.end {
		return 0
	}
	var neighbors []vec2
	addNeighbors := func(ns ...vec2) {
		for _, n := range ns {
			if !h.g.contains(n) {
				continue
			}
			if h.g.at(n) == '#' {
				continue
			}
			if _, ok := seen[n]; ok {
				continue
			}
			neighbors = append(neighbors, n)
		}
	}
	switch h.g.at(v) {
	case '^':
		addNeighbors(v.add(northv))
	case '>':
		addNeighbors(v.add(eastv))
	case 'v':
		addNeighbors(v.add(southv))
	case '<':
		addNeighbors(v.add(westv))
	case '.':
		addNeighbors(v.neighbors4()...)
	default:
		panic("bad")
	}
	switch len(neighbors) {
	case 0:
		return -1
	case 1:
		best := h.longest(neighbors[0], seen)
		if best >= 0 {
			best++
		}
		return best
	}
	best := -1
	seen[v] = struct{}{}
	for _, n := range neighbors {
		best = max(best, h.longest(n, seen))
	}
	delete(seen, v)
	if best >= 0 {
		best++
	}
	return best
}

type hikeSegment struct {
	n      int
	end    vec2 // typically a junction, but could also be endpoint (or dead end)
	endIdx int  // index in junction slice
}

type hikeSegmentGraph struct {
	junctions [][]*hikeSegment
	start     *hikeSegment
}

func (h *hike) buildSegmentGraph() *hikeSegmentGraph {
	junctions := make(map[vec2][]*hikeSegment)
	var g hikeSegmentGraph
	h.g.forEach(func(v vec2, c byte) bool {
		if c != '.' {
			return true
		}
		var numNeighbors int
		for _, n := range v.neighbors4() {
			if h.g.contains(n) && h.g.at(n) != '#' {
				numNeighbors++
			}
		}
		if numNeighbors > 2 {
			junctions[v] = nil
		}
		return true
	})
	for v := range junctions {
		h.fillJunction(v, junctions)
	}
	g.start = h.findSegment(h.start, junctions)
	g.start.n-- // starting point doesn't count

	// Assign bitset indices to each junction.
	junctionIdx := make(map[vec2]int)
	for v, segs := range junctions {
		junctionIdx[v] = len(g.junctions)
		g.junctions = append(g.junctions, segs)
		if len(g.junctions) > 64 {
			panic("too many junctions")
		}
	}
	// Tack on indices for start and end.
	junctionIdx[h.start] = len(g.junctions)
	g.junctions = append(g.junctions, nil)
	junctionIdx[h.end] = len(g.junctions)
	g.junctions = append(g.junctions, nil)
	for _, segs := range g.junctions {
		for _, seg := range segs {
			seg.endIdx = junctionIdx[seg.end]
		}
	}
	g.start.endIdx = junctionIdx[g.start.end]

	return &g
}

func (h *hike) fillJunction(v vec2, junctions map[vec2][]*hikeSegment) {
	for _, n := range v.neighbors4() {
		if !h.g.contains(n) {
			continue
		}
		if h.g.at(n) == '#' {
			continue
		}
		junctions[v] = append(junctions[v], h.findSegment(n, junctions))
	}
}

func (h *hike) findSegment(start vec2, junctions map[vec2][]*hikeSegment) *hikeSegment {
	seg := &hikeSegment{end: start}
	prev := vec2{-1, -1}
searchLoop:
	for {
		seg.n++
		for _, n := range seg.end.neighbors4() {
			if n == prev {
				continue
			}
			if !h.g.contains(n) {
				continue
			}
			if h.g.at(n) == '#' {
				continue
			}
			if _, ok := junctions[n]; ok {
				if seg.end == start {
					continue
				}
				seg.end = n
				break
			}
			prev = seg.end
			seg.end = n
			continue searchLoop
		}
		return seg
	}
}

func (h *hike) part2() int {
	g := h.buildSegmentGraph()

	var best int
	st := hikeSearchState{
		seg:  g.start,
		seen: 0,
	}
	stk := []hikeSearchState{st}
	for len(stk) > 0 {
		st := StackPop(&stk)
		n := st.n + st.seg.n
		if st.seg.end == h.end {
			best = max(best, n)
			continue
		}
		st.n = n + 1 // count junction
		st.seen |= uint64(1) << st.seg.endIdx
		for _, next := range g.junctions[st.seg.endIdx] {
			if st.seen&(uint64(1)<<next.endIdx) > 0 {
				continue
			}
			st.seg = next
			StackPush(&stk, st)
		}
	}
	return best
}

type hikeSearchState struct {
	n    int
	seg  *hikeSegment
	seen uint64
}
