package main

import (
	"slices"
)

func init() {
	addSolutions(1, problem1)
}

func problem1(ctx *problemContext) {
	var ns []int64
	var lines []string
	scanner := ctx.scanner()
	for scanner.scan() {
		d0 := int64(-1)
		d1 := int64(0)
		for _, r := range scanner.text() {
			if r < '0' || r > '9' {
				continue
			}
			d := int64(r - '0')
			d1 = d
			if d0 < 0 {
				d0 = d
			}
		}
		ns = append(ns, 10*d0+d1)
		lines = append(lines, scanner.text())
	}
	ctx.reportLoad()

	ctx.reportPart1(SliceSum(ns))

	searcher := buildDay1Searcher(
		[]namedDigit{
			{"1", 1},
			{"2", 2},
			{"3", 3},
			{"4", 4},
			{"5", 5},
			{"6", 6},
			{"7", 7},
			{"8", 8},
			{"9", 9},
			{"one", 1},
			{"two", 2},
			{"three", 3},
			{"four", 4},
			{"five", 5},
			{"six", 6},
			{"seven", 7},
			{"eight", 8},
			{"nine", 9},
		},
	)
	var part2Sum int64
	for _, line := range lines {
		first, last := searcher.search(line)
		part2Sum += 10*int64(first) + int64(last)
	}
	ctx.reportPart2(part2Sum)
}

// day1Searcher implements an Aho-Corasick search for all elements of a set
// (the "dictionary") within an input string.
type day1Searcher struct {
	root *acNode
}

type acNode struct {
	c          byte   // transition character
	val        string // full value
	inDict     bool
	parent     *acNode
	children   []*acNode
	suffix     *acNode
	dictSuffix *acNode

	digit int
}

type namedDigit struct {
	name  string
	digit int
}

func buildDay1Searcher(dict []namedDigit) *day1Searcher {
	root := new(acNode)

	// Build main tree structure.
	for _, nd := range dict {
		n := root
		for i := 0; i < len(nd.name); i++ {
			child := n.findChild(nd.name[i])
			if child == nil {
				child = &acNode{
					c:      nd.name[i],
					val:    nd.name[:i+1],
					parent: n,
				}
				n.children = append(n.children, child)
			}
			n = child
		}
		n.inDict = true
		n.digit = nd.digit
	}

	// Fill in suffix edges with BFS.
	q := slices.Clone(root.children)
	for len(q) > 0 {
		n0 := SlicePop(&q)
		for n := n0.parent.suffix; n != nil; n = n.suffix {
			if child := n.findChild(n0.c); child != nil {
				n0.suffix = child
				break
			}
		}
		if n0.suffix == nil {
			n0.suffix = root
		}
		q = append(q, n0.children...)
	}

	// Fill in dictSuffix edges and best match with DFS.
	var fillMatch func(n *acNode)
	seen := make(map[*acNode]struct{}) // so we can cache nil results too
	fillMatch = func(n *acNode) {
		if _, ok := seen[n]; ok {
			return
		}
		seen[n] = struct{}{}
		if n.inDict {
		}
		if n.suffix != nil {
			if n.suffix.inDict {
				n.dictSuffix = n.suffix
			} else {
				fillMatch(n.suffix)
				n.dictSuffix = n.suffix.dictSuffix
			}
		}
		for _, child := range n.children {
			fillMatch(child)
		}
	}
	fillMatch(root)

	return &day1Searcher{
		root: root,
	}
}

func (s *day1Searcher) search(text string) (first, last int) {
	found := false
	emit := func(match int) {
		if !found {
			first = match
			found = true
		}
		last = match
	}
	n := s.root
	for i := 0; i < len(text); i++ {
		c := text[i]
		for {
			if next := n.findChild(c); next != nil {
				n = next
				break
			}
			n = n.suffix
			if n == nil {
				break
			}
		}
		if n == nil {
			n = s.root
			continue
		}
		for n1 := n; n1 != nil; n1 = n1.dictSuffix {
			if n1.inDict {
				emit(n1.digit)
				break
			}
		}
	}
	return first, last
}

func (n *acNode) findChild(c byte) *acNode {
	for _, child := range n.children {
		if child.c == c {
			return child
		}
	}
	return nil
}
