package main

import (
	"maps"
	"math/rand"
	"slices"
	"strings"
)

func init() {
	addSolutions(25, problem25)
}

func problem25(ctx *problemContext) {
	var g graph
	scanner := ctx.scanner()
	for scanner.scan() {
		g.addNode(scanner.text())
	}
	ctx.reportLoad()
	// pretty.Println(g)

	ctx.reportPart1(g.cut3Score())
}

type graph struct {
	nodes map[string]map[string]struct{}
}

func (g *graph) clone() *graph {
	g1 := &graph{
		nodes: make(map[string]map[string]struct{}),
	}
	for src, dsts := range g.nodes {
		g1.nodes[src] = maps.Clone(dsts)
	}
	return g1
}

func (g *graph) addNode(s string) {
	src, rest, ok := strings.Cut(s, ":")
	if !ok {
		panic("bad")
	}
	dsts := make(map[string]struct{})
	for _, dst := range strings.Fields(rest) {
		dsts[dst] = struct{}{}
	}
	if g.nodes == nil {
		g.nodes = make(map[string]map[string]struct{})
	}
	g.nodes[src] = dsts
}

func (g *graph) cut3Score() int64 {
	// https://en.wikipedia.org/wiki/Karger%27s_algorithm
	for {
		ag := g.toAdj()
		// pretty.Println(ag)
		for ag.numNodes() > 2 {
			// fmt.Printf("\033[01;34m>>>> ag.numNodes(): %v\x1B[m\n", ag.numNodes())
			// fmt.Printf("\033[01;34m>>>> ag.numEdges: %v\x1B[m\n", ag.numEdges)
			ag.contract()
		}
		if ag.numEdges < 3 {
			panic("got better min cut")
		}
		if ag.numEdges == 3 {
			partitions := ag.partitionSizes()
			return partitions[0] * partitions[1]
		}
	}
}

func (g *graph) toAdj() *adjGraph {
	idxs := make(map[string]int)
	var ag adjGraph
	addNode := func(n string) {
		if _, ok := idxs[n]; ok {
			return
		}
		idx := len(idxs)
		idxs[n] = idx
		ag.nodes = append(ag.nodes, nodeIdx{n, idx})
		ag.multi = append(ag.multi, 1)
	}
	for src, dsts := range g.nodes {
		addNode(src)
		for dst := range dsts {
			addNode(dst)
		}
	}
	ag.edges = make([][]int, len(ag.nodes))
	for i := range ag.edges {
		ag.edges[i] = make([]int, len(ag.nodes))
	}
	for src, dsts := range g.nodes {
		srcIdx := idxs[src]
		for dst := range dsts {
			dstIdx := idxs[dst]
			ag.addEdge(srcIdx, dstIdx, 1)
		}
	}
	return &ag
}

type adjGraph struct {
	nodes    []nodeIdx
	multi    []int
	edges    [][]int
	numEdges int
}

type nodeIdx struct {
	name string
	idx  int
}

func (ag *adjGraph) contract() {
	// Select a random edge.
	edgeIdx := rand.Intn(ag.numEdges)
	u, v := -1, -1
edgeLoop:
	for a, dsts := range ag.edges {
		for b, n := range dsts {
			edgeIdx -= n
			if edgeIdx < 0 {
				u, v = a, b
				break edgeLoop
			}
		}
	}
	if u == -1 {
		panic("no random edge selected")
	}

	// Merge v into u.
	for _, ni := range ag.nodes {
		if ni.idx == u || ni.idx == v {
			continue
		}
		n := ag.deleteEdges(ni.idx, v)
		if n == 0 {
			continue
		}
		ag.addEdge(ni.idx, u, n)
	}
	ag.deleteEdges(u, v)
	ag.multi[u] += ag.multi[v]
	ag.multi[v] = 0
	ag.nodes = slices.DeleteFunc(ag.nodes, func(ni nodeIdx) bool {
		return ni.idx == v
	})
}

func (ag *adjGraph) addEdge(a, b, n int) {
	if a > b {
		a, b = b, a
	}
	ag.edges[a][b] += n
	ag.numEdges += n
}

func (ag *adjGraph) getEdge(a, b int) int {
	if a > b {
		a, b = b, a
	}
	return ag.edges[a][b]
}

func (ag *adjGraph) deleteEdges(a, b int) int {
	if a > b {
		a, b = b, a
	}
	n := ag.edges[a][b]
	ag.edges[a][b] = 0
	ag.numEdges -= n
	return n
}

func (ag *adjGraph) numNodes() int {
	return len(ag.nodes)
}

func (ag *adjGraph) partitionSizes() []int64 {
	var sizes []int64
	for _, n := range ag.multi {
		if n == 0 {
			continue
		}
		sizes = append(sizes, int64(n))
	}
	return sizes
}

// type graphCut struct {
// 	numEdges int
// 	s, t     string
// }

// func (g *graph) minCutScore() int64 {
// 	// Stoer-Wagner
// 	orig := g.clone()
// 	minCut := graphCut{
// 		numEdges: math.MaxInt,
// 	}
// 	for len(g.nodes) > 1 {
// 		var v string
// 		for x := range g.nodes {
// 			v = x
// 			break
// 		}
// 		cut := g.minimumCutPhase()
// 		if cut.numEdges < minCut.numEdges {
// 			minCut = cut
// 		}
// 		g.merge(cut.s, cut.t)
// 	}
// 	if minCut.numEdges != 3 {
// 		panic(fmt.Sprintf("min cut was %d", minCut.numEdges))
// 	}
// 	part0, part1 := orig.componentSize(minCut.s), orig.componentSize(minCut.t)
// 	return part0 * part1
// }

// func (g *graph) minimumCutPhase() graphCut {
// 	panic("unimplemented")
// }

// func (g *graph) merge(s, t string) {
// 	panic("unimplemented")
// }

// func (g *graph) componentSize(v string) int64 {
// 	seen := map[string]struct{}{v: {}}
// 	q := []string{v}
// 	for len(q) > 0 {
// 		v := SlicePop(&q)
// 		for dst := range g.nodes[v] {
// 			if _, ok := seen[dst]; !ok {
// 				seen[dst] = struct{}{}
// 				q = append(q, dst)
// 			}
// 		}
// 	}
// 	return int64(len(seen))
// }
