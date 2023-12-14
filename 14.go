package main

import (
	"crypto/sha256"
)

func init() {
	addSolutions(14, problem14)
}

func problem14(ctx *problemContext) {
	var d dish
	scanner := ctx.scanner()
	for scanner.scan() {
		d.g.addRow([]byte(scanner.text()))
	}
	ctx.reportLoad()

	d1 := d.clone()
	d1.tilt(north)
	ctx.reportPart1(d1.load())

	d2 := d.clone()
	m := make(map[string]int64)
	targ := int64(-1)
	for i := int64(0); i < 1e9; i++ {
		if targ >= 0 {
			if i == targ {
				break
			}
		} else {
			key := d2.cacheKey()
			if prev, ok := m[key]; ok {
				cycleLen := i - prev
				targ = ((1e9 - prev) % cycleLen) + i
			} else {
				m[key] = int64(i)
			}
		}
		d2.spin()
	}
	ctx.reportPart2(d2.load())
}

type dish struct {
	g grid[byte]
}

func (d *dish) clone() *dish {
	g1 := d.g.clone()
	return &dish{g: *g1}
}

func (d *dish) cacheKey() string {
	h := sha256.New()
	for _, row := range d.g.g {
		h.Write(row)
	}
	sum := h.Sum(nil)
	return string(sum)
}

func (d *dish) spin() {
	d.tilt(north)
	d.tilt(west)
	d.tilt(south)
	d.tilt(east)
}

func (d *dish) tilt(dir cardinal) {
	dv := cardToVec(dir)
	majorv := dv.mul(-1)
	absv := vec2{abs(dv.x), abs(dv.y)}
	minorv := vec2{1, 1}.sub(absv)
	select1 := func(x int64) int64 {
		if x == 1 {
			return 1
		}
		return 0
	}
	v := vec2{d.g.cols - 1, d.g.rows - 1}.eltMul(vec2{select1(dv.x), select1(dv.y)})
	for d.g.contains(v) {
		if d.g.at(v) == 'O' {
			d.moveRock(v, dv)
		}
		v = v.add(minorv)
		if !d.g.contains(v) {
			v = v.eltMul(absv)
			v = v.add(majorv)
		}
	}
}

func (d *dish) moveRock(v, dv vec2) {
	v0 := v
	for {
		v1 := v.add(dv)
		if !d.g.contains(v1) {
			break
		}
		if d.g.at(v1) != '.' {
			break
		}
		v = v1
	}
	d.g.set(v0, '.')
	d.g.set(v, 'O')
}

func (d *dish) load() int64 {
	var load int64
	d.g.forEach(func(v vec2, c byte) bool {
		if c == 'O' {
			load += d.g.rows - v.y
		}
		return true

	})
	return load
}
