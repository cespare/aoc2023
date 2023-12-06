package main

import (
	"slices"
	"strings"
)

func init() {
	addSolutions(5, problem5)
}

func problem5(ctx *problemContext) {
	var a almanac
	var curMap []seedMapRange
	scanner := ctx.scanner()
	for scanner.scan() {
		line := scanner.text()
		switch {
		case len(a.seeds) == 0:
			s := strings.TrimPrefix(line, "seeds: ")
			for _, sn := range strings.Fields(s) {
				a.seeds = append(a.seeds, parseInt(sn))
			}
		case line == "":
		case strings.Contains(line, "map:"):
			if len(curMap) > 0 {
				a.maps = append(a.maps, curMap)
			}
			curMap = nil
		default:
			var ns []int64
			for _, sn := range strings.Fields(line) {
				ns = append(ns, parseInt(sn))
			}
			if len(ns) != 3 {
				panic("bad range")
			}
			rnge := seedMapRange{dest: ns[0], source: ns[1], len: ns[2]}
			curMap = append(curMap, rnge)
		}
	}
	if len(curMap) > 0 {
		a.maps = append(a.maps, curMap)
	}
	a.sortMaps()
	ctx.reportLoad()

	ctx.reportPart1(a.solve1())

	ctx.reportPart2(a.solve2())
}

type almanac struct {
	seeds []int64
	maps  [][]seedMapRange
}

type seedMapRange struct {
	dest   int64
	source int64
	len    int64
}

func (a *almanac) sortMaps() {
	for _, m := range a.maps {
		slices.SortFunc(m, func(r0, r1 seedMapRange) int {
			return int(r0.source - r1.source)
		})
	}
}

func (a *almanac) solve1() int64 {
	best := int64(-1)
	for _, seed := range a.seeds {
		if loc := a.getLocation(seed); best < 0 || loc < best {
			best = loc
		}
	}
	return best
}

func (a *almanac) getLocation(seed int64) int64 {
	n := seed
	for _, m := range a.maps {
		n = convertSeed(n, m)
	}
	return n
}

func convertSeed(source int64, m []seedMapRange) int64 {
	for _, r := range m {
		if source < r.source {
			return source
		}
		if source < r.source+r.len {
			return r.dest + (source - r.source)
		}
	}
	return source
}

func (a *almanac) solve2() int64 {
	best := int64(-1)
	for i := 0; i < len(a.seeds); i += 2 {
		start := a.seeds[i]
		end := a.seeds[i] + a.seeds[i+1]
		best1 := bestLocationForRange(a.maps, start, end)
		if best < 0 || best1 < best {
			best = best1
		}
	}
	return best
}

func bestLocationForRange(maps [][]seedMapRange, start, end int64) int64 {
	if len(maps) == 0 {
		return start
	}
	best := int64(-1)
	updateBest := func(candidate int64) {
		if best < 0 || candidate < best {
			best = candidate
		}
	}
	m := maps[0]
	maps = maps[1:]
	for start < end {
		if len(m) == 0 {
			updateBest(bestLocationForRange(maps, start, end))
			break
		}
		if start < m[0].source {
			next := min(end, m[0].source)
			updateBest(bestLocationForRange(maps, start, next))
			start = next
		}
		sourceEnd := m[0].source + m[0].len
		if start < sourceEnd {
			start1 := m[0].dest + (start - m[0].source)
			next := min(end, sourceEnd)
			end1 := m[0].dest + (next - m[0].source)
			updateBest(bestLocationForRange(maps, start1, end1))
			start = next
		}
		m = m[1:]
	}
	return best
}
