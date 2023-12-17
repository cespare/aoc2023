package main

import (
	"github.com/cespare/next/container/heap"
)

func init() {
	addSolutions(17, problem17)
}

func problem17(ctx *problemContext) {
	var c city
	scanner := ctx.scanner()
	for scanner.scan() {
		c.g.addRow([]byte(scanner.text()))
	}
	ctx.reportLoad()

	ctx.reportPart1(c.bestPath(1, 3))
	ctx.reportPart2(c.bestPath(4, 10))
}

type city struct {
	g grid[byte]
}

type cityState struct {
	p vec2
	v vec2 // previous move (in reverse search)
	n int  // number of consecutive steps taken in the v direction
}

type cityQueueState struct {
	state cityState
	cost  int64
	idx   int
}

func (c *city) bestPath(minMove, maxMove int) int64 {
	h := heap.New(func(c0, c1 *cityQueueState) bool {
		return c0.cost < c1.cost
	})
	byState := make(map[cityState]*cityQueueState)
	h.SetIndex = func(e *cityQueueState, i int) {
		e.idx = i
	}
	seed := &cityQueueState{
		state: cityState{
			p: vec2{0, 0},
		},
	}
	targ := vec2{c.g.cols - 1, c.g.rows - 1}
	byState[seed.state] = seed
	h.Push(seed)
	for h.Len() > 0 {
		e := h.Pop()
		if e.state.p == targ && e.state.n >= minMove {
			return e.cost
		}
		delete(byState, e.state)

		for _, v := range nesw {
			if v == e.state.v.mul(-1) {
				continue
			}
			if e.state.n == maxMove && v == e.state.v {
				continue
			}
			if e.state.n > 0 && e.state.n < minMove && v != e.state.v {
				continue
			}
			p := e.state.p.add(v)
			if !c.g.contains(p) {
				continue
			}
			state := cityState{
				p: p,
				v: v,
				n: 1,
			}
			if v == e.state.v {
				state.n = e.state.n + 1
			}
			cost := e.cost + int64(c.g.at(p)-'0')
			if e1, ok := byState[state]; ok {
				if cost < e1.cost {
					e1.cost = cost
					h.Fix(e1.idx)
				}
				continue
			}
			qs := &cityQueueState{state: state, cost: cost}
			byState[state] = qs
			h.Push(qs)
		}
	}
	panic("no solution")
}
