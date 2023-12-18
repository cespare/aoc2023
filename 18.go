package main

import (
	"bytes"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

func init() {
	addSolutions(18, problem18)
}

func problem18(ctx *problemContext) {
	var insts []digInstruction
	scanner := ctx.scanner()
	for scanner.scan() {
		inst := parseDigInstruction(scanner.text())
		insts = append(insts, inst)
	}
	ctx.reportLoad()

	t := newTrench(insts)
	t.flood()
	ctx.reportPart1(t.countArea())

	for i, inst := range insts {
		insts[i] = inst.fix()
	}
	t = newTrench(insts)
	t.flood()
	ctx.reportPart2(t.countArea())
}

type digInstruction struct {
	v     vec2
	n     int64
	color string
}

func (d digInstruction) fix() digInstruction {
	var inst digInstruction
	switch d.color[5] {
	case '3':
		inst.v = northv
	case '0':
		inst.v = eastv
	case '1':
		inst.v = southv
	case '2':
		inst.v = westv
	default:
		panic("bad")
	}
	n, err := strconv.ParseInt(d.color[:5], 16, 64)
	if err != nil {
		panic(err)
	}
	inst.n = n
	return inst
}

func parseDigInstruction(s string) digInstruction {
	var inst digInstruction
	parts := strings.Fields(s)
	switch parts[0][0] {
	case 'U':
		inst.v = northv
	case 'R':
		inst.v = eastv
	case 'D':
		inst.v = southv
	case 'L':
		inst.v = westv
	default:
		panic("bad")
	}
	inst.n = parseInt(parts[1])
	inst.color = strings.TrimSuffix(strings.TrimPrefix(parts[2], "(#"), ")")
	return inst
}

type trench struct {
	// All vectors are indexes into gridx/gridy.
	gridx []int64
	gridy []int64

	g grid[byte] // in grid coords
}

func newTrench(insts []digInstruction) *trench {
	setx := map[int64]struct{}{0: {}}
	sety := map[int64]struct{}{0: {}}
	var v vec2
	var vecs []vec2
	for _, inst := range insts {
		vecs = append(vecs, v)
		v = v.add(inst.v.mul(inst.n))
		setx[v.x] = struct{}{}
		sety[v.y] = struct{}{}
	}
	if v != vecs[0] {
		panic("not a loop")
	}
	xvals := maps.Keys(setx)
	slices.Sort(xvals)
	yvals := maps.Keys(sety)
	slices.Sort(yvals)
	var t trench
	for i, x := range xvals {
		t.gridx = append(t.gridx, x)
		if i == len(xvals)-1 || xvals[i+1] > x+1 {
			t.gridx = append(t.gridx, x+1)
		}
	}
	for i, y := range yvals {
		t.gridy = append(t.gridy, y)
		if i == len(yvals)-1 || yvals[i+1] > y+1 {
			t.gridy = append(t.gridy, y+1)
		}
	}
	for i := 0; i < len(t.gridy); i++ {
		t.g.addRow(bytes.Repeat([]byte("."), len(t.gridx)))
	}

	for i, inst := range insts {
		v := vecs[i]
		cur := t.mapToGrid(v)
		if i < len(vecs)-1 {
			v = vecs[i+1]
		} else {
			v = vecs[0]
		}
		next := t.mapToGrid(v)

		for cur != next {
			t.g.set(cur, '#')
			cur = cur.add(inst.v)
		}
	}
	return &t
}

func (t *trench) mapToGrid(v vec2) vec2 {
	ix, _ := slices.BinarySearch(t.gridx, v.x)
	iy, _ := slices.BinarySearch(t.gridy, v.y)
	return vec2{int64(ix), int64(iy)}
}

func (t *trench) gridToMap(v vec2) (start, end vec2) {
	start = vec2{t.gridx[v.x], t.gridy[v.y]}
	end = vec2{t.gridx[v.x+1], t.gridy[v.y+1]}
	return start, end
}

func (t *trench) flood() {
	outside := make(map[vec2]struct{})
	t.g.forEach(func(v vec2, c byte) bool {
		if c != '.' {
			return true
		}
		if _, ok := outside[v]; ok {
			return true
		}
		region := make(map[vec2]struct{})
		out := t.floodRec(v, region)
		for v1 := range region {
			if out {
				outside[v1] = struct{}{}
			} else {
				t.g.set(v1, '#')
			}
		}
		return true
	})
}

func (t *trench) floodRec(v vec2, region map[vec2]struct{}) (outside bool) {
	if !t.g.contains(v) {
		return true
	}
	if t.g.at(v) != '.' {
		return false
	}
	if _, ok := region[v]; ok {
		return false
	}
	region[v] = struct{}{}
	for _, n := range v.neighbors4() {
		outside = outside || t.floodRec(n, region)
	}
	return outside
}

func (t *trench) countArea() int64 {
	var n int64
	t.g.forEach(func(v vec2, c byte) bool {
		if c == '.' {
			return true
		}
		start, end := t.gridToMap(v)
		d := end.sub(start)
		n += d.x * d.y
		return true
	})
	return n
}
