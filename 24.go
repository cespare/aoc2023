package main

import (
	"strings"
)

func init() {
	addSolutions(24, problem24)
}

func problem24(ctx *problemContext) {
	var stones []hailstone
	scanner := ctx.scanner()
	for scanner.scan() {
		stones = append(stones, parseHailstone(scanner.text()))
	}
	ctx.reportLoad()

	const bound0 = 200e12
	const bound1 = 400e12

	var part1 int
	for i, stone0 := range stones {
		for j := i + 1; j < len(stones); j++ {
			stone1 := stones[j]
			_, ok := intersection2d(stone0, stone1, bound0, bound1)
			if ok {
				part1++
			}
		}
	}
	ctx.reportPart1(part1)

	// For part 2:
	//
	// We are looking for (px, py, pz), (vx, vy,vz)
	// We are given hailstones a, b, c, ...
	// Let's say that the thrown stone hits hailstone a at ta, b at tb, ...
	//
	// If we consider only hailstone a, b, and c, we can make this nonlinear
	// system of equations:
	//
	//   px + vx*ta = pax + vfx*ta
	//   py + vy*ta = pay + vfy*ta
	//   pz + vz*ta = paz + vfz*ta
	//
	//   px + vx*tb = pbx + vgx*tb
	//   py + vy*tb = pby + vgy*tb
	//   pz + vz*tb = pbz + vgz*tb
	//
	//   px + vx*tc = pcx + vhx*tc
	//   py + vy*tc = pcy + vhy*tc
	//   pz + vz*tc = pcz + vhz*tc
	//
	// We have 9 equations and 9 variables. After attempting with wolfram
	// alpha and octave, I finally got it working in Z3:
	//
	//   (declare-const px Real)
	//   (declare-const py Real)
	//   (declare-const pz Real)
	//   (declare-const vx Real)
	//   (declare-const vy Real)
	//   (declare-const vz Real)
	//   (declare-const ta Real)
	//   (declare-const tb Real)
	//   (declare-const tc Real)
	//   (assert (= (+ px (* vx ta)) (+ 308205470708820 (* 42  ta))))
	//   (assert (= (+ py (* vy ta)) (+ 82023714100543  (* 274 ta))))
	//   (assert (= (+ pz (* vz ta)) (- 475164418926765 (* 194 ta))))
	//   (assert (= (+ px (* vx tb)) (+ 242904857760501 (* 147 tb))))
	//   (assert (= (+ py (* vy tb)) (- 351203053017504 (* 69  tb))))
	//   (assert (= (+ pz (* vz tb)) (+ 247366253386570 (* 131 tb))))
	//   (assert (= (+ px (* vx tc)) (- 258124591360173 (* 84  tc))))
	//   (assert (= (+ py (* vy tc)) (- 252205185038992 (* 5   tc))))
	//   (assert (= (+ pz (* vz tc)) (+ 113896142591148 (* 409 tc))))
	//   (check-sat)
	//   (get-model)
	//   sat
	//   (
	//     (define-fun tb () Real
	//       854610412103.0)
	//     (define-fun vy () Real
	//       75.0)
	//     (define-fun vz () Real
	//       221.0)
	//     (define-fun vx () Real
	//       245.0)
	//     (define-fun tc () Real
	//       300825392054.0)
	//     (define-fun ta () Real
	//       734248440071.0)
	//     (define-fun px () Real
	//       159153037374407.0)
	//     (define-fun pz () Real
	//       170451316297300.0)
	//     (define-fun py () Real
	//       228139153674672.0)
	//   )
}

type hailstone struct {
	p vec3
	v vec3
}

func parseHailstone(s string) hailstone {
	s = strings.ReplaceAll(s, " ", "")
	ps, vs, ok := strings.Cut(s, "@")
	if !ok {
		panic("bad")
	}
	pp := strings.Split(ps, ",")
	vp := strings.Split(vs, ",")
	return hailstone{
		p: vec3{
			parseInt(pp[0]),
			parseInt(pp[1]),
			parseInt(pp[2]),
		},
		v: vec3{
			parseInt(vp[0]),
			parseInt(vp[1]),
			parseInt(vp[2]),
		},
	}
}

func intersection2d(h0, h1 hailstone, bound0, bound1 int64) (vec2f, bool) {
	// intersection at p0 + v0*u0
	//
	// p0 + v0*u0 = p1 + v1*u1
	//
	// px0 + vx0*u0 = px1 + vx1*u1
	// py0 + vy0*u0 = py1 + vy1*u1
	//
	// u1 = (px0 + vx0*u0 - px1) / vx1
	//
	// py0 + vy0*u0 = py1 + vy1*(px0 + vx0*u0 - px1)/vx1
	//
	// py0 + vy0*u0 = py1 + (vy1*px0 + vy1*vx0*u0 - vy1*px1)/vx1
	//
	// py0 + vy0*u0 = py1 + vy1*px0/vx1 + (vy1*vx0/vx1)*u0 - vy1*px1/vx1
	//
	// vy0*u0 - (vy1*vx0/vx1)*u0 = py1 + vy1*px0/vx1 - vy1*px1/vx1 - py0
	//
	//      py1 + vy1*px0/vx1 - vy1*px1/vx1 - py0
	// u0 = -------------------------------------
	//             vy0 - vy1*vx0/vx1

	var (
		px0 = float64(h0.p.x)
		py0 = float64(h0.p.y)
		vx0 = float64(h0.v.x)
		vy0 = float64(h0.v.y)
		px1 = float64(h1.p.x)
		py1 = float64(h1.p.y)
		vx1 = float64(h1.v.x)
		vy1 = float64(h1.v.y)
	)
	if vx0 == 0 {
		panic("vertical")
	}
	if vx1 == 0 {
		panic("vertical")
	}
	slope0 := vy0 / vx0
	slope1 := vy1 / vx1
	if slope0 == slope1 {
		return vec2f{}, false // parallel
	}
	u0 := (py1 + vy1*px0/vx1 - vy1*px1/vx1 - py0) / (vy0 - vy1*vx0/vx1)
	if u0 < 0 {
		return vec2f{}, false
	}
	u1 := (px0 + vx0*u0 - px1) / vx1
	if u1 < 0 {
		return vec2f{}, false
	}
	x := px0 + vx0*u0
	if x < float64(bound0) || x > float64(bound1) {
		return vec2f{}, false
	}
	y := py0 + vy0*u0
	if y < float64(bound0) || y > float64(bound1) {
		return vec2f{}, false
	}
	return vec2f{x, y}, true
}

type vec2f struct {
	x float64
	y float64
}
