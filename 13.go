package main

import (
	"math/bits"
)

func init() {
	addSolutions(13, problem13)
}

func problem13(ctx *problemContext) {
	var maps []*mirrorMap
	m := new(mirrorMap)
	scanner := ctx.scanner()
	for scanner.scan() {
		row := scanner.text()
		if row == "" {
			if len(m.rows) > 0 {
				maps = append(maps, m)
				m = new(mirrorMap)
			}
			continue
		}
		m.addRow(row)

	}
	if len(m.rows) > 0 {
		maps = append(maps, m)
	}
	ctx.reportLoad()

	var part1Sum int64
	for _, m := range maps {
		part1Sum += m.summarize(0)
	}
	ctx.reportPart1(part1Sum)

	var part2Sum int64
	for _, m := range maps {
		part2Sum += m.summarize(1)
	}
	ctx.reportPart2(part2Sum)
}

type mirrorMap struct {
	rows []uint64 // bitmap
	cols int
}

func (m *mirrorMap) addRow(s string) {
	if m.cols == 0 {
		m.cols = len(s)
	}
	var n uint64
	for i, c := range s {
		if c == '#' {
			n |= 1 << i
		}
	}
	m.rows = append(m.rows, n)
}

func (m *mirrorMap) summarize(diff int) int64 {
	n := m.findHorizontalMirror(diff)
	if n < 0 {
		n = m.rotate().findHorizontalMirror(diff)
		if n < 0 {
			panic("cannot find mirror")
		}
	} else {
		n *= 100
	}
	return int64(n)
}

func (m *mirrorMap) findHorizontalMirror(diff int) int {
rowLoop:
	for y := 0; y < len(m.rows)-1; y++ {
		var d int
		y0 := y
		y1 := y + 1
		for y0 >= 0 && y1 < len(m.rows) {
			d += bits.OnesCount64(m.rows[y0] ^ m.rows[y1])
			if d > diff {
				continue rowLoop
			}
			y0--
			y1++
		}
		if d == diff {
			return y + 1
		}
	}
	return -1
}

func (m *mirrorMap) rotate() *mirrorMap {
	m1 := &mirrorMap{
		rows: make([]uint64, m.cols),
		cols: len(m.rows),
	}
	for y, row := range m.rows {
		for x := 0; x < m.cols; x++ {
			m1.rows[x] |= ((row >> x) & 1) << y
		}
	}
	return m1
}
