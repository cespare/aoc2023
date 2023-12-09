package main

import (
	"strings"
)

func init() {
	addSolutions(8, problem8)
}

func problem8(ctx *problemContext) {
	nw := &lrNetwork{
		nodes: make(map[string][2]string),
	}
	scanner := ctx.scanner()
	for scanner.scan() {
		line := scanner.text()
		if line == "" {
			continue
		}
		if nw.insts == "" {
			nw.insts = line
			continue
		}
		from, pair, ok := strings.Cut(line, " = ")
		if !ok {
			panic("bad")
		}
		pair = strings.TrimPrefix(pair, "(")
		pair = strings.TrimSuffix(pair, ")")
		to0, to1, ok := strings.Cut(pair, ", ")
		if !ok {
			panic("bad")
		}
		nw.nodes[from] = [2]string{to0, to1}

	}
	ctx.reportLoad()

	ctx.reportPart1(nw.part1())
	ctx.reportPart2(nw.part2())
}

type lrNetwork struct {
	insts string
	nodes map[string][2]string

	i   int
	cur string
}

func (n *lrNetwork) part1() int64 {
	cur := "AAA"
	i := 0
	for steps := int64(0); ; steps++ {
		if cur == "ZZZ" {
			return steps
		}
		switch n.insts[i] {
		case 'L':
			cur = n.nodes[cur][0]
		case 'R':
			cur = n.nodes[cur][1]
		default:
			panic("bad inst")
		}
		i = (i + 1) % len(n.insts)
	}
}

func (n *lrNetwork) part2() int64 {
	// The actual problem input we are given is much simpler than the
	// general problem as described: each starting node only connects to a
	// single exit node, and the cycle length is the entire chain. So we can
	// simply compute the LCM of all the cycles.
	res := int64(1)
	for node := range n.nodes {
		if strings.HasSuffix(node, "A") {
			res = lcm(res, n.cycleLength(node))
		}
	}
	return res
}

func (n *lrNetwork) cycleLength(start string) int64 {
	cur := start
	i := 0
	for steps := int64(0); ; steps++ {
		if strings.HasSuffix(cur, "Z") {
			return steps
		}
		switch n.insts[i] {
		case 'L':
			cur = n.nodes[cur][0]
		case 'R':
			cur = n.nodes[cur][1]
		default:
			panic("bad inst")
		}
		i = (i + 1) % len(n.insts)
	}
}

func lcm(a, b int64) int64 {
	return a * b / gcd(a, b)
}

func gcd(a, b int64) int64 {
	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}
	return a
}
