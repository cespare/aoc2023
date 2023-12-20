package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kr/pretty"
)

func init() {
	addSolutions(20, problem20)
}

func problem20(ctx *problemContext) {
	var m pulseMachine
	scanner := ctx.scanner()
	for scanner.scan() {
		m.parseLine(scanner.text())
	}
	m.init()
	ctx.reportLoad()

	pretty.Println(m)

	if err := os.WriteFile("out.dot", m.toDot(), 0o644); err != nil {
		log.Fatal(err)
	}

	// TODO: Describe how to use the visualization to get the input sizes.
	// Then we just take the LCM.

	m.run2()

	// var lo, hi int64
	// for i := 0; i < 1000; i++ {
	// 	lo1, hi1, _ := m.run()
	// 	lo += lo1
	// 	hi += hi1
	// }
	// ctx.reportPart1(lo * hi)

	// for i := 1; ; i++ {
	// 	if i%1e5 == 0 {
	// 		fmt.Println(i)
	// 	}
	// 	_, _, outputLo := m.run()
	// 	if outputLo {
	// 		ctx.reportPart2(i)
	// 		break
	// 	}
	// }
}

type pulseMachine struct {
	broadcast []string
	ms        map[string]module

	inputs map[string][]string
}

func (m *pulseMachine) parseLine(s string) {
	src, dst, ok := strings.Cut(s, " -> ")
	if !ok {
		panic("bad")
	}
	dsts := strings.Split(dst, ", ")
	if src == "broadcaster" {
		m.broadcast = dsts
		return
	}
	if m.ms == nil {
		m.ms = make(map[string]module)
		m.inputs = make(map[string][]string)
	}
	var mod module
	var name string
	if name, ok = strings.CutPrefix(src, "%"); ok {
		mod = &moduleFF{dsts: dsts}
	} else if name, ok = strings.CutPrefix(src, "&"); ok {
		mod = &moduleConj{dsts: dsts}

	} else {
		panic("bad")
	}
	for _, dst := range dsts {
		m.inputs[dst] = append(m.inputs[dst], name)
	}
	m.ms[name] = mod
}

func (m *pulseMachine) toDot() []byte {
	trimName := func(s string) string {
		return strings.Trim(s, "%&")
	}
	var b bytes.Buffer
	fmt.Fprintln(&b, "digraph G {")
	fmt.Fprintln(&b, "broadcaster [shape=box, style=filled]")
	for _, dst := range m.broadcast {
		fmt.Fprintf(&b, "broadcaster -> %s\n", trimName(dst))
	}
	for src, mod := range m.ms {
		if _, ok := mod.(*moduleConj); ok {
			fmt.Fprintf(&b, "%s [shape=box]\n", trimName(src))
		}
		for _, dst := range mod.dests() {
			fmt.Fprintf(&b, "%s -> %s\n", trimName(src), trimName(dst))
		}
	}
	fmt.Fprintln(&b, "}")
	return b.Bytes()
}

func (m *pulseMachine) init() {
	for name, mod := range m.ms {
		mod, ok := mod.(*moduleConj)
		if !ok {
			continue
		}
		mod.inputs = make(map[string]pulse)
		for _, dst := range m.inputs[name] {
			mod.inputs[dst] = pulseLo
		}
	}
}

type module interface {
	send(string, pulse) []pulse
	dests() []string
}

type pulse int

const (
	pulseLo pulse = iota
	pulseHi
)

func (p pulse) String() string {
	if p == pulseLo {
		return "lo"
	}
	return "hi"
}

type moduleFF struct {
	state bool
	dsts  []string
}

func (m *moduleFF) send(src string, p pulse) []pulse {
	if p == pulseHi {
		return nil
	}
	var out []pulse
	if m.state {
		out = []pulse{pulseLo}
	} else {
		out = []pulse{pulseHi}
	}
	m.state = !m.state
	return out
}

func (m *moduleFF) dests() []string { return m.dsts }

type moduleConj struct {
	inputs map[string]pulse
	dsts   []string
}

func (m *moduleConj) send(src string, p pulse) []pulse {
	m.inputs[src] = p
	for _, p1 := range m.inputs {
		if p1 == pulseLo {
			return []pulse{pulseHi}
		}
	}
	return []pulse{pulseLo}
}

func (m *moduleConj) dests() []string { return m.dsts }

type pulseAction struct {
	src string
	p   pulse
	dst string
}

func (m *pulseMachine) run() (lo, hi int64, outputLo bool) {
	lo++
	var q []pulseAction
	for _, dst := range m.broadcast {
		q = append(q, pulseAction{src: "broadcaster", p: pulseLo, dst: dst})
	}
	for len(q) > 0 {
		action := SlicePop(&q)
		if action.dst == "dn" {
			fmt.Println(action)
		}
		if action.p == pulseLo {
			lo++
		} else {
			hi++
		}
		// fmt.Println(action)
		// pretty.Println(m)
		mod := m.ms[action.dst]
		if mod == nil {
			if action.p == pulseLo {
				outputLo = true
			}
			continue
		}
		out := mod.send(action.src, action.p)
		for _, dst := range mod.dests() {
			for _, p := range out {
				q = append(q, pulseAction{src: action.dst, p: p, dst: dst})
			}
		}
	}
	return lo, hi, outputLo
}

func (m *pulseMachine) run2() {
	for i := 1; ; i++ {
		var q []pulseAction
		for _, dst := range m.broadcast {
			q = append(q, pulseAction{src: "broadcaster", p: pulseLo, dst: dst})
		}
		for len(q) > 0 {
			action := SlicePop(&q)
			if action.p == pulseLo {
				switch action.dst {
				// case "bd", "pm", "rs", "cc":
				case "xp", "fh", "dd", "fc":
					fmt.Println("steps:", i, "action:", action)
				}
			}
			mod := m.ms[action.dst]
			if mod == nil {
				continue
			}
			out := mod.send(action.src, action.p)
			for _, dst := range mod.dests() {
				for _, p := range out {
					q = append(q, pulseAction{src: action.dst, p: p, dst: dst})
				}
			}
		}
	}
}
