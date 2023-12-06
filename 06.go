package main

import (
	"strings"
)

func init() {
	addSolutions(6, problem6)
}

func problem6(ctx *problemContext) {
	scanner := ctx.scanner()
	var s [2]string
	scanner.scan()
	s[0] = scanner.text()
	scanner.scan()
	s[1] = scanner.text()
	rec := parseRaceRecords(s)
	ctx.reportLoad()

	var part1Prod int64 = 1
	for _, r := range rec {
		part1Prod *= r.waysToWin()
	}
	ctx.reportPart1(part1Prod)

	// Time:              44899691
	// Distance:   277113618901768
	//
	// h * (t-h) = d
	// ht - h^2 = d
	//
	// Wolfram Alpha says the roots are:
	//
	// h≈7387244.7
	// h≈37512446.3
	//
	// 7387245 ≤ h ≤ 37512446
	ctx.reportPart2(37512446 - 7387245 + 1)
}

type raceRecord struct {
	t int64
	d int64
}

func parseRaceRecords(s [2]string) []raceRecord {
	tf := strings.Fields(strings.TrimPrefix(s[0], "Time:"))
	td := strings.Fields(strings.TrimPrefix(s[1], "Distance:"))
	if len(tf) != len(td) {
		panic("bad")
	}
	var record []raceRecord
	for i, ts := range tf {
		record = append(record, raceRecord{
			t: parseInt(ts),
			d: parseInt(td[i]),
		})
	}
	return record
}

func (r raceRecord) waysToWin() int64 {
	var ways int64
	for hold := int64(1); hold <= r.t; hold++ {
		result := hold * (r.t - hold)
		if result > r.d {
			ways++
		}
	}
	return ways
}

// Time:              44899691
// Distance:   277113618901768
//
//
// h * (t-h) = d
// ht - h^2 = d
//
//
// roots: 7387245, 37512446
//
// h≈7.38724465976984×10^6
// h≈3.75124463402302×10^7
//
// h >= 7387246
// h <= 37512445
//
// 30125200 is wrong
//
// Time:      71530
// Distance:  940200
