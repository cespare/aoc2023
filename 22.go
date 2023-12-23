package main

import (
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

func init() {
	addSolutions(22, problem22)
}

func problem22(ctx *problemContext) {
	var t brickTower
	scanner := ctx.scanner()
	for scanner.scan() {
		t.parse(scanner.text())
	}
	t.init()
	ctx.reportLoad()

	t.settle()
	ctx.reportPart1(t.part1())
	ctx.reportPart1(t.part2())
}

type brickTower struct {
	m      map[vec3]*brick
	bricks []*brick
}

func (t *brickTower) clone() *brickTower {
	t1 := &brickTower{
		m:      make(map[vec3]*brick),
		bricks: make([]*brick, len(t.bricks)),
	}
	for i, b := range t.bricks {
		t1.bricks[i] = b.clone()
	}
	for v, b := range t.m {
		t1.m[v] = t1.bricks[b.id]
	}
	return t1
}

func (t *brickTower) parse(s string) {
	if t.m == nil {
		t.m = make(map[vec3]*brick)
	}
	b := parseBrick(s)
	b.id = len(t.bricks)
	for _, v := range b.points {
		t.m[v] = b
	}
	t.bricks = append(t.bricks, b)
}

func (t *brickTower) init() {
	slices.SortFunc(t.bricks, func(b0, b1 *brick) int {
		return int(b0.points[0].z - b1.points[0].z)
	})
}

type brick struct {
	id       int
	points   []vec3
	vertical bool
}

func (b *brick) clone() *brick {
	return &brick{
		id:       b.id,
		points:   slices.Clone(b.points),
		vertical: b.vertical,
	}
}

func parseBrick(s string) *brick {
	s0, s1, ok := strings.Cut(s, "~")
	if !ok {
		panic("bad")
	}
	v0, v1 := parseVec3(s0), parseVec3(s1)
	if v0.x > v1.x || v0.y > v1.y || v0.z > v1.z {
		v0, v1 = v1, v0
	}
	d := v1.sub(v0)
	if d.x > 0 {
		d.x = 1
	}
	if d.y > 0 {
		d.y = 1
	}
	if d.z > 0 {
		d.z = 1
	}
	v := v0
	var b brick
	for {
		b.points = append(b.points, v)
		if v == v1 {
			break
		}
		v = v.add(d)
	}
	// Put the lowest point first for later purposes.
	slices.SortFunc(b.points, func(p0, p1 vec3) int {
		return int(p0.z - p1.z)
	})
	if len(b.points) > 1 && b.points[0].z != b.points[1].z {
		b.vertical = true
	}
	return &b
}

func parseVec3(s string) vec3 {
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		panic(len(parts))
	}
	return vec3{
		x: parseInt(parts[0]),
		y: parseInt(parts[1]),
		z: parseInt(parts[2]),
	}
}

func (t *brickTower) settle() {
	for {
		changed := false
		for _, b := range t.bricks {
			changed = changed || t.settle1(b)
		}
		if !changed {
			return
		}
	}
}

func (t *brickTower) settle1(b *brick) (moved bool) {
	var dz int64
	if b.vertical {
		dz = t.lowerBlock(b.points[0])
		if dz == 0 {
			return false
		}
	} else {
		dz = -(b.points[0].z - 1)
		for _, p := range b.points {
			dz = max(dz, t.lowerBlock(p))
			if dz == 0 {
				return false
			}
		}
	}
	dv := vec3{0, 0, dz}
	t.remove(b)
	for i, v := range b.points {
		v = v.add(dv)
		t.m[v] = b
		b.points[i] = v
	}
	return true
}

func (t *brickTower) lowerBlock(v vec3) int64 {
	dz := int64(-1)
	for ; ; dz-- {
		n := v
		n.z += dz
		if n.z == 0 {
			break
		}
		if _, ok := t.m[n]; ok {
			break
		}
	}
	return dz + 1
}

func (t *brickTower) part1() int {
	singles := make(map[*brick]struct{})
	for _, b := range t.bricks {
		supporters := t.restsOn(b)
		if len(supporters) != 1 {
			continue
		}
		singles[supporters[0]] = struct{}{}
	}
	return len(t.bricks) - len(singles)
}

func (t *brickTower) restsOn(b *brick) []*brick {
	supporters := make(map[*brick]struct{})
	for _, p := range b.points {
		p = p.add(vec3{0, 0, -1})
		if b1, ok := t.m[p]; ok {
			if b1 != b {
				supporters[b1] = struct{}{}
			}
		}
	}
	return maps.Keys(supporters)
}

func (t *brickTower) part2() int {
	var n int
	for i := range t.bricks {
		t1 := t.clone()
		b := t1.bricks[i]
		t1.remove(b)
		moved := make(map[*brick]struct{})
		for {
			changed := false
			for _, b1 := range t1.bricks {
				if b1 == b {
					continue
				}
				if t1.settle1(b1) {
					moved[b1] = struct{}{}
					changed = true
				}
			}
			if !changed {
				break
			}
		}
		n += len(moved)
	}
	return n
}

func (t *brickTower) remove(b *brick) {
	for _, v := range b.points {
		delete(t.m, v)
	}
}
