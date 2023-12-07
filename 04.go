package main

import (
	"log"
	"strings"
)

func init() {
	addSolutions(4, problem4)
}

func problem4(ctx *problemContext) {
	cs := scanSlice(ctx, func(s string) *card {
		c, ok := parseCard(s)
		if !ok {
			log.Fatalf("Bad card: %q", s)
		}
		return c
	})
	ctx.reportLoad()

	var part1Score int64
	for _, c := range cs {
		part1Score += c.score()
	}
	ctx.reportPart1(part1Score)

	numCards := make(map[int64]int64)
	for _, c := range cs {
		numCards[c.id] = 1
	}
	for _, c := range cs {
		n := numCards[c.id]
		m := c.numMatches()
		for i := int64(0); i < m; i++ {
			id := c.id + i + 1
			numCards[id] += n
		}
	}
	var totalCards int64
	for _, n := range numCards {
		totalCards += n
	}
	ctx.reportPart2(totalCards)
}

type card struct {
	id      int64
	winning []int64
	have    []int64
	matches int64
}

func (c *card) score() int64 {
	w := make(map[int64]struct{})
	for _, n := range c.winning {
		w[n] = struct{}{}
	}
	var s int64
	for _, n := range c.have {
		if _, ok := w[n]; ok {
			if s == 0 {
				s = 1
			} else {
				s *= 2
			}
		}
	}
	return s
}

func (c *card) numMatches() int64 {
	if c.matches >= 0 {
		return c.matches
	}
	w := make(map[int64]struct{})
	for _, n := range c.winning {
		w[n] = struct{}{}
	}
	c.matches = 0
	for _, n := range c.have {
		if _, ok := w[n]; ok {
			c.matches++
		}
	}
	return c.matches
}

func parseCard(s string) (*card, bool) {
	c := &card{matches: -1}
	label, rest, ok := strings.Cut(s, ": ")
	if !ok {
		return nil, false
	}
	_, idStr, ok := strings.Cut(label, " ")
	if !ok {
		return nil, false
	}
	c.id = parseInt(strings.TrimSpace(idStr))
	winningStr, haveStr, ok := strings.Cut(rest, " | ")
	if !ok {
		return nil, false
	}
	c.winning = parseInts(winningStr)
	c.have = parseInts(haveStr)
	return c, true
}

func parseInts(s string) []int64 {
	var ns []int64
	for _, field := range strings.Fields(s) {
		ns = append(ns, parseInt(field))
	}
	return ns
}
