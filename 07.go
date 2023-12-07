package main

import (
	"slices"
	"strings"
)

func init() {
	addSolutions(7, problem7)
}

func problem7(ctx *problemContext) {
	hbs := scanSlice(ctx, func(s string) handBid {
		hand, bidStr, ok := strings.Cut(s, " ")
		if !ok {
			panic("bad")
		}
		return handBid{
			hand: pokerHand([]byte(hand)),
			bid:  parseInt(bidStr),
		}
	})
	ctx.reportLoad()

	slices.SortFunc(hbs, func(hb0, hb1 handBid) int {
		return hb0.hand.cmp(hb1.hand)
	})
	calcScore := func() int64 {
		var score int64
		for i, hb := range hbs {
			score += int64(i+1) * hb.bid
		}
		return score
	}
	ctx.reportPart1(calcScore())

	slices.SortFunc(hbs, func(hb0, hb1 handBid) int {
		return hb0.hand.cmpJ(hb1.hand)
	})
	ctx.reportPart2(calcScore())
}

type handBid struct {
	hand pokerHand
	bid  int64
}

type pokerHand [5]byte

func (h pokerHand) cmp(h1 pokerHand) int {
	p0, p1 := h.power(), h1.power()
	if p0 != p1 {
		return p0 - p1
	}
	for i, c := range h {
		p0, p1 := cardPower(c), cardPower(h1[i])
		if p0 != p1 {
			return p0 - p1
		}
	}
	return 0
}

func (h pokerHand) power() int {
	m := make(map[byte]int)
	for _, c := range h {
		m[c]++
	}
	var counts []int
	for _, n := range m {
		counts = append(counts, n)
	}
	slices.SortFunc(counts, func(n0, n1 int) int { return n1 - n0 })
	switch counts[0] {
	case 5:
		return 6 // five of a kind
	case 4:
		return 5 // four of a kind
	case 3:
		if counts[1] == 2 {
			return 4 // full house
		} else {
			return 3 // three of a kind
		}
	case 2:
		if counts[1] == 2 {
			return 2 // two pair
		} else {
			return 1 // pair
		}
	default:
		return 0 // high card
	}
}

func cardPower(c byte) int {
	switch c {
	case 'A':
		return 14
	case 'K':
		return 13
	case 'Q':
		return 12
	case 'J':
		return 11
	case 'T':
		return 10
	default:
		return int(c - '0')
	}
}

func (h pokerHand) cmpJ(h1 pokerHand) int {
	p0, p1 := h.powerJ(), h1.powerJ()
	if p0 != p1 {
		return p0 - p1
	}
	for i, c := range h {
		p0, p1 := cardPowerJ(c), cardPowerJ(h1[i])
		if p0 != p1 {
			return p0 - p1
		}
	}
	return 0
}

func (h pokerHand) powerJ() int {
	m := make(map[byte]int)
	for _, c := range h {
		m[c]++
	}
	jokers := m['J']
	if jokers == 5 {
		return 6 // five of a kind
	}
	delete(m, 'J')
	var counts []int
	for _, n := range m {
		counts = append(counts, n)
	}
	slices.SortFunc(counts, func(n0, n1 int) int { return n1 - n0 })
	switch counts[0] + jokers {
	case 5:
		return 6 // five of a kind
	case 4:
		return 5 // four of a kind
	case 3:
		if counts[1] == 2 {
			return 4 // full house
		} else {
			return 3 // three of a kind
		}
	case 2:
		if counts[1] == 2 {
			return 2 // two pair
		} else {
			return 1 // pair
		}
	default:
		return 0 // high card
	}
}

func cardPowerJ(c byte) int {
	switch c {
	case 'A':
		return 13
	case 'K':
		return 12
	case 'Q':
		return 11
	case 'J':
		return 1
	case 'T':
		return 10
	default:
		return int(c - '0')
	}
}
