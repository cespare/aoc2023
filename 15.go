package main

import (
	"slices"
	"strings"
)

func init() {
	addSolutions(15, problem15)
}

func problem15(ctx *problemContext) {
	input := ctx.readAll()
	ctx.reportLoad()

	var part1Sum uint64
	for _, s := range strings.Split(strings.TrimSpace(string(input)), ",") {
		part1Sum += uint64(hash(s))
	}
	ctx.reportPart1(part1Sum)

	var h hashMap
	for _, s := range strings.Split(strings.TrimSpace(string(input)), ",") {
		label, lensStr, ok := strings.Cut(s, "=")
		if !ok {
			label = strings.TrimSuffix(s, "-")
			h.delete(label)
			continue
		}
		lens := int(parseInt(lensStr))
		h.set(label, lens)
	}
	ctx.reportPart2(h.focusingPower())
}

func hash(s string) uint8 {
	var n int
	for i := 0; i < len(s); i++ {
		n += int(s[i])
		n *= 17
		n = n % 256
	}
	return uint8(n)
}

type hashMap struct {
	boxes [256][]hashSlot
}

type hashSlot struct {
	label string
	lens  int // 1 - 9
}

func (h *hashMap) delete(label string) {
	i := hash(label)
	box := h.boxes[i]
	for j, slot := range box {
		if slot.label == label {
			h.boxes[i] = slices.Delete(box, j, j+1)
			return
		}
	}
}

func (h *hashMap) set(label string, lens int) {
	i := hash(label)
	box := h.boxes[i]
	for j, slot := range box {
		if slot.label == label {
			box[j] = hashSlot{label, lens}
			return
		}
	}
	h.boxes[i] = append(box, hashSlot{label, lens})
}

func (h *hashMap) focusingPower() int64 {
	var p int64
	for i, box := range h.boxes {
		for j, slot := range box {
			p += int64(i+1) * int64(j+1) * int64(slot.lens)
		}

	}
	return p
}
