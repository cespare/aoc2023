package main

import (
	"fmt"
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

type machineInterval struct {
	start, end int64 // inclusive
}

func (iv *machineInterval) count() int64 {
	return iv.end - iv.start + 1
}

func (iv *machineInterval) clearGT(n int64) {
	switch {
	case n >= iv.end:
	case n < iv.start:
		iv.end = iv.start
	default:
		iv.end = n
	}
}

func (iv *machineInterval) clearGEQ(n int64) {
	switch {
	case n > iv.end:
	case n <= iv.start:
		iv.end = iv.start
	default:
		iv.end = n - 1
	}
}

func (iv *machineInterval) clearLT(n int64) {
	switch {
	case n <= iv.start:
	case n > iv.end:
		iv.start = iv.end
	default:
		iv.start = n
	}
}

func (iv *machineInterval) clearLEQ(n int64) {
	switch {
	case n < iv.start:
	case n >= iv.end:
		iv.start = iv.end
	default:
		iv.start = n + 1
	}
}

type partIntervals struct {
	x, m, a, s machineInterval
}

func newPartRanges() *partIntervals {
	return &partIntervals{
		x: machineInterval{1, 4000},
		m: machineInterval{1, 4000},
		a: machineInterval{1, 4000},
		s: machineInterval{1, 4000},
	}
}

func (pi *partIntervals) clone() *partIntervals {
	pi1 := *pi
	return &pi1
}

func (pi *partIntervals) valid() bool {
	return pi.x.count() > 0 && pi.m.count() > 0 && pi.a.count() > 0 && pi.s.count() > 0
}

func (pi *partIntervals) combos() int64 {
	return pi.x.count() * pi.m.count() * pi.a.count() * pi.s.count()
}

func (a *avalanche) multiExec(coord workflowCoord, pi *partIntervals) int64 {
	if !pi.valid() {
		return 0
	}
	r := a.flow[coord.name].rules[coord.i]
	var combos int64
	if r.v != "" {
		piNeg := pi.clone()
		var iv, ivNeg *machineInterval
		switch r.v {
		case "x":
			iv = &pi.x
			ivNeg = &piNeg.x
		case "m":
			iv = &pi.m
			ivNeg = &piNeg.m
		case "a":
			iv = &pi.a
			ivNeg = &piNeg.a
		case "s":
			iv = &pi.s
			ivNeg = &piNeg.s
		}
		if r.lt {
			iv.clearGEQ(r.targ)
			ivNeg.clearLT(r.targ)
		} else {
			iv.clearLEQ(r.targ)
			ivNeg.clearGT(r.targ)
		}
		coordNeg := coord
		coordNeg.i++
		combos += a.multiExec(coordNeg, piNeg)
	}
	switch r.dest {
	case "R":
	case "A":
		combos += pi.combos()
	default:
		coord := workflowCoord{r.dest, 0}
		combos += a.multiExec(coord, pi)
	}
	return combos
}

func (a *avalanche) part2() int64 {
	coord := workflowCoord{"in", 0}
	return a.multiExec(coord, newPartRanges())
}
