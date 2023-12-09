package main

import (
	"slices"
	"strings"
)

func init() {
	addSolutions(9, problem9)
}

func problem9(ctx *problemContext) {
	var seqs [][]int64
	scanner := ctx.scanner()
	for scanner.scan() {
		fields := strings.Fields(scanner.text())
		seq := SliceMap(fields, parseInt)
		seqs = append(seqs, seq)

	}
	ctx.reportLoad()

	var part1Sum int64
	for _, seq := range seqs {
		part1Sum += predictNext(seq)
	}
	ctx.reportPart1(part1Sum)

	var part2Sum int64
	for _, seq := range seqs {
		part2Sum += predictPrev(seq)
	}
	ctx.reportPart2(part2Sum)
}

func predictNext(seq []int64) int64 {
	seq = slices.Clone(seq)
	end := len(seq)
	for end > 0 {
		nonzero := false
		for i := 0; i < end-1; i++ {
			seq[i] = seq[i+1] - seq[i]
			if seq[i] != 0 {
				nonzero = true
			}
		}
		end--
		if nonzero {
			continue
		}
		var next int64
		for i := end; i < len(seq); i++ {
			next += seq[i]
		}
		return next
	}
	panic("seq too short to predict")
}

func predictPrev(seq []int64) int64 {
	seq = slices.Clone(seq)
	start := 0
	for start < len(seq) {
		nonzero := false
		for i := len(seq) - 1; i > start; i-- {
			seq[i] -= seq[i-1]
			if seq[i] != 0 {
				nonzero = true
			}
		}
		start++
		if nonzero {
			continue
		}
		var next int64
		for i := start; i >= 0; i-- {
			next = seq[i] - next
		}
		return next
	}
	panic("seq too short to predict")
}
