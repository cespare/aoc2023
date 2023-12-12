package main

import (
	"strings"
)

func init() {
	addSolutions(12, problem12)
}

func problem12(ctx *problemContext) {
	var recs []*springRecord
	scanner := ctx.scanner()
	for scanner.scan() {
		recs = append(recs, parseSpringRecord(scanner.text()))

	}
	ctx.reportLoad()

	var part1Sum int64
	for _, rec := range recs {
		part1Sum += rec.combos()
	}
	ctx.reportPart1(part1Sum)

	var part2Sum int64
	for _, rec := range recs {
		rec = rec.unfold()
		part2Sum += rec.combos()
	}
	ctx.reportPart2(part2Sum)
}

type springRecord struct {
	s      string
	groups []int
}

func parseSpringRecord(s string) *springRecord {
	r, rest, ok := strings.Cut(s, " ")
	if !ok {
		panic("bad")
	}
	var groups []int
	for _, f := range strings.Split(rest, ",") {
		groups = append(groups, int(parseInt(f)))
	}
	return &springRecord{s: r, groups: groups}
}

func (r *springRecord) combos() int64 {
	return r.combosRec(0, 0, 0, make(map[springState]int64))
}

type springState struct {
	i         int
	groupIdx  int
	groupSize int
}

func (r *springRecord) combosRec(i, groupIdx, groupSize int, memo map[springState]int64) (res int64) {
	state := springState{i: i, groupIdx: groupIdx, groupSize: groupSize}
	if n, ok := memo[state]; ok {
		return n
	}
	defer func() {
		memo[state] = res
	}()

	// groupCmp returns: 1 if we've exceeded len(r.groups) or the size of
	// the current group; otherwise -1/0 if we are less than/equal to the
	// size of the current group.
	groupCmp := func() int {
		if groupIdx >= len(r.groups) {
			return 1
		}
		if groupSize < r.groups[groupIdx] {
			return -1
		}
		if groupSize == r.groups[groupIdx] {
			return 0
		}
		return 1
	}

	if i == len(r.s) {
		if groupSize > 0 {
			if groupCmp() != 0 {
				return 0
			}
			groupIdx++
		}
		if groupIdx != len(r.groups) {
			return 0
		}
		return 1
	}

	switch r.s[i] {
	case '.':
		if groupSize == 0 {
			return r.combosRec(i+1, groupIdx, 0, memo)
		}
		if groupCmp() != 0 {
			return 0
		}
		return r.combosRec(i+1, groupIdx+1, 0, memo)
	case '#':
		groupSize++
		if groupCmp() > 0 {
			return 0
		}
		return r.combosRec(i+1, groupIdx, groupSize, memo)
	case '?':
		var res int64

		// .
		if groupSize == 0 {
			res += r.combosRec(i+1, groupIdx, 0, memo)
		} else if groupCmp() == 0 {
			res += r.combosRec(i+1, groupIdx+1, 0, memo)
		}

		// #
		groupSize++
		if groupCmp() <= 0 {
			res += r.combosRec(i+1, groupIdx, groupSize, memo)
		}
		return res
	default:
		panic("bad")
	}
}

func (r *springRecord) unfold() *springRecord {
	var r1 springRecord
	for i := 0; i < 5; i++ {
		r1.s += r.s
		if i < 4 {
			r1.s += "?"
		}
		r1.groups = append(r1.groups, r.groups...)
	}
	return &r1
}
