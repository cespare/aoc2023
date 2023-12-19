package main

import (
	"fmt"
	"slices"
	"strings"
)

func init() {
	addSolutions(19, problem19)
}

func problem19(ctx *problemContext) {
	a := &avalanche{flow: make(map[string]*machineWorkflow)}
	firstSection := true
	scanner := ctx.scanner()
	for scanner.scan() {
		line := scanner.text()
		if line == "" {
			firstSection = false
			continue
		}
		if firstSection {
			w := parseMachineWorkflow(line)
			a.flow[w.name] = w
		} else {
			a.parts = append(a.parts, parseMachinePart(line))
		}
	}
	a.init()
	ctx.reportLoad()

	ctx.reportPart1(a.part1())
	ctx.reportPart2(a.part2())
}

type avalanche struct {
	flow  map[string]*machineWorkflow
	parts []machinePart

	backLinks map[string][]workflowCoord
}

func (a *avalanche) init() {
	a.backLinks = make(map[string][]workflowCoord)
	for _, w := range a.flow {
		for i, r := range w.rules {
			coord := workflowCoord{w.name, i}
			a.backLinks[r.dest] = append(a.backLinks[r.dest], coord)
		}
	}
}

type workflowCoord struct {
	name string
	i    int // rule index
}

func (c workflowCoord) String() string {
	return fmt.Sprintf("%s[%d]", c.name, c.i)
}

type machineWorkflow struct {
	name  string
	rules []machineRule
}

func parseMachineWorkflow(s string) *machineWorkflow {
	name, rest, ok := strings.Cut(s, "{")
	if !ok {
		panic("bad")
	}
	w := &machineWorkflow{name: name}
	rest = strings.TrimSuffix(rest, "}")
	for _, s := range strings.Split(rest, ",") {
		var r machineRule
		if cond, dest, ok := strings.Cut(s, ":"); ok {
			v, targ, ok := strings.Cut(cond, ">")
			if !ok {
				v, targ, ok = strings.Cut(cond, "<")
				if !ok {
					panic("bad")
				}
				r.lt = true
			}
			r.v = v
			r.targ = parseInt(targ)
			r.dest = dest
		} else {
			r.dest = s
		}
		w.rules = append(w.rules, r)
	}
	return w
}

type machineRule struct {
	v    string
	lt   bool
	targ int64
	dest string
}

func (r machineRule) matches(p machinePart) bool {
	var n int64
	switch r.v {
	case "x":
		n = p.x
	case "m":
		n = p.m
	case "a":
		n = p.a
	case "s":
		n = p.s
	case "":
		return true
	default:
		panic("bad")
	}
	if r.lt {
		return n < r.targ
	}
	return n > r.targ
}

type machinePart struct {
	x, m, a, s int64
}

func parseMachinePart(s string) machinePart {
	s = strings.TrimPrefix(strings.TrimSuffix(s, "}"), "{")
	var p machinePart
	for _, part := range strings.Split(s, ",") {
		v, ns, ok := strings.Cut(part, "=")
		if !ok {
			panic("bad")
		}
		n := parseInt(ns)
		switch v {
		case "x":
			p.x = n
		case "m":
			p.m = n
		case "a":
			p.a = n
		case "s":
			p.s = n
		default:
			panic("bad")
		}
	}
	return p
}

func (a *avalanche) exec(p machinePart) bool {
	cur := "in"
flowLoop:
	for {
		switch cur {
		case "A":
			return true
		case "R":
			return false
		}
		w := a.flow[cur]
		for _, r := range w.rules {
			if r.matches(p) {
				cur = r.dest
				continue flowLoop
			}
		}
		panic("unreached")
	}
}

func (a *avalanche) part1() int64 {
	var total int64
	for _, p := range a.parts {
		if a.exec(p) {
			total += p.x + p.m + p.a + p.s
		}
	}
	return total
}

type rangeSet struct {
	ranges []rnge // non-overlapping, ordered
}

type rnge struct {
	start, end int64 // inclusive
}

func newRangeSet(start, end int64) *rangeSet {
	return &rangeSet{
		ranges: []rnge{{start, end}},
	}
}

func (rs *rangeSet) clone() *rangeSet {
	return &rangeSet{
		ranges: slices.Clone(rs.ranges),
	}
}

func (rs *rangeSet) count() int64 {
	var total int64
	for _, r := range rs.ranges {
		total += r.end - r.start + 1
	}
	return total
}

func (rs *rangeSet) clearGT(n int64) {
	for i, r := range rs.ranges {
		if n >= r.end {
			continue
		}
		if n < r.start {
			rs.ranges = rs.ranges[:i]
			return
		}
		rs.ranges[i].end = n
		rs.ranges = rs.ranges[:i+1]
		return
	}
}

func (rs *rangeSet) clearGEQ(n int64) {
	for i, r := range rs.ranges {
		if n > r.end {
			continue
		}
		if n <= r.start {
			rs.ranges = rs.ranges[:i]
			return
		}
		rs.ranges[i].end = n - 1
		rs.ranges = rs.ranges[:i+1]
		return
	}
}

func (rs *rangeSet) clearLT(n int64) {
	for i := len(rs.ranges) - 1; i >= 0; i-- {
		r := rs.ranges[i]
		if n <= r.start {
			continue
		}
		if n > r.end {
			rs.ranges = rs.ranges[i+1:]
			return
		}
		rs.ranges[i].start = n
		rs.ranges = rs.ranges[i:]
		return
	}
}

func (rs *rangeSet) clearLEQ(n int64) {
	for i := len(rs.ranges) - 1; i >= 0; i-- {
		r := rs.ranges[i]
		if n < r.start {
			continue
		}
		if n >= r.end {
			rs.ranges = rs.ranges[i+1:]
			return
		}
		rs.ranges[i].start = n + 1
		rs.ranges = rs.ranges[i:]
		return
	}
}

type partRanges struct {
	x, m, a, s *rangeSet
}

func newPartRanges() partRanges {
	return partRanges{
		x: newRangeSet(1, 4000),
		m: newRangeSet(1, 4000),
		a: newRangeSet(1, 4000),
		s: newRangeSet(1, 4000),
	}
}

func (pb partRanges) clone() partRanges {
	return partRanges{
		x: pb.x.clone(),
		m: pb.m.clone(),
		a: pb.a.clone(),
		s: pb.s.clone(),
	}
}

func (pb partRanges) valid() bool {
	return pb.x.count() > 0 && pb.m.count() > 0 && pb.a.count() > 0 && pb.s.count() > 0
}

func (pb partRanges) combos() int64 {
	return pb.x.count() * pb.m.count() * pb.a.count() * pb.s.count()
}

func (a *avalanche) multiExec(coord workflowCoord, pb partRanges) int64 {
	if !pb.valid() {
		return 0
	}
	r := a.flow[coord.name].rules[coord.i]
	var combos int64
	if r.v != "" {
		pbNeg := pb.clone()
		var rs, rsNeg *rangeSet
		switch r.v {
		case "x":
			rs = pb.x
			rsNeg = pbNeg.x
		case "m":
			rs = pb.m
			rsNeg = pbNeg.m
		case "a":
			rs = pb.a
			rsNeg = pbNeg.a
		case "s":
			rs = pb.s
			rsNeg = pbNeg.s
		}
		if r.lt {
			rs.clearGEQ(r.targ)
			rsNeg.clearLT(r.targ)
		} else {
			rs.clearLEQ(r.targ)
			rsNeg.clearGT(r.targ)
		}
		coordNeg := coord
		coordNeg.i++
		combos += a.multiExec(coordNeg, pbNeg)
	}
	switch r.dest {
	case "R":
	case "A":
		combos += pb.combos()
	default:
		coord := workflowCoord{r.dest, 0}
		combos += a.multiExec(coord, pb)
	}
	return combos
}

func (a *avalanche) part2() int64 {
	coord := workflowCoord{"in", 0}
	return a.multiExec(coord, newPartRanges())
}
